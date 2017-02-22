package groovecoaster

// Sales is the structure that represents shop sales info
type Sales struct {
	Avatar  IntToBool `json:"avatar_sale"`
	Item    IntToBool `json:"item_sale"`
	Message IntToBool `json:"message_sale"`
	Music   IntToBool `json:"music_sale"`
	Skin    IntToBool `json:"skin_sale"`
	Sound   IntToBool `json:"sound_sale"`
}

// Shop is the structure that contains sales info, whether is openning
type Shop struct {
	Sales  Sales `json:"shop_data"`
	IsOpen bool  `json:"open_flg"`
	Coin   int   `json:"current_coin"`
}

// ShopSummary fetches shop information
func (c *APIClient) ShopSummary() (Shop, error) {
	const uri = "mypage.groovecoaster.jp/sp/json/shop_sales_data.php"

	data, err := c.get(uri)
	if err != nil {
		return Shop{}, err
	}

	var s Shop
	c.unmarshal(data, &s)

	return s, nil
}
