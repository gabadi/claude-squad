# Validation Templates Task

## Purpose
Standardized templates for consistent validation across stories and validation types.

## Agent
**SM** (Scrum Master)

## Templates Included

## 1. User Journey Validation Template

### Pre-Validation Setup
```
Story: [Epic X.Story Y]
Validation Date: [Date]
QA Agent: [Name]
UX Collaboration: [Yes/No]
Time Allocated: [X hours]
Environment: [Test/Staging/Local]
```

### Journey Test Template
```
ACCEPTANCE CRITERIA MAPPING:
AC1: [Criteria text]
  - User Action Required: [Specific steps]
  - Expected Result: [What should happen]
  - Test Status: [Pass/Fail/Blocked]
  - Evidence: [Screenshot/Video link]

AC2: [Criteria text]
  - User Action Required: [Specific steps]
  - Expected Result: [What should happen]  
  - Test Status: [Pass/Fail/Blocked]
  - Evidence: [Screenshot/Video link]
```

### Usability Heuristics Checklist
```
□ Visibility of system status
□ Match between system and real world
□ User control and freedom
□ Consistency and standards
□ Error prevention
□ Recognition rather than recall
□ Flexibility and efficiency of use
□ Aesthetic and minimalist design
□ Help users recognize, diagnose, recover from errors
□ Help and documentation
```

### Journey Validation Report Template
```
## User Journey Validation Report

**Story**: [Epic X.Story Y]
**Date**: [Date]
**QA Agent**: [Name]
**Duration**: [X hours]

### Summary
- **Overall Status**: [Pass/Fail/Conditional]
- **Acceptance Criteria Met**: [X of Y]
- **Critical Issues**: [Number]
- **Usability Concerns**: [Number]

### Detailed Results
[For each AC - detailed test results]

### Issues Found
1. **Issue**: [Description]
   - **Severity**: [Critical/High/Medium/Low]
   - **Reproduction**: [Steps]
   - **Recommendation**: [Fix suggestion]

### Usability Observations
[UX-Expert collaboration notes]

### Evidence
[Links to screenshots, recordings, etc.]

### Recommendations
[Actionable feedback for Dev team]
```

## 2. Implementation Verification Template

### Pre-Verification Setup
```
Story: [Epic X.Story Y]
Verification Date: [Date]
Architect: [Name]
Story Type: [Simple/Complex]
Time Allocated: [X hours]
Code Review Scope: [Files/Components]
```

### Acceptance Criteria Mapping Template
```
AC1: [Criteria text]
  - Implementation Location: [File/Function/Component]
  - Verification Method: [Code review/Test execution/Demo]
  - Status: [Implemented/Partial/Missing]
  - Notes: [Technical details]

AC2: [Criteria text]
  - Implementation Location: [File/Function/Component]
  - Verification Method: [Code review/Test execution/Demo]
  - Status: [Implemented/Partial/Missing]
  - Notes: [Technical details]
```

### Code Quality Checklist
```
□ Follows project coding standards
□ Proper error handling
□ Performance considerations
□ Security best practices
□ Appropriate abstractions
□ Code maintainability
□ Documentation updated
□ Tests included/updated
```

### Architectural Review Checklist
```
□ Follows established patterns
□ No architectural violations
□ Dependencies properly managed
□ Integration points correct
□ Future extensibility considered
□ Technical debt minimized
□ Design decisions documented
```

### Implementation Verification Report Template
```
## Implementation Verification Report

**Story**: [Epic X.Story Y]
**Date**: [Date]
**Architect**: [Name]
**Duration**: [X hours]

### Summary
- **Overall Status**: [Approved/Conditional/Rejected]
- **Acceptance Criteria Implemented**: [X of Y]
- **Code Quality**: [Excellent/Good/Needs Improvement]
- **Architectural Alignment**: [Aligned/Minor Issues/Major Concerns]

### Acceptance Criteria Analysis
[For each AC - implementation details]

### Code Quality Assessment
[Detailed review notes]

### Architectural Review
[Alignment with patterns, concerns, recommendations]

### Recommendations
1. **Immediate Actions**: [Must fix before proceeding]
2. **Quality Improvements**: [Should address]
3. **Future Considerations**: [Could improve]

### Technical Debt Impact
[Assessment of any debt introduced]
```

## 3. Validation Failure Handling Template

### Escalation Matrix
```
ISSUE SEVERITY | HANDLER | TIMELINE | NEXT STEPS
Critical | PO + SM | Immediate | Back to story refinement
High | Architect + Dev | Same day | Direct collaboration fix
Medium | Original implementer | Next day | Standard fix cycle  
Low | Story backlog | Next sprint | Future improvement
```

### Retry Protocol Template
```
VALIDATION TYPE | MAX RETRIES | RETRY TRIGGER | ESCALATION
User Journey | 2 | Fix implementation | UX-Expert consult
Implementation | 3 | Address feedback | Architecture review
Quality Gates | 5 | Fix technical issues | SM intervention
```

## 4. Cross-Validation Consistency Template

### Validation Handoff Checklist
```
FROM DEV TO QA:
□ Implementation complete
□ Local testing passed
□ Documentation updated
□ Demo environment ready
□ Test data available

FROM QA TO ARCHITECT:
□ User journey validated
□ Issues documented
□ Evidence provided
□ Usability notes included

FROM ARCHITECT TO REVIEWS:
□ Implementation verified
□ Code quality confirmed
□ Architecture aligned
□ Technical debt assessed
```

## Usage Instructions

1. **SM** provides templates to QA and Architect at story start
2. **QA** uses journey validation template for user testing
3. **Architect** uses implementation template for code verification
4. **All agents** follow escalation matrix for issue handling
5. **Templates updated** based on retrospective feedback

## Continuous Improvement
- Monthly template review with team
- Update based on validation effectiveness
- Add new templates for emerging validation needs
- Version control template changes

## Notes
- Templates ensure consistency across different agents
- Reduce validation time through standardization
- Provide clear escalation paths for issues
- Enable metrics collection for process improvement