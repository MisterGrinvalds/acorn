# Update Forked Plugins

Sync all forked repositories with their upstream sources.

## Workflow

### Step 1: Identify Forked Plugins

Scan configuration for plugins using the `fork()` wrapper:

Search for patterns:
- `require('lib.forks').fork`
- `fork('owner/repo')`

Parse out original `owner/repo` names from each `fork('...')` call.

### Step 2: List All Forks

Create a list of forks to sync:

| Fork Repo | Upstream |
|-----------|----------|
| MisterGrinvalds/ThePrimeagen.harpoon | ThePrimeagen/harpoon |
| MisterGrinvalds/nvim-telescope.telescope.nvim | nvim-telescope/telescope.nvim |

### Step 3: Sync Each Fork

For each fork, sync with upstream using GitHub CLI:

```bash
gh repo sync MisterGrinvalds/{owner}.{repo} --source {owner}/{repo}
```

This fetches from upstream and merges into the fork's default branch.

### Step 4: Handle Conflicts

If sync fails due to conflicts:

1. Report which fork has conflicts
2. Suggest manual resolution:
   ```bash
   # Clone the fork locally
   gh repo clone MisterGrinvalds/{owner}.{repo}
   cd {owner}.{repo}

   # Add upstream remote
   git remote add upstream https://github.com/{owner}/{repo}.git

   # Fetch and merge manually
   git fetch upstream
   git merge upstream/main  # or upstream/master

   # Resolve conflicts, then push
   git push origin main
   ```

### Step 5: Update Lazy.nvim

After syncing forks on GitHub, update local plugins:

Suggest running in Neovim:
```
:Lazy update
```

### Step 6: Report Results

| Fork | Upstream | Status |
|------|----------|--------|
| MisterGrinvalds/ThePrimeagen.harpoon | ThePrimeagen/harpoon | Synced |
| MisterGrinvalds/nvim-telescope.telescope.nvim | nvim-telescope/telescope.nvim | Synced |

Summary:
- Total forks: X
- Successfully synced: Y
- Conflicts: Z (list if any)

## Arguments (Optional)

$ARGUMENTS - Specific fork to update (e.g., `ThePrimeagen/harpoon`)

If provided, only sync that specific fork instead of all forks.

## Notes

- `gh repo sync` requires the fork to have upstream configured (done automatically by `gh repo fork`)
- If upstream branch name differs (main vs master), sync may need adjustment
- Forks with custom commits will have those commits preserved (merged with upstream)
- Run this periodically to stay up-to-date with plugin improvements and bug fixes
