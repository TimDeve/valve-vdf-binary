package valve_vdf_binary

import (
	"errors"
	"io"
	"strconv"
)

type Shortcut struct {
	AppId    uint32
	AppName  string
	Exe      string
	Icon     string
	IsHidden bool // Doesn't seem to be used by steam anymore
	StartDir string
	Tags     []string
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

		appId, ok := s.GetUint("appid")
		if !ok {
			return []Shortcut{}, errors.New("Could not get key 'appid' for one of the shortcuts")
		}

		appName, ok := s.GetString("AppName")
		if !ok {
			return []Shortcut{}, errors.New("Could not get key 'AppName' for one of the shortcuts")
		}

		exe, ok := s.GetString("Exe")
		if !ok {
			return []Shortcut{}, errors.New("Could not get key 'Exe' for one of the shortcuts")
		}

		startDir, ok := s.GetString("StartDir")
		if !ok {
			return []Shortcut{}, errors.New("Could not get key 'StartDir' for one of the shortcuts")
		}

		icon, ok := s.GetString("icon")
		if !ok {
			return []Shortcut{}, errors.New("Could not get key 'icon' for one of the shortcuts")
		}

		isHidden, ok := s.GetBool("IsHidden")
		if !ok {
			return []Shortcut{}, errors.New("Could not get key 'IsHidden' for one of the shortcuts")
		}

		tagsMap, ok := s.GetMap("tags")
		if !ok {
			return []Shortcut{}, errors.New("Could not get key 'tags' for one of the shortcuts")
		}

		tags := []string{}
		for i := 0; i < len(tagsMap); i++ {
			key := strconv.Itoa(i)
			t, ok := tagsMap[key]
			if !ok {
				return []Shortcut{}, errors.New("vdf that should be an array does not have the corresponding index")
			}

			ts, ok := t.AsString()
			if !ok {
				return []Shortcut{}, errors.New("tag should be a string but wasn't")
			}

			tags = append(tags, ts)
		}

		shortcuts[i] = Shortcut{
			AppId:    appId,
			AppName:  appName,
			Exe:      exe,
			Icon:     icon,
			IsHidden: isHidden,
			StartDir: startDir,
			Tags:     tags,
		}
	}

	return shortcuts, nil
}
