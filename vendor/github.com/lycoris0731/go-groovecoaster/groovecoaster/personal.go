package groovecoaster

import "fmt"

// Statistics is the structure that represents stage statistics of all musics
type Statistics struct {
	All       int
	Clear     int
	FullChain int
	NoMiss    int
	Perfect   int
	S         int
	SS        int
	SSS       int
}

// Personal is the structure that represents personal configration
type Personal struct {
	AverageScore      string `json:"average_score"`
	Avatar            string
	FriendApplication bool
	Level             int
	Name              string `json:"player_name"`
	Rank              int
	Title             string
	TotalMusic        int    `json:"total_music"`
	TotalPlayMusic    int    `json:"total_play_music"`
	TotalScore        string `json:"total_score"`
	TotalTrophy       string `json:"total_trophy"`
	TrophyRank        string `json:"trophy_rank"`
}

// Personal fetches player profile
func (c *APIClient) Personal() (Personal, error) {
	const uri = "mypage.groovecoaster.jp/sp/json/player_data.php"

	data, err := c.get(uri)
	if err != nil {
		return Personal{}, err
	}

	var pd struct {
		Personal   Personal   `json:"player_data"`
		Statistics Statistics `json:"stage"`
	}
	c.unmarshal(data, &pd)

	if pd.Personal.Title == "" {
		return Personal{}, fmt.Errorf("Invalid JSON structure")
	}

	return pd.Personal, nil
}

// Statistics fetches music statistics
func (c *APIClient) Statistics() (Statistics, error) {
	const uri = "mypage.groovecoaster.jp/sp/json/player_data.php"

	data, err := c.get(uri)
	if err != nil {
		return Statistics{}, err
	}

	var pd struct {
		Personal   Personal   `json:"player_data"`
		Statistics Statistics `json:"stage"`
	}
	c.unmarshal(data, &pd)

	if pd.Personal.Title != "" && pd.Statistics.All == 0 {
		return Statistics{}, fmt.Errorf("Invalid JSON structure")
	}

	return pd.Statistics, nil
}
