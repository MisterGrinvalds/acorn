package configfile

// ITerm2ColorScheme represents an iTerm2 color palette.
// Colors are stored as hex strings (#RRGGBB) and converted to
// iTerm2's normalized RGB format (0.0-1.0) during generation.
type ITerm2ColorScheme struct {
	Background   string
	Foreground   string
	Bold         string
	Cursor       string
	CursorText   string
	Selection    string
	SelectedText string
	// Ansi colors 0-15 (standard terminal palette)
	Ansi [16]string
}

// ITerm2ColorSchemes contains predefined color schemes.
var ITerm2ColorSchemes = map[string]ITerm2ColorScheme{
	"catppuccin-mocha": {
		Background:   "#1e1e2e",
		Foreground:   "#cdd6f4",
		Bold:         "#cdd6f4",
		Cursor:       "#f38ba8",
		CursorText:   "#1e1e2e",
		Selection:    "#3a3c53",
		SelectedText: "#cdd6f4",
		Ansi: [16]string{
			"#45475a", // 0: Black
			"#f38ba8", // 1: Red
			"#a6e3a1", // 2: Green
			"#f9e2af", // 3: Yellow
			"#89b4fa", // 4: Blue
			"#f5c2e7", // 5: Magenta
			"#94e2d5", // 6: Cyan
			"#bac2de", // 7: White
			"#585b70", // 8: Bright Black
			"#f38ba8", // 9: Bright Red
			"#a6e3a1", // 10: Bright Green
			"#f9e2af", // 11: Bright Yellow
			"#89b4fa", // 12: Bright Blue
			"#f5c2e7", // 13: Bright Magenta
			"#94e2d5", // 14: Bright Cyan
			"#cdd6f4", // 15: Bright White
		},
	},
	"catppuccin-latte": {
		Background:   "#eff1f5",
		Foreground:   "#4c4f69",
		Bold:         "#4c4f69",
		Cursor:       "#d20f39",
		CursorText:   "#eff1f5",
		Selection:    "#acb0be",
		SelectedText: "#4c4f69",
		Ansi: [16]string{
			"#5c5f77", // 0: Black (Subtext 1)
			"#d20f39", // 1: Red
			"#40a02b", // 2: Green
			"#df8e1d", // 3: Yellow
			"#1e66f5", // 4: Blue
			"#ea76cb", // 5: Magenta (Pink)
			"#179299", // 6: Cyan (Teal)
			"#acb0be", // 7: White (Surface 2)
			"#6c6f85", // 8: Bright Black (Subtext 0)
			"#d20f39", // 9: Bright Red
			"#40a02b", // 10: Bright Green
			"#df8e1d", // 11: Bright Yellow
			"#1e66f5", // 12: Bright Blue
			"#ea76cb", // 13: Bright Magenta
			"#179299", // 14: Bright Cyan
			"#4c4f69", // 15: Bright White (Text)
		},
	},
	"catppuccin-frappe": {
		Background:   "#303446",
		Foreground:   "#c6d0f5",
		Bold:         "#c6d0f5",
		Cursor:       "#ea999c",
		CursorText:   "#303446",
		Selection:    "#51576d",
		SelectedText: "#c6d0f5",
		Ansi: [16]string{
			"#51576d", // 0: Black (Surface 1)
			"#e78284", // 1: Red
			"#a6d189", // 2: Green
			"#e5c890", // 3: Yellow
			"#8caaee", // 4: Blue
			"#f4b8e4", // 5: Magenta (Pink)
			"#81c8be", // 6: Cyan (Teal)
			"#b5bfe2", // 7: White (Subtext 1)
			"#626880", // 8: Bright Black (Surface 2)
			"#e78284", // 9: Bright Red
			"#a6d189", // 10: Bright Green
			"#e5c890", // 11: Bright Yellow
			"#8caaee", // 12: Bright Blue
			"#f4b8e4", // 13: Bright Magenta
			"#81c8be", // 14: Bright Cyan
			"#c6d0f5", // 15: Bright White (Text)
		},
	},
	"catppuccin-macchiato": {
		Background:   "#24273a",
		Foreground:   "#cad3f5",
		Bold:         "#cad3f5",
		Cursor:       "#ee99a0",
		CursorText:   "#24273a",
		Selection:    "#494d64",
		SelectedText: "#cad3f5",
		Ansi: [16]string{
			"#494d64", // 0: Black (Surface 1)
			"#ed8796", // 1: Red
			"#a6da95", // 2: Green
			"#eed49f", // 3: Yellow
			"#8aadf4", // 4: Blue
			"#f5bde6", // 5: Magenta (Pink)
			"#8bd5ca", // 6: Cyan (Teal)
			"#b8c0e0", // 7: White (Subtext 1)
			"#5b6078", // 8: Bright Black (Surface 2)
			"#ed8796", // 9: Bright Red
			"#a6da95", // 10: Bright Green
			"#eed49f", // 11: Bright Yellow
			"#8aadf4", // 12: Bright Blue
			"#f5bde6", // 13: Bright Magenta
			"#8bd5ca", // 14: Bright Cyan
			"#cad3f5", // 15: Bright White (Text)
		},
	},
}
