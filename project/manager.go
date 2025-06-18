package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// ProjectStorage defines the interface for project persistence
type ProjectStorage interface {
	SaveProjects(projectsJSON json.RawMessage) error
	GetProjects() json.RawMessage
	DeleteProject(projectID string) error
	SetActiveProject(projectID string) error
	GetActiveProject() string
	// Project history methods
	SaveProjectHistory(history *ProjectHistory) error
	GetProjectHistory() *ProjectHistory
}

// ProjectManager manages multiple projects and their state
type ProjectManager struct {
	projects      map[string]*Project
	activeProject *Project
	storage       ProjectStorage
	history       *ProjectHistory
}

// NewProjectManager creates a new project manager with the given storage backend
func NewProjectManager(storage ProjectStorage) (*ProjectManager, error) {
	if storage == nil {
		return nil, fmt.Errorf("storage cannot be nil")
	}

	pm := &ProjectManager{
		projects: make(map[string]*Project),
		storage:  storage,
		history:  storage.GetProjectHistory(),
	}

	// Load existing projects from storage
	if err := pm.loadProjects(); err != nil {
		return nil, fmt.Errorf("failed to load projects: %w", err)
	}

	// Set active project if one was stored
	activeProjectID := storage.GetActiveProject()
	if activeProjectID != "" {
		pm.setActiveProjectByID(activeProjectID)
	}

	return pm, nil
}

// AddProject adds a new project to the manager or reactivates existing one
func (pm *ProjectManager) AddProject(path, name string) (*Project, error) {
	// Clean path for comparison
	cleanPath := filepath.Clean(path)
	
	// Check if project with same path already exists - if so, reactivate it
	for _, existing := range pm.projects {
		if existing.Path == cleanPath {
			// Reactivate existing project
			existing.LastAccessed = time.Now()
			pm.SetActiveProject(existing.ID)
			return existing, nil
		}
	}

	// Create new project
	project, err := NewProject(cleanPath, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Validate path exists
	if _, err := os.Stat(project.Path); os.IsNotExist(err) {
		return nil, fmt.Errorf("project path does not exist: %s", project.Path)
	}

	// Add to manager
	pm.projects[project.ID] = project

	// If this is the first project, make it active
	if len(pm.projects) == 1 {
		pm.SetActiveProject(project.ID)
	}

	// Save to storage
	if err := pm.saveProjects(); err != nil {
		// Remove from memory if save failed
		delete(pm.projects, project.ID)
		return nil, fmt.Errorf("failed to save project: %w", err)
	}

	return project, nil
}

// GetProject retrieves a project by ID
func (pm *ProjectManager) GetProject(projectID string) (*Project, bool) {
	project, exists := pm.projects[projectID]
	return project, exists
}

// GetActiveProject returns the currently active project
func (pm *ProjectManager) GetActiveProject() *Project {
	return pm.activeProject
}

// SetActiveProject sets the active project by ID
func (pm *ProjectManager) SetActiveProject(projectID string) error {
	project, exists := pm.projects[projectID]
	if !exists {
		return fmt.Errorf("project not found: %s", projectID)
	}

	// Deactivate current active project
	if pm.activeProject != nil {
		pm.activeProject.SetInactive()
	}

	// Set new active project
	pm.activeProject = project
	project.SetActive()

	// Save to storage
	if err := pm.storage.SetActiveProject(projectID); err != nil {
		return fmt.Errorf("failed to save active project: %w", err)
	}

	return nil
}

// setActiveProjectByID is an internal method that doesn't save to storage
func (pm *ProjectManager) setActiveProjectByID(projectID string) {
	if project, exists := pm.projects[projectID]; exists {
		if pm.activeProject != nil {
			pm.activeProject.SetInactive()
		}
		pm.activeProject = project
		project.SetActive()
	}
}

// ListProjects returns all projects sorted by last accessed time (most recent first)
func (pm *ProjectManager) ListProjects() []*Project {
	projects := make([]*Project, 0, len(pm.projects))
	for _, project := range pm.projects {
		projects = append(projects, project)
	}

	// Sort by last accessed time (most recent first)
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].LastAccessed.After(projects[j].LastAccessed)
	})

	return projects
}

// RemoveProject removes a project from the manager
func (pm *ProjectManager) RemoveProject(projectID string) error {
	project, exists := pm.projects[projectID]
	if !exists {
		return fmt.Errorf("project not found: %s", projectID)
	}

	// If this is the active project, clear active state
	if pm.activeProject != nil && pm.activeProject.ID == projectID {
		pm.activeProject = nil
		pm.storage.SetActiveProject("")
	}

	// Remove from memory
	delete(pm.projects, projectID)

	// Remove from storage
	if err := pm.storage.DeleteProject(projectID); err != nil {
		// Re-add to memory if storage delete failed
		pm.projects[projectID] = project
		return fmt.Errorf("failed to delete project from storage: %w", err)
	}

	// Save updated projects
	if err := pm.saveProjects(); err != nil {
		// Re-add to memory if save failed
		pm.projects[projectID] = project
		return fmt.Errorf("failed to save projects after deletion: %w", err)
	}

	return nil
}

