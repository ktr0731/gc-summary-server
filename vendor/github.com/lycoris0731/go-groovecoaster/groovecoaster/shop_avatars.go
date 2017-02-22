package groovecoaster

// ShopAvatars is the structure that represents shop sales info about avatars
type ShopAvatars struct {
	GC        int `json:",string"`
	ID        int
	Name      string
	New       StringToBool
	ProductID string `json:"product_id"`
}

// ShopAvatars fetches all avatars in shop
func (c *APIClient) ShopAvatars() ([]ShopAvatars, error) {
	const uri = "mypage.groovecoaster.jp/sp/json/shop_avatar_list.php"

	data, err := c.get(uri)
	if err != nil {
		return nil, err
	}

	var sa struct {
		ShopAvatars []ShopAvatars `json:"product_list"`
	}
	c.unmarshal(data, &sa)

	return sa.ShopAvatars, nil
}
