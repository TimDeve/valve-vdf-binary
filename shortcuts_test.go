package valve_vdf_binary_test

import (
	"bytes"
	_ "embed"
	"testing"

	vdf "github.com/TimDeve/valve-vdf-binary"
	test "github.com/stretchr/testify/require"
)

//go:embed test-files/shortcuts.vdf
var shortcutVdf []byte

func TestParseShortcut(t *testing.T) {
	shortcuts, err := parseShortcutsFromBytes(shortcutVdf)

	test.Nil(t, err)

	expected := []vdf.Shortcut{
		{
			AppId:    3414143657,
			AppName:  "Control",
			Exe:      "\"C:\\Program Files\\Epic Games\\Control\\Control_DX12.exe\"",
			Icon:     "",
			IsHidden: true,
			StartDir: "\"C:\\Program Files\\Epic Games\\Control\\\"",
			Tags:     []string{},
		},
		{
			AppId:    3022575626,
			AppName:  "Cyberpunk 2077",
			Exe:      "\"C:\\Program Files (x86)\\GOG Galaxy\\GalaxyClient.exe\"",
			Icon:     "C:\\Users\\user\\Downloads\\cyberpunk.ico",
			IsHidden: false,
			StartDir: "\"C:\\Program Files (x86)\\GOG Galaxy\\\"",
			Tags:     []string{"favorite"},
		},
		{
			AppId:    3043193801,
			AppName:  "Skate 3",
			Exe:      "\"C:\\Users\\user\\scoop\\apps\\RPCS3\\current\\rpcs3.exe\"",
			Icon:     "",
			IsHidden: false,
			StartDir: "\"C:\\Users\\user\\scoop\\apps\\RPCS3\\current\\\"",
			Tags:     []string{"Sport", "Action", "Skate"},
		},
	}

	test.Equal(t, expected, shortcuts)
}

func parseShortcutsFromBytes(buf []byte) ([]vdf.Shortcut, error) {
	return vdf.ParseShortcuts(bytes.NewReader(buf))
}
