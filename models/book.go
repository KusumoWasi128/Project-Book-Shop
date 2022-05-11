package models

type Book struct {
	BookCode    string `json:"book_code" bson:"bookCode"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Rating      int    `json:"rating" bson:"rating"`
	Discount    int    `json:"discount" bson:"discount"`
}
