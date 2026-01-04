---
description: Design comprehensive database schema for issue tracking system
---

Design a complete database schema for the issue tracking system with the following entities:

## Core Tables

### 1. Projects
- id (UUID, primary key)
- name (text, required)
- key (text, unique, e.g., "PROJ" for issue prefixes)
- description (text)
- created_at, updated_at (timestamps)

**Notes**: Keep simple for now; will be used for RBAC later

### 2. Statuses
- id (UUID, primary key)
- project_id (UUID, nullable - null means standard status)
- name (text, required, e.g., "To Do", "In Progress", "Done")
- category (text, e.g., "todo", "in_progress", "done" - for semantic grouping)
- position (integer, for ordering)
- is_standard (boolean)
- created_at (timestamp)

**Notes**: Projects define their own statuses. Support both standard statuses (project_id = null) and custom per-project statuses

### 3. Issue Types
- id (UUID, primary key)
- name (text, e.g., "Epic", "Story", "Task", "Bug")
- description (text)
- icon (text, emoji or icon identifier)
- hierarchy_level (integer, e.g., 0=Epic, 1=Story, 2=Task)
- allowed_parent_types (UUID[], array of type IDs that can be parents)
- is_standard (boolean)
- created_at (timestamp)

**Notes**: Defines type-based hierarchy rules (Epic → Story → Task)

### 4. Priorities
- id (UUID, primary key)
- name (text, e.g., "Critical", "High", "Medium", "Low")
- level (integer, for sorting)
- color (text, hex color code)
- icon (text)
- is_standard (boolean)
- created_at (timestamp)

### 5. Issues (Core Entity)
- id (UUID, primary key)
- project_id (UUID, foreign key → projects)
- issue_number (integer, auto-increment per project)
- issue_key (text, computed: PROJECT_KEY-NUMBER, e.g., "PROJ-123")
- type_id (UUID, foreign key → issue_types)
- status_id (UUID, foreign key → statuses)
- priority_id (UUID, foreign key → priorities)
- parent_id (UUID, nullable, self-referential foreign key)
- title (text, required)
- description (text, markdown)
- reporter_id (UUID, user reference - simple for now)
- assignee_id (UUID, nullable, user reference)
- created_at, updated_at (timestamps)
- resolved_at (timestamp, nullable)

**Constraints**: Validate parent_id type compatibility based on issue_types.allowed_parent_types

### 6. Comments
- id (UUID, primary key)
- issue_id (UUID, foreign key → issues)
- author_id (UUID, user reference)
- content (text, markdown)
- created_at, updated_at (timestamps)
- is_deleted (boolean, soft delete)

### 7. Work Log (Time Tracking)
- id (UUID, primary key)
- issue_id (UUID, foreign key → issues)
- user_id (UUID, user reference)
- time_spent_seconds (integer)
- description (text)
- started_at (timestamp)
- created_at (timestamp)

### 8. Activity Log (Audit Trail)
- id (UUID, primary key)
- issue_id (UUID, foreign key → issues)
- user_id (UUID, user reference)
- action (text, e.g., "created", "updated", "status_changed", "commented")
- field_name (text, nullable, e.g., "status", "assignee")
- old_value (text, nullable, JSON)
- new_value (text, nullable, JSON)
- created_at (timestamp)

**Notes**: Automatically track all changes to issues

### 9. Custom Field Definitions
- id (UUID, primary key)
- project_id (UUID, nullable - null means global)
- name (text, required)
- field_key (text, unique, e.g., "story_points")
- field_type (text, enum: "text_short", "text_long", "number", "date", "datetime", "select_single", "select_multi", "user_single", "user_multi", "issue_link")
- description (text)
- is_required (boolean)
- default_value (text, JSON)
- config (jsonb, for field-specific config like select options)
- created_at (timestamp)

### 10. Custom Field Values
- id (UUID, primary key)
- issue_id (UUID, foreign key → issues)
- field_definition_id (UUID, foreign key → custom_field_definitions)
- value_text (text, nullable)
- value_number (numeric, nullable)
- value_date (timestamp, nullable)
- value_json (jsonb, nullable - for arrays/objects)
- created_at, updated_at (timestamps)

**Notes**: Polymorphic storage - use appropriate column based on field_type

### 11. Boards
- id (UUID, primary key)
- project_id (UUID, foreign key → projects)
- name (text, required)
- description (text)
- is_default (boolean)
- created_at, updated_at (timestamps)

**Notes**: Boards are views within a project showing issues organized by status

### 12. Board Columns
- id (UUID, primary key)
- board_id (UUID, foreign key → boards)
- status_id (UUID, foreign key → statuses)
- position (integer, for ordering)
- wip_limit (integer, nullable)
- created_at (timestamp)

**Notes**: Maps statuses to board columns, defines status progression path

### 13. Filter Definitions
- id (UUID, primary key)
- project_id (UUID, nullable - null means cross-project filter)
- name (text, required)
- description (text)
- owner_id (UUID, user reference)
- is_shared (boolean)
- filter_expression (jsonb, structured query)
- created_at, updated_at (timestamps)

**Filter Expression Structure** (JSON):
```json
{
  "operator": "AND",
  "conditions": [
    {"field": "project_id", "operator": "equals", "value": "uuid"},
    {"field": "status.category", "operator": "in", "value": ["todo", "in_progress"]},
    {"field": "assignee_id", "operator": "equals", "value": "uuid"},
    {"field": "custom.story_points", "operator": "gte", "value": 5}
  ]
}
```

**Supported Operators**:
- Basic: equals, not_equals, is_null, is_not_null
- Text: contains, starts_with, ends_with (case-insensitive)
- Comparison: gt, gte, lt, lte, between
- Set: in, not_in

## Indexes
- issues(project_id, issue_number) - for issue key generation
- issues(project_id, status_id) - for board queries
- issues(parent_id) - for hierarchy traversal
- custom_field_values(issue_id, field_definition_id) - for custom field lookups
- activity_log(issue_id, created_at) - for timeline queries
- comments(issue_id, created_at) - for comment ordering

## Migration Files
Create PostgreSQL migration files (use golang-migrate or similar):
- Use UUIDs (uuid-ossp extension)
- Proper foreign key constraints with CASCADE/RESTRICT
- Indexes for performance
- Include both up and down migrations

Also provide Go structs with appropriate tags (`json`, `db`).
