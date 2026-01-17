package iterm2

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Writer implements the configfile.Writer interface for iTerm2 Dynamic Profiles.
// It generates JSON profile files from declarative YAML configuration.
type Writer struct{}

// NewWriter creates a new iTerm2 profile writer.
func NewWriter() *Writer {
	return &Writer{}
}

// Format returns the format identifier.
func (w *Writer) Format() string {
	return "iterm2"
}

// Write generates iTerm2 dynamic profile JSON from values.
func (w *Writer) Write(values map[string]any) ([]byte, error) {
	profile := make(map[string]any)

	// Process each section
	if p, ok := values["profile"].(map[string]any); ok {
		w.processProfile(profile, p)
	}
	if f, ok := values["font"].(map[string]any); ok {
		w.processFont(profile, f)
	}
	if t, ok := values["terminal"].(map[string]any); ok {
		w.processTerminal(profile, t)
	}
	if c, ok := values["cursor"].(map[string]any); ok {
		w.processCursor(profile, c)
	}
	if i, ok := values["input"].(map[string]any); ok {
		w.processInput(profile, i)
	}
	if m, ok := values["mouse"].(map[string]any); ok {
		w.processMouse(profile, m)
	}
	if b, ok := values["behavior"].(map[string]any); ok {
		w.processBehavior(profile, b)
	}
	if c, ok := values["colors"].(map[string]any); ok {
		w.processColors(profile, c)
	}
	if km, ok := values["keyboard_maps"].([]any); ok {
		w.processKeyboardMaps(profile, km)
	}
	if tags, ok := values["tags"].([]any); ok {
		w.processTags(profile, tags)
	}

	// Wrap in Profiles array as iTerm2 expects
	output := map[string]any{
		"Profiles": []any{profile},
	}

	return json.MarshalIndent(output, "", "  ")
}

// processProfile handles profile identity fields.
func (w *Writer) processProfile(profile, p map[string]any) {
	if name, ok := p["name"].(string); ok {
		profile["Name"] = name
	}
	if guid, ok := p["guid"].(string); ok {
		profile["Guid"] = guid
	}
	if parent, ok := p["parent"].(string); ok {
		profile["Dynamic Profile Parent Name"] = parent
	}
	if desc, ok := p["description"].(string); ok {
		profile["Description"] = desc
	}
}

// processFont handles font configuration.
func (w *Writer) processFont(profile, f map[string]any) {
	family := "Menlo"
	size := 14

	if fam, ok := f["family"].(string); ok {
		family = fam
	}
	if s, ok := f["size"].(int); ok {
		size = s
	}

	fontSpec := fmt.Sprintf("%s %d", family, size)
	profile["Normal Font"] = fontSpec
	profile["Non Ascii Font"] = fontSpec
	profile["Use Non-ASCII Font"] = false

	if aa, ok := f["anti_aliased"].(bool); ok {
		profile["ASCII Anti Aliased"] = aa
		profile["Non-ASCII Anti Aliased"] = aa
	} else {
		profile["ASCII Anti Aliased"] = true
		profile["Non-ASCII Anti Aliased"] = true
	}

	if h, ok := f["horizontal_spacing"]; ok {
		profile["Horizontal Spacing"] = toFloat(h)
	} else {
		profile["Horizontal Spacing"] = 1
	}
	if v, ok := f["vertical_spacing"]; ok {
		profile["Vertical Spacing"] = toFloat(v)
	} else {
		profile["Vertical Spacing"] = 1.1
	}

	if ub, ok := f["use_bold"].(bool); ok {
		profile["Use Bold Font"] = ub
	} else {
		profile["Use Bold Font"] = true
	}
	if ubb, ok := f["use_bright_bold"].(bool); ok {
		profile["Use Bright Bold"] = ubb
	} else {
		profile["Use Bright Bold"] = true
	}
	if ui, ok := f["use_italic"].(bool); ok {
		profile["Use Italic Font"] = ui
	} else {
		profile["Use Italic Font"] = true
	}
}

