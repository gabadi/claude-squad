# Multi-Project Architecture Specification

## System Overview

### Current Architecture Summary

Claude Squad is a terminal-based project management tool built with Go and Bubble Tea TUI framework. The current architecture manages single-project instances with the following core components:

- **Application Layer** (`app/`): Main application state and UI coordination using Bubble Tea
- **Session Management** (`session/`): Instance lifecycle, tmux session handling, and git worktree management
- **UI Components** (`ui/`): Bubble Tea components for list, preview, diff, console, and overlays
- **Configuration** (`config/`): Application configuration and persistent state management
- **Storage** (`session/storage.go`): JSON-based persistence for instances and application state

### Multi-Project Feature Scope

Transform Claude Squad from single-project instance management to hierarchical multi-project support:

1. **Hierarchical Project-Instance Structure**: Projects contain multiple instances
2. **Smart Project Discovery**: Auto-detect and add projects via intelligent path resolution
3. **Project Context Switching**: Quick navigation between project contexts
4. **Cross-Project Instance Management**: Maintain existing instance operations within project boundaries
5. **Enhanced Storage**: Extend current storage to support project grouping

### Integration Points and Dependencies

- **Bubble Tea UI Framework**: Hierarchical list rendering and input handling
- **Git Worktree System**: Per-project git operations and branch management  
- **Tmux Session Management**: Console and CLI session isolation per project
- **File System Operations**: Project path resolution and validation
- **JSON Configuration**: Extended storage schema for projects and instances

## Implementation Specification

### Data Model Transformations

#### 1. Project Structure
```go
// project/project.go
type Project struct {
    ID           string    `json:"id"`           // Unique identifier
    Name         string    `json:"name"`         // Display name
    Path         string    `json:"path"`         // Absolute filesystem path
    LastAccessed time.Time `json:"last_accessed"`
    CreatedAt    time.Time `json:"created_at"`
    IsActive     bool      `json:"is_active"`    // Current project context
    Instances    []string  `json:"instances"`    // Instance IDs belonging to this project
}

type ProjectManager struct {
    projects      map[string]*Project
    activeProject *Project
    storage       ProjectStorage
}
```

#### 2. Enhanced Instance Structure
```go
// session/instance.go - Add project context
type Instance struct {
    // Existing fields...
    ProjectID string `json:"project_id"` // Reference to parent project
    
    // Existing implementation remains unchanged
}

// session/storage.go - Add project context to InstanceData
type InstanceData struct {
    // Existing fields...
    ProjectID string `json:"project_id"`
}
```

#### 3. Storage Schema Extensions
```go
// config/state.go - Add project storage
type State struct {
    HelpScreensSeen uint32          `json:"help_screens_seen"`
    InstancesData   json.RawMessage `json:"instances"`
    ProjectsData    json.RawMessage `json:"projects"`     // New field
    ActiveProject   string          `json:"active_project"` // New field
}

// project/storage.go
type ProjectStorage interface {
    SaveProjects(projectsJSON json.RawMessage) error
    GetProjects() json.RawMessage
    DeleteProject(projectID string) error
    SetActiveProject(projectID string) error
    GetActiveProject() string
}
```

### Code Structure Changes Required

#### 1. New Package: `project/`
```
project/
├── project.go          # Project struct and core operations
├── manager.go          # ProjectManager for multi-project operations
├── storage.go          # Project storage interface and implementation
├── discovery.go        # Smart project path resolution and discovery
└── validation.go       # Project path validation and conflict resolution
```

#### 2. Enhanced UI Components
```
ui/
├── project_list.go     # Hierarchical project-instance list component
├── project_input.go    # Smart project input overlay with autocomplete
├── project_switcher.go # Quick project context switcher overlay
```

#### 3. Updated Application State
```go
// app/app.go - Enhanced home struct
type home struct {
    // Existing fields...
    projectManager    *project.ProjectManager
    projectInputOverlay *ui.ProjectInputOverlay
    projectSwitchOverlay *ui.ProjectSwitchOverlay
}

// New application states
const (
    stateAddProject state = iota + 5  // After existing states
    stateSwitchProject
)
```

### Key Functions and Structs to Implement