// ValidateProjectPath checks if a path is valid for a new project
func (pm *ProjectManager) ValidateProjectPath(path string) error {
	if path == "" {
		return fmt.Errorf("project path cannot be empty")
	}

	cleanPath := filepath.Clean(path)
	if !filepath.IsAbs(cleanPath) {
		return fmt.Errorf("project path must be absolute: %s", path)
	}

	// Check if path exists
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return fmt.Errorf("project path does not exist: %s", cleanPath)
	}

	// Note: We allow existing projects to be re-added (they will be reactivated)
	return nil
}

// GetProjectInstances returns the instance IDs for a given project
func (pm *ProjectManager) GetProjectInstances(projectID string) ([]string, error) {
	project, exists := pm.projects[projectID]
	if !exists {
		return nil, fmt.Errorf("project not found: %s", projectID)
	}

	// Return a copy to prevent external modification
	instances := make([]string, len(project.Instances))
	copy(instances, project.Instances)
	return instances, nil
}

// AddInstanceToProject adds an instance to a project
func (pm *ProjectManager) AddInstanceToProject(projectID, instanceID string) error {
	project, exists := pm.projects[projectID]
	if !exists {
		return fmt.Errorf("project not found: %s", projectID)
	}

	project.AddInstance(instanceID)

	// Save to storage
	return pm.saveProjects()
}

// RemoveInstanceFromProject removes an instance from a project
func (pm *ProjectManager) RemoveInstanceFromProject(projectID, instanceID string) error {
	project, exists := pm.projects[projectID]
	if !exists {
		return fmt.Errorf("project not found: %s", projectID)
	}

	if !project.RemoveInstance(instanceID) {
		return fmt.Errorf("instance not found in project: %s", instanceID)
	}

	// Save to storage
	return pm.saveProjects()
}

// ProjectCount returns the total number of projects
func (pm *ProjectManager) ProjectCount() int {
	return len(pm.projects)
}

// loadProjects loads projects from storage
func (pm *ProjectManager) loadProjects() error {
	projectsJSON := pm.storage.GetProjects()
	if len(projectsJSON) == 0 {
		pm.projects = make(map[string]*Project) // Initialize empty map
		return nil                              // No projects to load
	}

	var projects map[string]*Project
	if err := json.Unmarshal(projectsJSON, &projects); err != nil {
		return fmt.Errorf("failed to unmarshal projects: %w", err)
	}

	// Validate loaded projects
	for id, project := range projects {
		if err := project.Validate(); err != nil {
			return fmt.Errorf("invalid project %s: %w", id, err)
		}
	}

	pm.projects = projects
	return nil
}

// saveProjects saves projects to storage
func (pm *ProjectManager) saveProjects() error {
	projectsJSON, err := json.Marshal(pm.projects)
	if err != nil {
		return fmt.Errorf("failed to marshal projects: %w", err)
	}

	return pm.storage.SaveProjects(projectsJSON)
}

// Project History Management Methods

// GetProjectHistory returns the project history
func (pm *ProjectManager) GetProjectHistory() *ProjectHistory {
	return pm.history
}

// UpdateProjectHistory adds a project path to the history and saves it
func (pm *ProjectManager) UpdateProjectHistory(projectPath string) error {
	if pm.history == nil {
		pm.history = NewProjectHistory()
	}

	pm.history.AddProject(projectPath)
	return pm.storage.SaveProjectHistory(pm.history)
}

// ClearProjectHistory removes all but the last N projects from history
func (pm *ProjectManager) ClearProjectHistory(keepLast int) error {
	if pm.history == nil {
		pm.history = NewProjectHistory()
		return nil
	}

	pm.history.ClearHistory(keepLast)
	return pm.storage.SaveProjectHistory(pm.history)
}

// GetRecentProjectPaths returns recent project paths for the UI
func (pm *ProjectManager) GetRecentProjectPaths() []string {
	if pm.history == nil {
		return []string{}
	}
	return pm.history.GetRecentProjects()
}

// GetTopProjectPaths returns the first N recent project paths (for quick select)
func (pm *ProjectManager) GetTopProjectPaths(count int) []string {
	if pm.history == nil {
		return []string{}
	}
	return pm.history.GetTopProjects(count)
}

// FilterProjectPaths returns project paths that match the search query
func (pm *ProjectManager) FilterProjectPaths(query string) []string {
	if pm.history == nil {
		return []string{}
	}
	return pm.history.FilterProjects(query)
}

// CleanupNonExistentProjects removes project paths that no longer exist
func (pm *ProjectManager) CleanupNonExistentProjects() error {
	if pm.history == nil {
		return nil
	}

	removedCount := pm.history.RemoveNonExistentPaths()
	if removedCount > 0 {
		return pm.storage.SaveProjectHistory(pm.history)
	}
	return nil
}
