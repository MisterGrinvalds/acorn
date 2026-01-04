---
description: Create a new Cobra CLI router with proper structure
argument-hint: <router-name> [description]
model: claude-sonnet-4-5-20250929
allowed-tools: Write, Edit, Bash(go fmt:*), Bash(go mod tidy:*)
---

# Create Cobra Router

Create a new Cobra CLI router at `cmd/$1/$1.go` with the description: $2

## Steps

1. **Extract router information**
   - Router name: `$1` (will be used for directory and file names)
   - Description: `$2` (or use a default description if not provided)
   - Ensure router name is in lowercase

2. **Create directory structure**
   - Create `cmd/$1/` directory if it doesn't exist
   - Verify the path is created successfully

3. **Generate the router file**
   - Create `cmd/$1/$1.go` with proper Cobra command structure
   - Include:
     - Package declaration
     - Necessary imports (github.com/spf13/cobra)
     - Command variable with Use, Short, Long, and Run fields
     - init() function to add flags if needed
     - Proper error handling
     - Example usage comment

4. **Code structure template**
   ```go
   package <router-name>

   import (
       "fmt"
       
       "github.com/spf13/cobra"
   )

   var Cmd = &cobra.Command{
       Use:   "<router-name>",
       Short: "<short description>",
       Long:  `<longer description>`,
       Run: func(cmd *cobra.Command, args []string) {
           fmt.Println("<router-name> called")
           // TODO: Implement router logic
       },
   }

   func init() {
       // Add flags here
       // Cmd.Flags().StringP("flag", "f", "", "flag description")
   }
   ```

5. **Format and verify**
   - Run `go fmt` on the created file
   - Confirm the file was created successfully
   - Show the file path and next steps

6. **Provide integration instructions**
   - Remind to import the new router in main.go or root command
   - Show example: `rootCmd.AddCommand(<router-name>.Cmd)`
   - Suggest testing with: `go run main.go <router-name> --help`

## Notes
- Use kebab-case for multi-word router names
- Follow Go naming conventions
- Include helpful comments for future modifications