#### 1. Project Management Core
```go
// project/manager.go
func NewProjectManager(storage ProjectStorage) (*ProjectManager, error)
func (pm *ProjectManager) AddProject(path, name string) (*Project, error)
func (pm *ProjectManager) SetActiveProject(projectID string) error
func (pm *ProjectManager) GetActiveProject() *Project
func (pm *ProjectManager) DiscoverSiblingProjects(currentPath string) ([]*Project, error)
func (pm *ProjectManager) ValidateProjectPath(path string) error
func (pm *ProjectManager) GetProjectInstances(projectID string) ([]*session.Instance, error)

// project/discovery.go
func ResolvePath(input, currentProjectPath string) (string, error)
func FindSiblingDirectories(basePath string) ([]string, error)
func AutocompleteProjectPath(input, context string) ([]string, error)
```

#### 2. Enhanced UI Components
```go
// ui/project_list.go
type ProjectList struct {
    projects        []*project.Project
    instances       map[string][]*session.Instance
    selectedProject int
    selectedInstance int
    isProjectMode   bool
}

func (pl *ProjectList) RenderHierarchical() string
func (pl *ProjectList) NavigateProjects() 
func (pl *ProjectList) NavigateInstances()
func (pl *ProjectList) ToggleProjectCollapse(projectID string)

// ui/project_input.go
type ProjectInputOverlay struct {
    textInput     textinput.Model
    suggestions   []string
    currentPath   string
}

func (pio *ProjectInputOverlay) HandleInput(msg tea.KeyMsg) (bool, error)
func (pio *ProjectInputOverlay) GenerateSuggestions(input string) []string
func (pio *ProjectInputOverlay) ValidateAndResolve(input string) (string, error)

// ui/project_switcher.go
type ProjectSwitchOverlay struct {
    projects      []*project.Project
    filteredList  []*project.Project
    selectedIndex int
    filterInput   string
}

func (pso *ProjectSwitchOverlay) FilterProjects(input string)
func (pso *ProjectSwitchOverlay) HandleNavigation(msg tea.KeyMsg) bool
```

#### 3. Application Integration
```go
// app/app.go - New key handlers
func (m *home) handleAddProject(msg tea.KeyMsg) (tea.Model, tea.Cmd)
func (m *home) handleSwitchProject(msg tea.KeyMsg) (tea.Model, tea.Cmd)
func (m *home) handleProjectInput(msg tea.KeyMsg) (tea.Model, tea.Cmd)
func (m *home) updateProjectContext(projectID string) tea.Cmd

// Enhanced initialization
func newHome(ctx context.Context, program string, autoYes bool) *home {
    // Existing initialization...
    
    // Initialize project manager
    projectManager, err := project.NewProjectManager(appState)
    if err != nil {
        fmt.Printf("Failed to initialize project manager: %v\n", err)
        os.Exit(1)
    }
    
    h.projectManager = projectManager
    
    // Load instances within project context
    // Modified to filter by active project
}
```

### Storage and Configuration Updates

#### 1. Enhanced State Management
```go
// config/state.go - Extended State struct
func (s *State) SaveProjects(projectsJSON json.RawMessage) error {
    s.ProjectsData = projectsJSON
    return SaveState(s)
}

func (s *State) GetProjects() json.RawMessage {
    if len(s.ProjectsData) == 0 {
        return json.RawMessage("[]")
    }
    return s.ProjectsData
}

func (s *State) SetActiveProject(projectID string) error {
    s.ActiveProject = projectID
    return SaveState(s)
}

func (s *State) GetActiveProject() string {
    return s.ActiveProject
}
```

#### 2. Migration Strategy
```go
// config/migration.go
func MigrateToMultiProject(state *State) error {
    // Create default project from current working directory
    // Move existing instances to default project
    // Set default project as active
}
```

## Development Phases

### Phase 1: Core Data Model (Days 1-3)
**Deliverables:**
- [ ] `project/project.go` - Project struct and basic operations
- [ ] `project/storage.go` - Storage interface implementation
- [ ] `config/state.go` - Extended state with project fields
- [ ] Migration logic for existing instances

**Files to Create/Modify:**
- `project/project.go` (new)
- `project/storage.go` (new)
- `config/state.go` (modify - add ProjectsData, ActiveProject fields)
- `session/instance.go` (modify - add ProjectID field)
- `session/storage.go` (modify - add ProjectID to InstanceData)

**Testing Checkpoints:**
- Project creation and storage persistence
- Instance-project association
- Migration from single-project to multi-project state

