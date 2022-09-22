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
		{AppName: "Control", Exe: "\"C:\\Program Files\\Epic Games\\Control\\Control_DX12.exe\""},
		{AppName: "Cyberpunk 2077", Exe: "\"C:\\Program Files (x86)\\GOG Galaxy\\GalaxyClient.exe\""},
		{AppName: "Skate 3", Exe: "\"C:\\Users\\user\\scoop\\apps\\RPCS3\\current\\rpcs3.exe\""},
	}

	test.Equal(t, expected, shortcuts)
}

func parseShortcutsFromBytes(buf []byte) ([]vdf.Shortcut, error) {
	return vdf.ParseShortcuts(bytes.NewReader(buf))
}
