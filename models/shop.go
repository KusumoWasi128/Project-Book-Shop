package models

type Shop struct {
	ShopCode string `json:"shop_code" bson:"shopCode"`
	Name     string `json:"name" bson:"name"`
	Address  string `json:"address" bson:"address"`
	City     string `json:"city" bson:"city"`
	Owner    string `json:"owner" bson:"owner"`
	Books    []Book `json:"books" bson:"books"`
}
