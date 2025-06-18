# User Journey Validation Task

## Purpose
QA validates that users can complete the workflow described in acceptance criteria with basic usability heuristics.

## Agent
**QA** (with UX-Expert collaboration)

## When to Execute
After story implementation is complete, before formal reviews begin.

## Time Box
**Maximum 4 hours** for complex journeys

## Inputs Required
- Story file with acceptance criteria
- UX validation template
- Implementation artifacts (screenshots, demos)

## Validation Criteria

### 1. Functional Journey Completion
- [ ] User can complete primary workflow described in acceptance criteria
- [ ] All acceptance criteria steps are achievable by end user
- [ ] Error scenarios handle gracefully
- [ ] Edge cases identified in story work correctly

### 2. Basic Usability Heuristics (UX-Expert collaboration)
- [ ] User interface is intuitive for target user
- [ ] Navigation flows logically
- [ ] Feedback provided for user actions
- [ ] Error messages are clear and actionable

### 3. Cross-Device/Browser Validation (if applicable)
- [ ] Functionality works on primary supported devices
- [ ] Core workflow accessible across main browsers
- [ ] Responsive design maintains usability

### 4. Accessibility Basic Checks
- [ ] Keyboard navigation possible for core workflow
- [ ] Screen reader compatibility for key elements
- [ ] Color contrast adequate for text/background

## Evidence Required
- [ ] Screenshots of successful user journey completion
- [ ] Documentation of any issues found
- [ ] Verification that acceptance criteria are met from user perspective
- [ ] Notes on usability observations

## Failure Conditions
- User cannot complete acceptance criteria workflow
- Critical usability issues prevent effective task completion
- Accessibility barriers block primary user scenarios
- Implementation doesn't match story description

## Success Criteria
- Complete user workflow demonstrated
- Basic usability heuristics validated
- Evidence provided for journey completion
- Any issues documented with clear reproduction steps

## Escalation Process
- **Minor issues**: Feedback to Dev for quick fixes
- **Major usability problems**: Escalate to UX-Expert for design consultation
- **Functionality gaps**: Back to story refinement with PO

## Output
User Journey Validation Report including:
- Journey completion evidence
- Usability observations
- Issues found (if any)
- Recommendations for improvement

## Notes
- Focus on acceptance criteria fulfillment + critical user paths
- Collaborate with UX-Expert for usability criteria definition
- Time-box validation to prevent blocking workflow
- Provide actionable feedback, not just pass/fail decisions