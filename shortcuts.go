package valve_vdf_binary

import (
	"errors"
	"io"
	"strconv"
)

type Shortcut struct {
	AppName string
	Exe     string
}

func ParseShortcuts(buf io.Reader) ([]Shortcut, error) {
	vdf, err := Parse(buf)
	if err != nil {
		return []Shortcut{}, err
	}

	shortcutsMap, ok := vdf.GetMap("shortcuts")
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

		appName, ok := s.GetString("AppName")
		if !ok {
			return []Shortcut{}, errors.New("Could not get key 'AppName' for one of the shortcuts")
		}

		exe, ok := s.GetString("Exe")
		if !ok {
			return []Shortcut{}, errors.New("Could not get key 'Exe' for one of the shortcuts")
		}

		shortcuts[i] = Shortcut{
			AppName: appName,
			Exe:     exe,
		}
	}

	return shortcuts, nil
}
