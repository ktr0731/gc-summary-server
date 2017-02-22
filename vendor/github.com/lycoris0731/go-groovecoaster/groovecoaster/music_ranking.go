package groovecoaster

import "fmt"

// Difficulty represents difficulty each of musics
type Difficulty int

// Enum for difficulty
const (
	Simple Difficulty = iota
	Normal
	Hard
	Extra
)

// RankingElement is the structure that represents a ranking detail of a player
type RankingElement struct {
	Score            string `json:"event_point"`
	Name             string `json:"player_name"`
	LastPlayedArcade string `json:"last_play_tenpo_name"`
	Prefecture       string `json:"pref"`
	Rank             int
	Title            string
}

// MusicRankingPageCount fetches last page number by music id and difficulty
func (c *APIClient) MusicRankingPageCount(id int, diff Difficulty) (int, error) {
	const uri = "mypage.groovecoaster.jp/sp/json/score_ranking_bymusic_bydifficulty.php?music_id=%d&difficulty=%d"

	data, err := c.get(fmt.Sprintf(uri, id, diff))
	if err != nil {
		return -1, err
	}

	var rd struct {
		Count   int              `json:",string"`
		Ranking []RankingElement `json:"score_rank"`
	}
	c.unmarshal(data, &rd)

	if rd.Ranking == nil {
		return -1, fmt.Errorf("invalid JSON structure")
	}

	return rd.Count / 100, nil
}

// MusicRanking fetches a music score ranking by music id and difficulty
func (c *APIClient) MusicRanking(id int, diff Difficulty, page int) ([]RankingElement, error) {
	const uri = "mypage.groovecoaster.jp/sp/json/score_ranking_bymusic_bydifficulty.php?music_id=%d&difficulty=%d&page=%d"

	data, err := c.get(fmt.Sprintf(uri, id, diff, page))
	if err != nil {
		return nil, err
	}

	var rd struct {
		Count   int              `json:",string"`
		Ranking []RankingElement `json:"score_rank"`
	}
	c.unmarshal(data, &rd)

	if rd.Ranking == nil {
		return nil, fmt.Errorf("invalid JSON structure")
	}

	return rd.Ranking, nil
}
