package groovecoaster

// Player is the structure that represents a part of players in online battle result
type Player struct {
	Name       string `json:"player_name"`
	Prefecture string `json:"pref"`
}

// OnlineBattleSummary is the structure that represents a row of online battle list
type OnlineBattleSummary struct {
	EID        int          `json:",string"`
	MID        int          `json:",string"`
	Star       int          `json:",string"`
	Players    []Player     `json:"detail"`
	IsFavorite StringToBool `json:"fav"`
	Date       string
	Rank       int
}

// OnlineBattleSummary fetches all online battle results
func (c *APIClient) OnlineBattleSummary() ([]OnlineBattleSummary, error) {
	const uri = "mypage.groovecoaster.jp/sp/json/online_battle_list.php"

	data, err := c.get(uri)
	if err != nil {
		return nil, err
	}

	var obs struct {
		OnlineBattleSummary []OnlineBattleSummary `json:"onlineBattleList"`
	}
	c.unmarshal(data, &obs)

	return obs.OnlineBattleSummary, nil
}
