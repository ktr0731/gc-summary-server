package groovecoaster

import "fmt"

// Messages is the structure that represents shop sales info about messages
type Messages struct {
	GC        int `json:",string"`
	ID        int
	Name      string
	New       StringToBool
	ProductID string `json:"product_id"`
}

// ShopMessages is the structure that is a set of each categories
type ShopMessages struct {
	Communication []Messages
	AA            []Messages
	Item          []Messages
	Music         []Messages
	Character     []Messages
}

// ShopMessages fetches all musics in shop
func (c *APIClient) ShopMessages() (ShopMessages, error) {
	type shopMessages struct {
		ShopMessages []Messages `json:"product_list"`
	}

	const uri = "mypage.groovecoaster.jp/sp/#/sp_me/%d"
	const (
		Communication = iota
		AA
		Item
		Music
		Character
	)

	var sm ShopMessages
	var messages = make([][]Messages, 5)
	for i := Communication; i <= Character; i++ {
		var message shopMessages
		data, err := c.get(fmt.Sprintf(uri, i+1))
		if err != nil {
			return ShopMessages{}, err
		}

		c.unmarshal(data, &message)
		messages[i] = message.ShopMessages
	}

	sm.Communication = messages[Communication]
	sm.AA = messages[AA]
	sm.Character = messages[Character]
	sm.Item = messages[Item]
	sm.Music = messages[Music]

	return sm, nil
}
