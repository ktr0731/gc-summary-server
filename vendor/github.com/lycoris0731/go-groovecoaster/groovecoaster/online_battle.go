package groovecoaster

import "fmt"

// OnlineBattleDetail is the structure that represents online battle result
type OnlineBattleDetail struct {
	Difficulty1 StringToDifficulty `json:"difficulty_1st"`
	Difficulty2 StringToDifficulty `json:"difficulty_2nd"`
	Difficulty3 StringToDifficulty `json:"difficulty_3rd"`
	EntryNo     string             `json:"entry_no"`
	Item1       string             `json:"item_1st"`
	Item2       string             `json:"item_2nd"`
	Item3       string             `json:"item_3rd"`
	Music1      string             `json:"music_1st"`
	Music2      string             `json:"music_2nd"`
	Music3      string             `json:"music_3rd"`
	Name        string             `json:"player_name"`
	Prefecture  string             `json:"pref"`
	Star        int                `json:"result_star,string"`
	Arcade      string             `json:"tenpo_name"`
	Rank        int
	Title       string
	Avatar      string
}

// OnlineBattle fetches a online battle detail by eid and mid
func (c *APIClient) OnlineBattle(eid int, mid int) ([]OnlineBattleDetail, error) {
	const uri = "mypage.groovecoaster.jp/sp/json/online_battle_detail.php?eid=%d&mid=%d"

	data, err := c.get(fmt.Sprintf(uri, eid, mid))
	if err != nil {
		return nil, err
	}

	var obd struct {
		OnlineBattle []OnlineBattleDetail
	}
	c.unmarshal(data, &obd)

	return obd.OnlineBattle, nil
}
