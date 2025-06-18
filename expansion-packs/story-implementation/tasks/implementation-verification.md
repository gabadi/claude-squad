# Implementation Verification Task

## Purpose
Architect verifies that implementation meets acceptance criteria functionality, code quality, and architectural alignment.

## Agent
**Architect**

## When to Execute
After story implementation is complete, can run parallel to user journey validation.

## Time Box
**Maximum 2 hours** for comprehensive verification
**Maximum 1 hour** for simple story verification

## Inputs Required
- Story file with acceptance criteria
- Architectural checklist
- Implementation code/artifacts
- Technical design requirements

## Verification Criteria

### 1. Functional Completeness
- [ ] Implementation fulfills all acceptance criteria requirements
- [ ] Code produces expected behavior described in story
- [ ] Edge cases handled appropriately
- [ ] Integration points working correctly

### 2. Code Quality
- [ ] Code follows project coding standards
- [ ] Proper error handling implemented
- [ ] Performance considerations addressed
- [ ] Security best practices followed

### 3. Architectural Alignment
- [ ] Implementation follows established patterns
- [ ] No architectural violations introduced
- [ ] Dependencies managed appropriately
- [ ] Future extensibility considered

### 4. Technical Debt Prevention
- [ ] No quick hacks or temporary solutions
- [ ] Proper abstractions used
- [ ] Code maintainability preserved
- [ ] Documentation updated where needed

## Evidence Required
- [ ] Code review notes with specific feedback
- [ ] Verification that each acceptance criteria maps to working code
- [ ] Documentation of architectural decisions made
- [ ] Performance or security impact assessment

## Feedback Type
**Actionable recommendations** rather than pass/fail decisions:
- Specific code improvements needed
- Architectural guidance for alignment
- Technical debt prevention suggestions
- Performance optimization opportunities

## Failure Conditions
- Acceptance criteria not met by implementation
- Significant architectural violations
- Critical technical debt introduced
- Security or performance issues identified

## Success Criteria
- All acceptance criteria verifiably implemented
- Code quality meets project standards
- Architectural consistency maintained
- Technical debt minimized

## Escalation Process
- **Minor issues**: Direct feedback to Dev for quick fixes
- **Architectural concerns**: Work directly with Dev team for solution
- **Major design issues**: Back to story analysis with PO involvement

## Output
Implementation Verification Report including:
- Acceptance criteria mapping to implementation
- Code quality assessment
- Architectural alignment notes
- Recommendations for improvement (if any)
- Approval status with justification

## Checklist Template

### For Simple Stories
- [ ] Change implements acceptance criteria
- [ ] Code quality acceptable
- [ ] No architectural violations
- [ ] Minimal technical debt impact

### For Complex Stories
- [ ] All acceptance criteria implemented
- [ ] Integration points verified
- [ ] Performance impact assessed
- [ ] Security considerations reviewed
- [ ] Architectural patterns followed
- [ ] Future maintainability preserved

## Notes
- Focus on technical debt prevention and architectural consistency
- Provide specific, actionable feedback
- Time-box to prevent workflow blocking
- Collaborate with Dev team for solution-oriented approach