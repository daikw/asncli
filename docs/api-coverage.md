# Asana API Coverage

Feature matrix comparing Asana REST API capabilities with asncli support.

Legend: Supported / Not supported

## Core Resources

| Resource | Operation | API | CLI | CLI Command |
|----------|-----------|-----|-----|-------------|
| **Tasks** | List | Yes | Yes | `tasks list` |
| | Search | Yes | Yes | `tasks search` |
| | Get | Yes | Yes | `tasks get` |
| | Create | Yes | Yes | `tasks create` |
| | Update | Yes | Yes | `tasks update` |
| | Delete | Yes | No | |
| | Duplicate | Yes | No | |
| | Set parent | Yes | No | |
| | Add/Remove project | Yes | No | |
| | Add/Remove tag | Yes | No | |
| | Add/Remove followers | Yes | No | |
| | Add/Remove dependencies | Yes | No | |
| **Subtasks** | List | Yes | Yes | `tasks subtasks` |
| | Create | Yes | No | |
| **Stories (Comments)** | List | Yes | Yes | `tasks comments` |
| | Get | Yes | No | |
| | Create | Yes | Yes | `tasks comment-add` |
| | Update | Yes | Yes | `tasks comment-update` |
| | Delete | Yes | Yes | `tasks comment-delete` |
| **Attachments** | List | Yes | Yes | `tasks attachments` |
| | Get | Yes | Yes | `tasks attachment-get` |
| | Upload | Yes | No | |
| | Delete | Yes | No | |
| **Projects** | List | Yes | Yes | `projects list` |
| | Get | Yes | Yes | `projects get` |
| | Create | Yes | No | |
| | Update | Yes | No | |
| | Delete | Yes | No | |
| | Duplicate | Yes | No | |
| | Search | Yes | No | |
| | Task count | Yes | No | |
| **Custom Fields** | List (workspace) | Yes | Yes | `custom-fields list` |
| | Get | Yes | Yes | `custom-fields get` |
| | Create | Yes | No | |
| | Update | Yes | Yes | `custom-fields update` |
| | Delete | Yes | No | |
| **Users** | Get me | Yes | Yes | `auth status` |
| | List | Yes | No | |
| | Get | Yes | No | |
| | Update | Yes | No | |

## Organization & Team Resources

| Resource | Operation | API | CLI | Notes |
|----------|-----------|-----|-----|-------|
| **Workspaces** | List / Get / Update | Yes | No | |
| **Workspace Memberships** | List (user) | Yes | Yes | Used internally by `auth status`, `config set-workspace` |
| **Teams** | CRUD / Members | Yes | No | |
| **Team Memberships** | List | Yes | No | |
| **Sections** | CRUD / Add task | Yes | No | |
| **Tags** | CRUD | Yes | No | |
| **Memberships** | CRUD | Yes | No | |

## Advanced Features

| Resource | Operation | API | CLI |
|----------|-----------|-----|-----|
| **Goals** | CRUD / Metrics / Relationships | Yes | No |
| **Portfolios** | CRUD / Items / Members | Yes | No |
| **Status Updates** | CRUD | Yes | No |
| **Project Statuses** | CRUD | Yes | No |
| **Project Briefs** | CRUD | Yes | No |
| **Rules** | Trigger | Yes | No |
| **Webhooks** | CRUD | Yes | No |
| **Events** | Get | Yes | No |
| **Time Tracking** | CRUD | Yes | No |
| **Typeahead** | Search | Yes | No |
| **Batch API** | Submit | Yes | No |
| **Audit Log** | Get | Yes | No |
| **Allocations** | CRUD | Yes | No |
| **Budgets** | CRUD | Yes | No |
| **Rates** | CRUD | Yes | No |
| **Roles** | CRUD | Yes | No |

## Summary

| Category | Supported | Not supported |
|----------|-----------|---------------|
| Core resource operations | 20 | 22 |
| Organization & team | 1 (internal) | 6 resources |
| Advanced features | 0 | 15 resources |

Source: [Asana REST API Reference](https://developers.asana.com/reference/rest-api-reference)
