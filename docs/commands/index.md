---
title: Command Reference
---

# Command Reference

This section contains auto-generated documentation for all Canvas CLI commands.

## Available Commands

| Command | Description |
|---------|-------------|
| [canvas](canvas.md) | Root command |
| [canvas accounts](canvas_accounts.md) | Account management |
| [canvas admins](canvas_admins.md) | Account administrator management |
| [canvas analytics](canvas_analytics.md) | Canvas analytics |
| [canvas announcements](canvas_announcements.md) | Announcement management |
| [canvas api](canvas_api.md) | Raw API requests |
| [canvas assignment-groups](canvas_assignment-groups.md) | Assignment group management |
| [canvas assignments](canvas_assignments.md) | Assignment management |
| [canvas auth](canvas_auth.md) | Authentication commands |
| [canvas blueprint](canvas_blueprint.md) | Blueprint course management |
| [canvas cache](canvas_cache.md) | Cache management |
| [canvas calendar](canvas_calendar.md) | Calendar management |
| [canvas config](canvas_config.md) | Configuration management |
| [canvas content-migrations](canvas_content-migrations.md) | Content migration management |
| [canvas conversations](canvas_conversations.md) | Conversations (inbox) management |
| [canvas courses](canvas_courses.md) | Course management |
| [canvas discussions](canvas_discussions.md) | Discussion management |
| [canvas enrollments](canvas_enrollments.md) | Enrollment management |
| [canvas external-tools](canvas_external-tools.md) | External tools (LTI) management |
| [canvas files](canvas_files.md) | File management |
| [canvas grades](canvas_grades.md) | Gradebook management |
| [canvas groups](canvas_groups.md) | Group management |
| [canvas modules](canvas_modules.md) | Module management |
| [canvas outcomes](canvas_outcomes.md) | Learning outcomes management |
| [canvas overrides](canvas_overrides.md) | Assignment override management |
| [canvas pages](canvas_pages.md) | Page management |
| [canvas peer-reviews](canvas_peer-reviews.md) | Peer review management |
| [canvas planner](canvas_planner.md) | Planner management |
| [canvas quizzes](canvas_quizzes.md) | Quiz management |
| [canvas roles](canvas_roles.md) | Account role management |
| [canvas rubrics](canvas_rubrics.md) | Rubric management |
| [canvas sections](canvas_sections.md) | Course section management |
| [canvas sis-imports](canvas_sis-imports.md) | SIS import management |
| [canvas submissions](canvas_submissions.md) | Submission management |
| [canvas sync](canvas_sync.md) | Course synchronization |
| [canvas users](canvas_users.md) | User management |
| [canvas webhook](canvas_webhook.md) | Webhook server |

## Global Flags

All commands support the following global flags:

| Flag | Description |
|------|-------------|
| `--config` | Config file path (default: $HOME/.canvas-cli/config.yaml) |
| `--instance` | Canvas instance URL (overrides config) |
| `-o, --output` | Output format: table, json, yaml, csv |
| `-v, --verbose` | Enable verbose output |
| `--as-user` | Masquerade as another user (admin feature) |
| `--no-cache` | Disable caching of API responses |

## Usage Pattern

```bash
canvas <resource> <action> [flags]
```

For example:
```bash
canvas courses list                           # List all courses
canvas assignments get 123 --course-id 456    # Get assignment details
canvas submissions grade 789 --score 95       # Grade a submission
```