// processTerminal handles terminal settings.
func (w *Writer) processTerminal(profile, t map[string]any) {
	if tt, ok := t["type"].(string); ok {
		profile["Terminal Type"] = tt
	} else {
		profile["Terminal Type"] = "xterm-256color"
	}

	profile["Character Encoding"] = 4 // UTF-8

	if sl, ok := t["scrollback_lines"].(int); ok {
		profile["Scrollback Lines"] = sl
	} else {
		profile["Scrollback Lines"] = 50000
	}

	if us, ok := t["unlimited_scrollback"].(bool); ok {
		profile["Unlimited Scrollback"] = us
	} else {
		profile["Unlimited Scrollback"] = false
	}

	profile["Scrollback With Status Bar"] = true
}

// processCursor handles cursor settings.
func (w *Writer) processCursor(profile, c map[string]any) {
	if ct, ok := c["type"].(string); ok {
		profile["Cursor Type"] = w.cursorTypeToInt(ct)
	} else {
		profile["Cursor Type"] = 0 // block
	}

	if b, ok := c["blinking"].(bool); ok {
		profile["Blinking Cursor"] = b
	} else {
		profile["Blinking Cursor"] = false
	}

	if boost, ok := c["boost"].(int); ok {
		profile["Cursor Boost"] = boost
	} else {
		profile["Cursor Boost"] = 0
	}

	profile["Minimum Contrast"] = 0
}

// cursorTypeToInt converts cursor type name to iTerm2 int.
func (w *Writer) cursorTypeToInt(ct string) int {
	switch strings.ToLower(ct) {
	case "underline":
		return 1
	case "bar", "ibeam", "vertical":
		return 2
	default: // block
		return 0
	}
}

// processInput handles keyboard/input settings.
func (w *Writer) processInput(profile, i map[string]any) {
	if opt, ok := i["option_key_sends"].(string); ok {
		profile["Option Key Sends"] = w.optionKeyToInt(opt)
	} else {
		profile["Option Key Sends"] = 2 // meta
	}

	if ropt, ok := i["right_option_key_sends"].(string); ok {
		profile["Right Option Key Sends"] = w.optionKeyToInt(ropt)
	} else {
		profile["Right Option Key Sends"] = 0 // normal
	}
}

// optionKeyToInt converts option key mode to iTerm2 int.
func (w *Writer) optionKeyToInt(mode string) int {
	switch strings.ToLower(mode) {
	case "meta":
		return 2
	case "esc+", "esc":
		return 1
	default: // normal
		return 0
	}
}

// processMouse handles mouse settings.
func (w *Writer) processMouse(profile, m map[string]any) {
	if mr, ok := m["reporting"].(bool); ok {
		profile["Mouse Reporting"] = mr
	} else {
		profile["Mouse Reporting"] = true
	}

	if aw, ok := m["allow_wheel"].(bool); ok {
		profile["Mouse Reporting Allow Mouse Wheel"] = aw
	} else {
		profile["Mouse Reporting Allow Mouse Wheel"] = true
	}
}

// processBehavior handles window/session behavior.
func (w *Writer) processBehavior(profile, b map[string]any) {
	if wd, ok := b["working_directory"].(string); ok {
		profile["Custom Directory"] = w.workingDirToString(wd)
	} else {
		profile["Custom Directory"] = "Recycle"
	}
	profile["Working Directory"] = ""

	if coe, ok := b["close_on_end"].(bool); ok {
		profile["Close Sessions On End"] = coe
	} else {
		profile["Close Sessions On End"] = true
	}

	profile["Prompt Before Closing 2"] = 0
	profile["Send Code When Idle"] = false
	profile["Idle Code"] = 0

	if sb, ok := b["silence_bell"].(bool); ok {
		profile["Silence Bell"] = sb
	} else {
		profile["Silence Bell"] = false
	}

	if ft, ok := b["flash_tab"].(bool); ok {
		profile["Flash Tab"] = ft
	} else {
		profile["Flash Tab"] = true
	}

	profile["BM Growl"] = true
}

// workingDirToString converts working directory mode to iTerm2 string.
func (w *Writer) workingDirToString(mode string) string {
	switch strings.ToLower(mode) {
	case "home":
		return "Home"
	case "custom":
		return "Yes"
	default: // recycle
		return "Recycle"
	}
}

