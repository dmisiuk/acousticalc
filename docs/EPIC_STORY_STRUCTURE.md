# AcoustiCalc Epic/Story Structure

## Overview

AcoustiCalc has been restructured to follow Epic/Story methodology instead of traditional GitHub milestone/issue structure. This provides better user-centric focus and clearer value delivery.

## Structure Organization

### Documents Location
- **Epics**: `docs/epics/` - High-level feature groupings
- **Stories**: `docs/stories/` - Individual user stories
- **PRD**: `docs/prd/` - Sharded Product Requirements (updated for Epic/Story)
- **Architecture**: `docs/architecture/` - Sharded technical documentation

### Epic/Story Mapping

#### Epic 1: Foundation & Core Engine
- **Story 1.1**: Project Setup ✅ Done
- **Story 1.2**: Core Calculator Engine ✅ Done
- **Story 1.3**: TUI Framework Integration (Pending)
- **Story 1.4**: MVP Release (Pending)

#### Epic 2: Enhanced User Experience
- **Story 2.1**: Audio System
- **Story 2.2**: Advanced Calculator Features
- **Story 2.3**: User Experience Polish
- **Story 2.4**: v1.0 Release

#### Epic 3: Distribution & Community
- **Story 3.1**: Package Management
- **Story 3.2**: Community Building
- **Story 3.3**: v1.1 Community Release

## Key Changes from Milestone Structure

### Before (Milestones)
- Task-focused: "Milestone 1.1: Project Setup"
- Technical orientation: Lists of implementation tasks
- No clear user value statement

### After (Epic/Story)
- User-focused: "As a user, I want... So that..."
- Value-driven: Clear benefit statements for each story
- Epic cohesion: Related stories grouped by user value theme

## B-Mad Agent Commands

Use these commands for Epic/Story management:

```bash
# Create new epics
/brownfield-create-epic [epic description]

# Create stories from requirements
/create-next-story [based on current epic]

# Review existing stories
/review-story [story reference]

# Create brownfield stories for enhancements
/create-brownfield-story [enhancement description]
```

## Story Format

All stories follow this structure:

```markdown
---
epic: [number]
story: [number]
title: "[Story Title]"
status: "[Draft|In Progress|Review|Done]"
---

### Story
**As a** [user type]
**I want** [goal/functionality]
**So that** [benefit/value]

### Acceptance Criteria
1. [Specific, testable criteria]
2. [Measurable outcomes]

### Story Tasks
- [ ] [Implementation task]
- [ ] [Testing requirement]
- [ ] [Documentation update]
```

## Epic Format

```markdown
---
epic: [number]
title: "[Epic Title]"
status: "[Planning|In Progress|Completed]"
---

### Epic
[Epic description and context]

### Epic Goals
[What this epic accomplishes]

### Stories
- **Story X.1**: [Title] - [Status]
- **Story X.2**: [Title] - [Status]

### Epic Definition of Done
- [ ] [Epic-level completion criteria]
- [ ] [Value delivery confirmation]
```

## Benefits of This Structure

1. **User-Centric**: Every story delivers user value
2. **Clear Communication**: Easy for stakeholders to understand
3. **Agile Compatibility**: Natural fit with agile development
4. **Prioritization**: Easier to reorder based on user feedback
5. **B-Mad Integration**: Works seamlessly with B-Mad methodology

## Migration Notes

- Existing milestone content has been converted to user stories
- Technical tasks are now story tasks within user-focused stories
- Epic structure groups related user value themes
- B-Mad agents are configured to work with new structure

---

*Updated: September 23, 2025*
*Structure Version: v1.0*