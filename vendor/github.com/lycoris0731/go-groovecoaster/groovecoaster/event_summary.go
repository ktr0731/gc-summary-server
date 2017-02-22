package groovecoaster

import "fmt"

// EventResult is your play data in an event
type EventResult struct {
	EP        string `json:"event_point"`
	HighScore int    `json:"high_score"`
	Rank      int
}

// EventSummary is a summary of an event that is being held now
type EventSummary struct {
	EventResult `json:"user_event_data"`
	StartDate   string `json:"open_date"`
	EndDate     string `json:"close_date"`
	Title       string `json:"title_name"`
}

// EventSummary fetches a summary of event
func (c *APIClient) EventSummary() (EventSummary, error) {
	const uri = "mypage.groovecoaster.jp/sp/json/event_data.php"

	data, err := c.get(uri)
	if err != nil {
		return EventSummary{}, err
	}

	var es struct {
		EventSummary `json:"event_data"`
	}
	c.unmarshal(data, &es)

	if es.EventSummary.Title == "" {
		return EventSummary{}, fmt.Errorf("invalid JSON structure")
	}

	return es.EventSummary, nil
}
