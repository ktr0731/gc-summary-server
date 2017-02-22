package groovecoaster

import "fmt"

// EventArchive is the structure that represents a row of all events
type EventArchive struct {
	OpenDate  string `json:"open_date"`
	CloseDate string `json:"close_date"`
	EventID   int    `json:"event_id"`
	Title     string
}

// EventArchiveSummary fetches all events that has been held until now
func (c *APIClient) EventArchiveSummary() ([]EventArchive, error) {
	const uri = "mypage.groovecoaster.jp/sp/json/event_info_list.php"

	data, err := c.get(uri)
	if err != nil {
		return nil, err
	}

	var ea struct {
		EventArchiveSummary []EventArchive `json:"event_info_list"`
	}
	c.unmarshal(data, &ea)

	if ea.EventArchiveSummary == nil {
		return nil, fmt.Errorf("invalid JSON structure: EventArchiveSummary()")
	}

	return ea.EventArchiveSummary, nil
}
