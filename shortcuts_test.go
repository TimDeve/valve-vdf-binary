package valve_vdf_binary

import (
	"bufio"
	"bytes"
	_ "embed"
	"github.com/stretchr/testify/require"
	"testing"
)

//go:embed test-files/shortcuts.vdf
var shortcutVdf []byte

func TestParseShortcut(t *testing.T) {
	parsedShortcuts, err := parseShortcutsFromBytes(shortcutVdf)
	if err != nil {
		t.Error(err)
		return
	}

	shortcuts := []Shortcut{{"Control"}, {"Cyberpunk 2077"}, {"Skate 3"}}

	require.Equal(t, parsedShortcuts, shortcuts)
}

func parseShortcutsFromBytes(buf []byte) ([]Shortcut, error) {
	return parseShortcuts(bufio.NewReader(bytes.NewReader(buf)))
}
