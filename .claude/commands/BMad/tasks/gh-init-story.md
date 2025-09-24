<!-- Powered by BMAD™ Core -->

# gh-init-story Task

## Purpose

Interactive GitHub story initialization with user confirmation prompts and configuration-based setup. Prevents unintended changes to GitHub and local environment.

## Configuration-Based Approach

### Story Configuration Template
```markdown
---
github_setup:
  epic: 0
  story: 0.2
  title: "Testing Infrastructure Foundation"
  status: "Planning"
  branch_name: "feature/story-0.2-testing-infrastructure-foundation"

issue_config:
  labels: ["story", "epic-0", "status-Planning"]
  assignee: "@me"

branch_config:
  create_branch: true
  switch_to_branch: false  # Safety: don't automatically switch
  push_to_remote: true
  set_upstream: true

safety_checks:
  require_confirmation: true
  dry_run: false
  prevent_force_push: true
```

## User Workflow

### Step 1: Configuration Loading
- Parse story file for GitHub setup configuration
- Extract epic, story, title, status
- Generate branch name if not provided
- Load default safety settings

### Step 2: Validation Checks
- ✅ GitHub CLI authenticated
- ✅ Story file exists and is valid
- ✅ No uncommitted changes (not required for issue/branch creation)
- ✅ Repository has proper remote configuration

### Step 3: Preview Changes
Show user exactly what will be created:
```
📋 GitHub Setup Preview for Story 0.2
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎯 Issue:
   Title: "Story {{epic}}.{{story}}: {{title}}"
   Labels: story, epic-0, status-Planning
   Assignee: @me
   Story File: [docs/stories/{{epic}}.{{story}}.story.md](docs/stories/{{epic}}.{{story}}.story.md)

🌿 Branch:
   Name: feature/story-{{epic}}.{{story}}-{{branch_name_suffix}}
   Create: Yes
   Switch: No (safety)
   Push: Yes
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Step 4: User Confirmation
```
❓ Confirm GitHub Setup
   [Y]es - Create these resources
   [N]o - Cancel operation
   [D]ry Run - Show commands without executing
   [C]ustomize - Modify configuration
```

### Step 5: Execute with Confirmation
- Only create issue after user confirmation
- Only create branch if explicitly confirmed
- Never automatically switch branches
- No automatic cleanup

## Safety Features

### 1. Confirmation Prompts
- Required before any GitHub API calls
- Required before any git operations
- Clear preview of all changes

### 2. Dry Run Mode
- Show all commands without executing
- Validate configuration syntax
- Test command templates

### 3. Change Prevention
- No automatic branch switching
- No automatic cleanup
- No force push operations
- Preserve working directory state

### 4. Error Handling
- Graceful failure with cleanup
- Preserve partial progress
- Clear error messages
- Recovery suggestions

## Usage Examples

### Basic Setup
```bash
# Interactive setup with confirmations
gh-init-story 0.2

# Dry run to preview
gh-init-story 0.2 --dry-run

# Custom configuration
gh-init-story 0.2 --config custom-config.yaml
```

### Advanced Usage
```bash
# Skip branch creation
gh-init-story 0.2 --no-branch

# Only issue creation
gh-init-story 0.2 --issue-only

# Custom branch name
gh-init-story 0.2 --branch "custom-branch-name"
```

## Story Data Extraction and Template Processing

### Automatic Story Data Extraction
The task automatically extracts story data from the story file and substitutes it into the issue template:

```bash
# Extract story data (simple approach)
STORY_FILE="docs/stories/0.2.story.md"
EPIC=$(grep "^epic:" "$STORY_FILE" | cut -d' ' -f2 | tr -d '"')
STORY=$(grep "^story:" "$STORY_FILE" | cut -d' ' -f2 | tr -d '"')
TITLE=$(grep "^title:" "$STORY_FILE" | cut -d' ' -f2- | tr -d '"')
STATUS=$(grep "^status:" "$STORY_FILE" | cut -d' ' -f2 | tr -d '"')

# Generate issue title with story number
ISSUE_TITLE="Story ${EPIC}.${STORY}: $TITLE"

# Generate branch name (simplified)
BRANCH_NAME="feature/story-${EPIC}.${STORY}-$(echo "$TITLE" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g' | sed 's/--*/-/g' | sed 's/^-//' | sed 's/-$//')"

# Use simple placeholders for story content
STORY_NARRATIVE="See story file for complete narrative"
ACCEPTANCE_CRITERIA="See story file for complete acceptance criteria"
```

### Template Processing
Replace placeholders in the issue template with actual story data:

```bash
# Process template with story data
BODY_CONTENT=$(cat .bmad-core/templates/issue-template.md | \
  sed "s/{{epic}}/$EPIC/g" | \
  sed "s/{{story}}/$STORY/g" | \
  sed "s/{{title}}/$TITLE/g" | \
  sed "s/{{status}}/$STATUS/g" | \
  sed "s/{{story_id}}/${EPIC}.${STORY}/g" | \
  sed "s|{{story_narrative}}|$STORY_NARRATIVE|g" | \
  sed "s|{{acceptance_criteria}}|$ACCEPTANCE_CRITERIA|g" | \
  sed "s|{{branch_name}}|$BRANCH_NAME|g" | \
  sed "s|{{dev_notes}}|See story file for development tasks|g" | \
  sed "s|{{technical_implementation}}|See story file for implementation details|g" | \
  sed "s|{{testing_requirements}}|See story file for testing requirements|g" | \
  sed "s|{{change_log}}|Initial issue creation|g" | \
  sed "s|{{test_info}}|Story file contains complete testing information|g")
```

### Issue Creation with Processed Data
```bash
gh issue create \
  --title "$ISSUE_TITLE" \
  --body "$BODY_CONTENT" \
  --label "story,epic-${EPIC},status-${STATUS}" \
  --assignee "@me"
```

### Branch Creation and Linking
```bash
# Create and push branch
git checkout -b $BRANCH_NAME
git push -u origin $BRANCH_NAME

# Link branch to issue (if issue was created)
ISSUE_NUMBER=$(gh issue list --limit 1 --json number --jq '.[0].number')
if [ -n "$ISSUE_NUMBER" ]; then
  gh issue edit $ISSUE_NUMBER --add-label "branch-$BRANCH_NAME"
  echo "Branch linked to issue #$ISSUE_NUMBER"
fi
```

## Error Recovery

### Common Issues
1. **Authentication**: `gh auth login`
2. **Permissions**: Check repo write access
3. **Network**: Verify internet connection
4. **Git State**: Uncommitted changes are acceptable for issue/branch creation

### Rollback Commands
```bash
# Delete issue (if created)
gh issue close <issue-number>

# Delete branch (if created)
git branch -D <branch-name>
git push origin --delete <branch-name>
```

## Integration Points

### 1. BMad Agent Integration
- Called from `develop-story` workflow
- Passes confirmation status back to agent
- Respects agent-level safety settings

### 2. Master Agent Integration
- Available as standalone command
- Supports batch operations with confirmations
- Provides detailed execution reports

### 3. CI/CD Integration
- Can be used in automated workflows
- Supports non-interactive mode with pre-approval
- Provides execution logs and status reporting

---

*This task provides safe, interactive GitHub story initialization with full user control*