// processColors handles color scheme or inline colors.
func (w *Writer) processColors(profile, c map[string]any) {
	// Check for scheme reference first
	if scheme, ok := c["scheme"].(string); ok {
		if cs, exists := ColorSchemes[scheme]; exists {
			w.applyColorScheme(profile, cs)
		}
	}

	// Inline colors override scheme
	if bg, ok := c["background"].(string); ok {
		profile["Background Color"] = w.hexToITermColor(bg)
	}
	if fg, ok := c["foreground"].(string); ok {
		profile["Foreground Color"] = w.hexToITermColor(fg)
	}
	if bold, ok := c["bold"].(string); ok {
		profile["Bold Color"] = w.hexToITermColor(bold)
	}
	if cursor, ok := c["cursor"].(string); ok {
		profile["Cursor Color"] = w.hexToITermColor(cursor)
	}
	if cursorText, ok := c["cursor_text"].(string); ok {
		profile["Cursor Text Color"] = w.hexToITermColor(cursorText)
	}
	if selection, ok := c["selection"].(string); ok {
		profile["Selection Color"] = w.hexToITermColor(selection)
	}
	if selectedText, ok := c["selected_text"].(string); ok {
		profile["Selected Text Color"] = w.hexToITermColor(selectedText)
	}

	// Handle ANSI colors array
	if ansi, ok := c["ansi"].([]any); ok && len(ansi) == 16 {
		for i, color := range ansi {
			if hex, ok := color.(string); ok {
				profile[fmt.Sprintf("Ansi %d Color", i)] = w.hexToITermColor(hex)
			}
		}
	}
}

// applyColorScheme applies a predefined color scheme to the profile.
func (w *Writer) applyColorScheme(profile map[string]any, cs ColorScheme) {
	profile["Background Color"] = w.hexToITermColor(cs.Background)
	profile["Foreground Color"] = w.hexToITermColor(cs.Foreground)
	profile["Bold Color"] = w.hexToITermColor(cs.Bold)
	profile["Cursor Color"] = w.hexToITermColor(cs.Cursor)
	profile["Cursor Text Color"] = w.hexToITermColor(cs.CursorText)
	profile["Selection Color"] = w.hexToITermColor(cs.Selection)
	profile["Selected Text Color"] = w.hexToITermColor(cs.SelectedText)

	for i, hex := range cs.Ansi {
		profile[fmt.Sprintf("Ansi %d Color", i)] = w.hexToITermColor(hex)
	}
}

// hexToITermColor converts #RRGGBB to iTerm2 color dict with normalized RGB.
func (w *Writer) hexToITermColor(hex string) map[string]any {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		// Return black for invalid colors
		return map[string]any{
			"Red Component":   0.0,
			"Green Component": 0.0,
			"Blue Component":  0.0,
			"Color Space":     "sRGB",
		}
	}

	r, _ := strconv.ParseInt(hex[0:2], 16, 64)
	g, _ := strconv.ParseInt(hex[2:4], 16, 64)
	b, _ := strconv.ParseInt(hex[4:6], 16, 64)

	return map[string]any{
		"Red Component":   float64(r) / 255.0,
		"Green Component": float64(g) / 255.0,
		"Blue Component":  float64(b) / 255.0,
		"Color Space":     "sRGB",
	}
}

// processKeyboardMaps handles keyboard mapping configuration.
func (w *Writer) processKeyboardMaps(profile map[string]any, kms []any) {
	keymap := make(map[string]any)

	for _, km := range kms {
		if m, ok := km.(map[string]any); ok {
			if key, ok := m["key"].(string); ok {
				entry := make(map[string]any)
				if action, ok := m["action"].(int); ok {
					entry["Action"] = action
				}
				if text, ok := m["text"].(string); ok {
					entry["Text"] = text
				}
				keymap[key] = entry
			}
		}
	}

	if len(keymap) > 0 {
		profile["Keyboard Map"] = keymap
	}
}

// processTags handles profile tags.
func (w *Writer) processTags(profile map[string]any, tags []any) {
	tagStrings := make([]string, 0, len(tags))
	for _, t := range tags {
		if s, ok := t.(string); ok {
			tagStrings = append(tagStrings, s)
		}
	}
	if len(tagStrings) > 0 {
		profile["Tags"] = tagStrings
	}
}

// toFloat converts various numeric types to float64.
func toFloat(v any) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	default:
		return 1.0
	}
}
