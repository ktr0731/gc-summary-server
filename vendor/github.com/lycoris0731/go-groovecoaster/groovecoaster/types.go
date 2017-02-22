package groovecoaster

import (
	"fmt"
	"strconv"
)

// IntToBool type is a type that regards int as bool for JSON
type IntToBool bool

// UnmarshalJSON unmarshals int to bool
func (i *IntToBool) UnmarshalJSON(bytes []byte) error {
	// If called from receiver that already marshaled from int to bool, bytes indicates bool.
	// So has to check it before strconv.
	switch string(bytes) {
	case "true":
		*i = true
		return nil
	case "false":
		*i = false
		return nil
	}

	n, err := strconv.Atoi(string(bytes))
	if err != nil {
		return err
	}

	switch n {
	case 1:
		*i = true
	case 0:
		*i = false
	default:
		return fmt.Errorf("Invalid value in IntToBool")
	}

	return nil
}

// StringToBool type is a type that regards string as bool for JSON
type StringToBool bool

// UnmarshalJSON unmarshals string to bool
func (s *StringToBool) UnmarshalJSON(bytes []byte) error {
	switch string(bytes) {
	case `"1"`:
		*s = true
	case `"0"`:
		*s = false
	default:
		return fmt.Errorf("Invalid value in StringToBool")
	}

	return nil
}

// StringToDifficulty type is a type that regards string as difficulty in GrooveCoaster for JSON
type StringToDifficulty string

// UnmarshalJSON unmarshals string to difficulty in GrooveCoaster
func (s *StringToDifficulty) UnmarshalJSON(bytes []byte) error {
	switch string(bytes) {
	case `"1"`:
		*s = "SIMPLE"
	case `"2"`:
		*s = "NORMAL"
	case `"3"`:
		*s = "HARD"
	case `"4"`:
		*s = "EXTRA"
	default:
		return fmt.Errorf("Invalid value in StringToDifficulty")
	}

	return nil
}

// AvatarAward is the type that represents avatar award of each event
type AvatarAward string

// ItemAward is the type that represents item award of each event
type ItemAward struct {
	Name   string `json:"item_name"`
	Number int    `json:"item_num"`
}

// MusicAward is the type that represents music award of each event
type MusicAward string

// TitleAward is the type that represents title award of each event
type TitleAward string

// Awards is the type that represents set of each awards
type Awards struct {
	Avatars []AvatarAward `json:"avatar_award"`
	Items   []ItemAward   `json:"item_award"`
	Musics  []MusicAward  `json:"music_award"`
	Titles  []TitleAward  `json:"title_award"`
	Trophy  int           `json:"trophy_num"`
}
