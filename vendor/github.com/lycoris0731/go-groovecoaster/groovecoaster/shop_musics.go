package groovecoaster

// ShopMusics is the structure that represents shop sales info about musics
type ShopMusics struct {
	ProductID string `json:"product_id"`
	GC        int    `json:",string"`
	Artist    string
	ID        int
	Name      string
	New       StringToBool
}

// ShopMusics fetches all musics in shop
func (c *APIClient) ShopMusics() ([]ShopMusics, error) {
	const uri = "mypage.groovecoaster.jp/sp/json/shop_music_list.php"

	data, err := c.get(uri)
	if err != nil {
		return nil, err
	}

	var sm struct {
		ShopMusics []ShopMusics `json:"product_list"`
	}
	c.unmarshal(data, &sm)

	return sm.ShopMusics, nil
}
