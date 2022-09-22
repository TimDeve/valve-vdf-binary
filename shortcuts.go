package valve_vdf_binary

import (
	"bufio"
	"errors"
	"strconv"
)

type Shortcut struct {
	AppName string
}

func parseShortcuts(buf *bufio.Reader) ([]Shortcut, error) {
	vdf, err := parse(buf)
	if err != nil {
		return []Shortcut{}, err
	}

	shortcutsMap, ok := vdf.getMap("shortcuts")
	if !ok {
		return []Shortcut{}, errors.New("Could not find 'shortcuts' in parsed vdf")
	}

	shortcuts := make([]Shortcut, len(shortcutsMap))

	for i := range shortcuts {
		key := strconv.Itoa(i)
		s, ok := shortcutsMap[key]
		if !ok {
			return []Shortcut{}, errors.New("vdf that should be an array does not have the corresponding index")
		}
		appName, ok := s.getString("AppName")
		shortcuts[i] = Shortcut{
			AppName: appName,
		}
	}

	return shortcuts, nil
}