### Phase 2: Project Management Logic (Days 4-6)
**Deliverables:**
- [ ] `project/manager.go` - ProjectManager implementation
- [ ] `project/discovery.go` - Smart path resolution
- [ ] `project/validation.go` - Path validation and conflict resolution
- [ ] Unit tests for project operations

**Files to Create/Modify:**
- `project/manager.go` (new)
- `project/discovery.go` (new)
- `project/validation.go` (new)
- Test files for project package

**Testing Checkpoints:**
- Project discovery from sibling directories
- Path resolution (relative, absolute, name-only)
- Project validation and error handling

### Phase 3: UI Component Integration (Days 7-10)
**Deliverables:**
- [ ] `ui/project_list.go` - Hierarchical project-instance list
- [ ] Enhanced navigation in main application
- [ ] Project context display in UI
- [ ] Updated list rendering with project grouping

**Files to Create/Modify:**
- `ui/project_list.go` (new)
- `ui/list.go` (modify - integrate project context)
- `app/app.go` (modify - integrate project manager)

**Testing Checkpoints:**
- Hierarchical display rendering correctly
- Navigation between projects and instances
- Project context updates reflected in UI

### Phase 4: Input Overlays and Smart Discovery (Days 11-13)
**Deliverables:**
- [ ] `ui/project_input.go` - Smart project input with autocomplete
- [ ] `ui/project_switcher.go` - Quick project switcher
- [ ] Enhanced keyboard handling for new overlays
- [ ] Real-time path suggestions and validation

**Files to Create/Modify:**
- `ui/project_input.go` (new)
- `ui/project_switcher.go` (new)
- `app/app.go` (modify - add new state handlers)
- `keys/keys.go` (modify - add project keys)

**Testing Checkpoints:**
- Smart input suggestions working
- Project switcher filtering and navigation
- Keyboard shortcuts properly integrated

### Phase 5: Integration and Polish (Days 14-16)
**Deliverables:**
- [ ] Complete keyboard shortcut integration
- [ ] Error handling and edge cases
- [ ] Help system updates
- [ ] Performance optimization
- [ ] Full integration testing

**Files to Create/Modify:**
- `keys/keys.go` (modify - add `p` and `ctrl+p` mappings)
- `app/help.go` (modify - update help text)
- All files - error handling and edge case fixes

**Testing Checkpoints:**
- All keyboard shortcuts working
- Error scenarios handled gracefully
- Performance acceptable with multiple projects
- Help system reflects new features

## Technical Patterns

### Go Code Patterns and Interfaces

#### 1. Project Interface Pattern
```go
type ProjectOperations interface {
    Create(path, name string) (*Project, error)
    Update(project *Project) error
    Delete(id string) error
    GetByID(id string) (*Project, error)
    List() ([]*Project, error)
}

type ProjectDiscovery interface {
    ResolvePath(input, context string) (string, error)
    FindSiblings(basePath string) ([]string, error)
    ValidatePath(path string) error
}
```

#### 2. Storage Pattern Extensions
```go
// Extend existing storage pattern to support projects
type Storage struct {
    state config.StateManager  // Existing interface enhanced
}

func (s *Storage) SaveProjects(projects []*project.Project) error
func (s *Storage) LoadProjects() ([]*project.Project, error)
func (s *Storage) LoadInstancesForProject(projectID string) ([]*Instance, error)
```

### Bubble Tea UI Patterns for Hierarchical Components

#### 1. Hierarchical List Pattern
```go
type HierarchicalItem interface {
    Render(selected bool, level int) string
    IsCollapsible() bool
    GetChildren() []HierarchicalItem
    GetType() ItemType
}

type ProjectItem struct {
    project   *project.Project
    instances []InstanceItem
    collapsed bool
}

func (pi ProjectItem) Render(selected bool, level int) string {
    // Render with appropriate indentation and collapse indicators
}
```

#### 2. Multi-State Overlay Pattern
```go
type OverlayManager struct {
    activeOverlay Overlay
    overlayStack  []Overlay
}

type Overlay interface {
    Render() string
    HandleKeyPress(msg tea.KeyMsg) (shouldClose bool, err error)
    SetSize(width, height int)
}
```

### Storage Patterns and Migration Strategies

