---
epic: 0
title: "Demo Infrastructure Foundation"
status: "Planning"
---

# Epic 0: Demo Infrastructure Foundation

## Epic Description

Establish comprehensive demo infrastructure to ensure every future story and PR includes automated testing, video recordings with visual input display, screenshots, and embedded demos. This foundational epic creates the tooling and workflows that make all subsequent development inherently demoable with professional video demonstrations.

**Critical Requirement**: No new feature development (Stories 1.3+) begins until Epic 0 is complete.

## Epic Goals

1. **GitHub Workflow Foundation**: Automated issue creation and branch management for stories
2. **Complete Testing Infrastructure**: Unit, integration, e2e, and UI testing with visual validation
3. **Video Recording System**: Automated demo videos with keystroke/mouse visualization
4. **PR Demo Automation**: Every PR automatically includes embedded demo videos and screenshots
5. **Professional Demo Quality**: Visual input overlays, progress indicators, and AC validation

## Stories (Logical Progression Order)

- **Story 0.1**: BMad Agent Development Workflow Automation - *Draft* (BMad agent uses gh commands to create GitHub issues and branches)
- **Story 0.2**: Testing Infrastructure Foundation - *Planning* (Expanded scope: Unit, integration, e2e testing + video recording setup)
- **Story 0.3**: Video Recording & Demo System - *Planning* (Visual input overlays, keystroke/mouse visualization, demo video generation)
- **Story 0.4**: PR Demo Automation & Integration - *Planning* (Automated PR embedding, demo video integration, screenshot comparison)
- **Story 0.5**: Documentation & Template Updates - *Planning* (Template integration, developer workflow documentation)

## Epic Definition of Done

- [ ] GitHub automatically creates issues and branches for new stories
- [ ] Complete testing pipeline: unit → integration → e2e → UI tests → demo video
- [ ] Video recordings include visual keystroke/mouse overlays
- [ ] Every PR automatically embeds demo videos and screenshots
- [ ] Bug fixes include before/after demo videos
- [ ] Screenshot comparison for UI regression testing
- [ ] All future stories inherit demo infrastructure automatically
- [ ] Developer workflow requires zero manual demo steps

## Demo Video Requirements

**Visual Input Display**:
- Floating overlays showing pressed keys ("j", "k", "space", "enter", etc.)
- Mouse click animations with circle indicators
- Input timeline showing sequence of actions
- Step progress indicators ("Step 3/7")

**Recording Quality**:
- Terminal session recording with timing data
- Clear visibility of application state changes
- Before/after comparisons for bug fixes
- Acceptance criteria validation markers (✅/❌)

**Output Formats**:
- MP4 videos for GitHub PR embedding
- GIF animations for README documentation
- Screenshots for quick reference
- Test coverage reports with visual proof

## Epic Success Metrics

- **100% automation**: Zero manual steps for demo creation
- **Speed**: < 2 minutes additional CI time for demo generation
- **Quality**: Professional video quality with input visualization
- **Coverage**: Every PR includes visual proof of functionality
- **Developer Experience**: Seamless integration with existing workflow

## Epic Dependencies

- GitHub Actions and workflow automation
- Go testing frameworks (testify, ginkgo, gomega)
- Terminal recording (asciinema, ttyrec, script)
- Video processing (ffmpeg, key-mon, screenkey)
- TUI testing frameworks for Go applications
- Screenshot comparison tools

---

*Epic 0 Priority: MUST complete before any new feature development*
*Foundation for demo-first development ensuring all future work is visually validated*