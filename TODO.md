# TODO - Acorn Development

**Last Updated:** 2026-01-17
**Current Branch:** feat/claude-tool-agents-commands

## 🔥 High Priority

### 1. Fix Pre-existing Test Failures
- [x] Fix `internal/ide/neovim/neovim.go:206` - non-constant format string in fmt.Errorf
- [x] Fix `internal/terminal/shell/shell_test.go:330` - test expects "shell" but gets "entrypoint"
- [x] Move config writers back to component packages to maintain proper organization

### 2. Document New Structure
- [ ] Add README.md to explain internal/ reorganization
- [ ] Add README.md in key categories (ai/, utils/, etc.) explaining their purpose
- [ ] Update CONTRIBUTING.md with new package organization guidelines

### 3. Commit Reorganization Work
- [x] Review all changes in git status
- [x] Create commit with clear message about reorganization
- [ ] Update .gitignore if needed for new structure

## 📋 Medium Priority

### 4. Code Quality
- [ ] Run `go vet ./...` and fix any issues
- [ ] Run `golangci-lint` if available
- [ ] Check for unused imports

### 5. Testing
- [ ] Ensure all existing tests pass
- [ ] Add tests for any untested packages
- [ ] Verify CI/CD pipeline works with new structure

### 6. Documentation Updates
- [ ] Update any docs referencing old internal/ structure
- [ ] Check if any examples need updating
- [ ] Update architecture diagrams if they exist

## 🎯 Future Enhancements

### 7. Consider Additional Categories
Based on the layout file, consider adding:
- [ ] `internal/cloud/aws/` (if AWS integration planned)
- [ ] `internal/cloud/azure/` (if Azure integration planned)
- [ ] `internal/devops/docker/` (if Docker tooling planned)
- [ ] `internal/devops/podman/` (if Podman support planned)
- [ ] `internal/database/postgres/` (move database.go into postgres/)

### 8. Code Organization
- [ ] Review if `database/` should have subdirectories
- [ ] Consider splitting large packages if needed
- [ ] Look for code that could be shared across categories

### 9. Build Improvements
- [ ] Optimize build times
- [ ] Add build tags if needed
- [ ] Review dependencies

## 📝 Notes

### Recent Changes
- ✅ **2026-01-17:** Completed internal/ folder reorganization
  - Moved 29 packages into 11 categorized directories
  - Updated all import paths (79 replacements)
  - Renamed cmd/gitcmd.go → cmd/git.go
  - Build passes, CLI works

### Architecture Decisions
- **Flat categories:** Chose single-level categorization over deep nesting
- **Utils over Core:** Named shared utilities "utils" for clarity
- **Database placement:** Kept at top level for now, can expand to database/postgres later

### Blockers
- None currently

### Questions for User
- Should old SESSION-*.md files be kept in root or moved to archive?
- Do we want to add more cloud providers (AWS, Azure)?
- Should database.go move into database/postgres/ subdirectory?

---

## 🗂️ Archive

### Completed (2026-01-17)
- ✅ Internal folder reorganization from 29 scattered packages to 11 categories
- ✅ Import path updates across entire codebase
- ✅ Session documentation and archiving