#### 1. Schema Versioning
```go
type StateVersion int

const (
    StateV1 StateVersion = 1  // Original single-project
    StateV2 StateVersion = 2  // Multi-project support
)

type VersionedState struct {
    Version StateVersion `json:"version"`
    Data    interface{}  `json:"data"`
}

func MigrateState(from StateVersion, data interface{}) (*State, error)
```

#### 2. Backward Compatibility
```go
func LoadStateWithMigration() *State {
    // Attempt to load as current version
    // If fails, try previous versions and migrate
    // Maintain data integrity throughout process
}
```

### Error Handling and Edge Case Management

#### 1. Project Path Conflicts
```go
func (pm *ProjectManager) ResolvePathConflict(path string) error {
    // Handle duplicate project paths
    // Provide resolution strategies
    // Maintain data consistency
}
```

#### 2. Instance Orphaning Prevention
```go
func (pm *ProjectManager) ValidateInstanceProjectAssociation() error {
    // Ensure all instances belong to valid projects
    // Handle orphaned instances
    // Provide repair mechanisms
}
```

## Implementation Guide

### Step-by-Step Development Instructions

#### Step 1: Set up Project Package Structure
```bash
mkdir project
cd project
touch project.go manager.go storage.go discovery.go validation.go
```

#### Step 2: Implement Core Project Struct
```go
// project/project.go
package project

import (
    "path/filepath"
    "time"
)

type Project struct {
    ID           string    `json:"id"`
    Name         string    `json:"name"`
    Path         string    `json:"path"`
    LastAccessed time.Time `json:"last_accessed"`
    CreatedAt    time.Time `json:"created_at"`
    IsActive     bool      `json:"is_active"`
    Instances    []string  `json:"instances"`
}

func NewProject(path, name string) (*Project, error) {
    absPath, err := filepath.Abs(path)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
    }
    
    return &Project{
        ID:           generateProjectID(absPath),
        Name:         name,
        Path:         absPath,
        CreatedAt:    time.Now(),
        LastAccessed: time.Now(),
        Instances:    make([]string, 0),
    }, nil
}

func generateProjectID(path string) string {
    // Generate unique ID based on path hash
    return fmt.Sprintf("proj_%x", sha256.Sum256([]byte(path)))[:12]
}
```

#### Step 3: Implement ProjectManager
```go
// project/manager.go
package project

type ProjectManager struct {
    projects      map[string]*Project
    activeProject *Project
    storage       ProjectStorage
}

func NewProjectManager(storage ProjectStorage) (*ProjectManager, error) {
    pm := &ProjectManager{
        projects: make(map[string]*Project),
        storage:  storage,
    }
    
    if err := pm.loadProjects(); err != nil {
        return nil, fmt.Errorf("failed to load projects: %w", err)
    }
    
    return pm, nil
}

func (pm *ProjectManager) AddProject(path, name string) (*Project, error) {
    // Validate path
    if err := pm.validateProjectPath(path); err != nil {
        return nil, err
    }
    
    // Create project
    project, err := NewProject(path, name)
    if err != nil {
        return nil, err
    }
    
    // Add to manager
    pm.projects[project.ID] = project
    
    // Save to storage
    if err := pm.saveProjects(); err != nil {
        return nil, err
    }
    
    return project, nil
}
```

#### Step 4: Extend Application State
```go
// config/state.go - Add to State struct
type State struct {
    HelpScreensSeen uint32          `json:"help_screens_seen"`
    InstancesData   json.RawMessage `json:"instances"`
    ProjectsData    json.RawMessage `json:"projects"`
    ActiveProject   string          `json:"active_project"`
}

// Add project storage methods
func (s *State) SaveProjects(projectsJSON json.RawMessage) error {
    s.ProjectsData = projectsJSON
    return SaveState(s)
}

func (s *State) GetProjects() json.RawMessage {
    if len(s.ProjectsData) == 0 {
        return json.RawMessage("[]")
    }
    return s.ProjectsData
}
```

#### Step 5: Integrate with Main Application
```go
// app/app.go - Add to home struct
type home struct {
    // Existing fields...
    projectManager       *project.ProjectManager
    projectInputOverlay  *ui.ProjectInputOverlay
    projectSwitchOverlay *ui.ProjectSwitchOverlay
}

// Add new states
const (
    stateAddProject state = iota + 5
    stateSwitchProject
)

// Modify newHome function
func newHome(ctx context.Context, program string, autoYes bool) *home {
    // Existing initialization...
    
    // Initialize project manager
    projectManager, err := project.NewProjectManager(appState)
    if err != nil {
        fmt.Printf("Failed to initialize project manager: %v\n", err)
        os.Exit(1)
    }
    
    h.projectManager = projectManager
    
    // Rest of initialization...
}
```

