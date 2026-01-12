---
description: Design RESTful API endpoints for issue tracking system
---

Design comprehensive REST API endpoints for the issue tracking system:

## Projects

- `GET /api/v1/projects` - List all projects (paginated)
- `POST /api/v1/projects` - Create new project
- `GET /api/v1/projects/:id` - Get project details
- `PUT /api/v1/projects/:id` - Update project
- `DELETE /api/v1/projects/:id` - Delete project
- `GET /api/v1/projects/:id/statuses` - Get project statuses (standard + custom)
- `POST /api/v1/projects/:id/statuses` - Create custom status for project
- `PUT /api/v1/projects/:id/statuses/:status_id` - Update status
- `DELETE /api/v1/projects/:id/statuses/:status_id` - Delete custom status

## Issues

- `GET /api/v1/projects/:project_id/issues` - List issues (with filters, pagination)
  - Query params: `type`, `status`, `assignee`, `priority`, `parent`, `custom_field_*`
- `POST /api/v1/projects/:project_id/issues` - Create new issue
- `GET /api/v1/issues/:id` - Get issue details (includes comments, work log, activity)
- `PUT /api/v1/issues/:id` - Update issue
- `DELETE /api/v1/issues/:id` - Delete issue
- `PATCH /api/v1/issues/:id/status` - Update issue status
- `PATCH /api/v1/issues/:id/assign` - Assign issue to user
- `GET /api/v1/issues/:id/children` - Get child issues (hierarchy)
- `GET /api/v1/issues/:id/ancestors` - Get parent chain (breadcrumbs)

## Comments

- `GET /api/v1/issues/:issue_id/comments` - List comments
- `POST /api/v1/issues/:issue_id/comments` - Add comment
- `PUT /api/v1/comments/:id` - Update comment
- `DELETE /api/v1/comments/:id` - Delete comment (soft delete)

## Work Log

- `GET /api/v1/issues/:issue_id/worklog` - List work log entries
- `POST /api/v1/issues/:issue_id/worklog` - Log time
- `PUT /api/v1/worklog/:id` - Update work log entry
- `DELETE /api/v1/worklog/:id` - Delete work log entry

## Activity Log

- `GET /api/v1/issues/:issue_id/activity` - Get issue activity timeline
- Note: Activity log is automatically created on issue changes (no POST/PUT/DELETE)

## Custom Fields

- `GET /api/v1/projects/:project_id/custom-fields` - List custom field definitions
- `POST /api/v1/projects/:project_id/custom-fields` - Create custom field
- `PUT /api/v1/custom-fields/:id` - Update custom field definition
- `DELETE /api/v1/custom-fields/:id` - Delete custom field definition
- Note: Custom field values are set/retrieved through issue endpoints

## Issue Types

- `GET /api/v1/issue-types` - List all issue types
- `POST /api/v1/issue-types` - Create custom issue type
- `PUT /api/v1/issue-types/:id` - Update issue type
- `DELETE /api/v1/issue-types/:id` - Delete custom issue type
- `GET /api/v1/issue-types/:id/hierarchy` - Get allowed parent/child types

## Priorities

- `GET /api/v1/priorities` - List all priorities
- `POST /api/v1/priorities` - Create custom priority
- `PUT /api/v1/priorities/:id` - Update priority
- `DELETE /api/v1/priorities/:id` - Delete custom priority

## Boards

- `GET /api/v1/projects/:project_id/boards` - List boards in project
- `POST /api/v1/projects/:project_id/boards` - Create board
- `GET /api/v1/boards/:id` - Get board details with columns and issues
- `PUT /api/v1/boards/:id` - Update board
- `DELETE /api/v1/boards/:id` - Delete board
- `POST /api/v1/boards/:id/columns` - Add column to board (map status)
- `PUT /api/v1/boards/:id/columns/:column_id` - Update column (position, WIP limit)
- `DELETE /api/v1/boards/:id/columns/:column_id` - Remove column from board

## Filters

- `GET /api/v1/filters` - List user's filters
- `GET /api/v1/projects/:project_id/filters` - List project filters
- `POST /api/v1/filters` - Create filter
- `GET /api/v1/filters/:id` - Get filter details
- `PUT /api/v1/filters/:id` - Update filter
- `DELETE /api/v1/filters/:id` - Delete filter
- `POST /api/v1/filters/:id/execute` - Execute filter and return matching issues

## Request/Response Schemas

For each endpoint, define:
- **Request body schemas** (JSON)
- **Response schemas** (JSON)
- **HTTP status codes**:
  - 200 OK - Successful GET/PUT/PATCH
  - 201 Created - Successful POST
  - 204 No Content - Successful DELETE
  - 400 Bad Request - Validation error
  - 404 Not Found - Resource not found
  - 409 Conflict - Constraint violation (e.g., hierarchy rules)
  - 422 Unprocessable Entity - Business logic error (e.g., WIP limit)
  - 500 Internal Server Error

- **Query parameters** for list endpoints:
  - `page`, `page_size` - Pagination
  - `sort` - Sort field and direction (e.g., `created_at:desc`)
  - Filter-specific params

- **Error response format**:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid issue type for parent relationship",
    "details": {
      "field": "parent_id",
      "constraint": "Epic cannot be child of Story"
    }
  }
}
```

## Examples

Provide example curl commands or HTTP requests for key endpoints:
- Creating a project
- Creating an issue with custom fields
- Moving an issue between statuses
- Executing a filter
- Viewing a board