#### Step 6: Add Keyboard Shortcuts
```go
// keys/keys.go - Add new key constants
const (
    // Existing keys...
    KeyAddProject
    KeySwitchProject
)

// Add to GlobalKeyStringsMap
var GlobalKeyStringsMap = map[string]KeyName{
    // Existing mappings...
    "p":        KeyAddProject,
    "ctrl+p":   KeySwitchProject,
}

// Add to GlobalkeyBindings
var GlobalkeyBindings = map[KeyName]key.Binding{
    // Existing bindings...
    KeyAddProject: key.NewBinding(
        key.WithKeys("p"),
        key.WithHelp("p", "add project"),
    ),
    KeySwitchProject: key.NewBinding(
        key.WithKeys("ctrl+p"),
        key.WithHelp("ctrl+p", "switch project"),
    ),
}
```

#### Step 7: Implement UI Components
```go
// ui/project_list.go
package ui

type ProjectList struct {
    projects        []*project.Project
    instances       map[string][]*session.Instance
    selectedProject int
    selectedInstance int
    inProjectMode   bool
    renderer        *ProjectRenderer
}

func NewProjectList(projects []*project.Project) *ProjectList {
    return &ProjectList{
        projects:  projects,
        instances: make(map[string][]*session.Instance),
        renderer:  &ProjectRenderer{},
    }
}

func (pl *ProjectList) RenderHierarchical() string {
    var b strings.Builder
    
    for i, proj := range pl.projects {
        // Render project header
        b.WriteString(pl.renderer.RenderProject(proj, i == pl.selectedProject))
        
        // Render instances if not collapsed
        if instances, ok := pl.instances[proj.ID]; ok {
            for j, inst := range instances {
                selected := i == pl.selectedProject && j == pl.selectedInstance && !pl.inProjectMode
                b.WriteString(pl.renderer.RenderInstance(inst, j, selected, true))
            }
        }
    }
    
    return b.String()
}
```

### Integration Points with Existing Codebase

#### 1. Instance Creation Integration
```go
// Modify existing instance creation to associate with active project
func (m *home) createNewInstance(title, program string) (*session.Instance, error) {
    activeProject := m.projectManager.GetActiveProject()
    if activeProject == nil {
        return nil, fmt.Errorf("no active project selected")
    }
    
    instance, err := session.NewInstance(session.InstanceOptions{
        Title:     title,
        Path:      activeProject.Path,  // Use project path
        Program:   program,
        ProjectID: activeProject.ID,    // Associate with project
    })
    
    return instance, err
}
```

#### 2. Storage Integration
```go
// Modify storage to filter instances by project
func (s *Storage) LoadInstancesForProject(projectID string) ([]*Instance, error) {
    allInstances, err := s.LoadInstances()
    if err != nil {
        return nil, err
    }
    
    projectInstances := make([]*Instance, 0)
    for _, instance := range allInstances {
        if instance.ProjectID == projectID {
            projectInstances = append(projectInstances, instance)
        }
    }
    
    return projectInstances, nil
}
```

### Performance and Scalability Considerations

#### 1. Lazy Loading
```go
// Load project instances only when needed
func (pm *ProjectManager) GetProjectInstances(projectID string) ([]*session.Instance, error) {
    if instances, cached := pm.instanceCache[projectID]; cached {
        return instances, nil
    }
    
    instances, err := pm.storage.LoadInstancesForProject(projectID)
    if err != nil {
        return nil, err
    }
    
    pm.instanceCache[projectID] = instances
    return instances, nil
}
```

#### 2. Efficient Rendering
```go
// Render only visible items in hierarchical list
func (pl *ProjectList) RenderVisible(startIdx, count int) string {
    // Calculate visible range based on scroll position
    // Render only items within viewport
    // Maintain smooth scrolling performance
}
```

This comprehensive architecture specification provides the complete technical foundation for implementing multi-project support in Claude Squad. The specification is structured to enable systematic development while maintaining compatibility with the existing codebase and ensuring optimal user experience.