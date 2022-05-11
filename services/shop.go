package services

import (
	"BOOK-SHOP/models"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShopServices interface {
	CreateShop(*models.Shop) error
	GetShopByKode(string) (*models.Shop, error)
	GetAll() ([]*models.Shop, error)
	UpdateShop(string, *models.Shop) error
	DeleteShop(string) error
	CreateBookByShop(string, *models.Book) (string, error)
	DeleteBookByShop(string, string) error
	UpdateBookByShop(string, string, *models.Book) error
}

type ShopServiceImpl struct {
	shopcollection *mongo.Collection
	ctx            context.Context
}

func NewShopService(shopcollection *mongo.Collection, ctx context.Context) ShopServices {
	return &ShopServiceImpl{
		shopcollection: shopcollection,
		ctx:            ctx,
	}
}

func (s *ShopServiceImpl) CreateShop(shop *models.Shop) error {
	_, err := s.shopcollection.InsertOne(s.ctx, shop)
	return err
}

func (s *ShopServiceImpl) GetShopByKode(kode string) (*models.Shop, error) {
	var shop *models.Shop
	query := bson.D{bson.E{Key: "shopCode", Value: kode}}
	err := s.shopcollection.FindOne(s.ctx, query).Decode(&shop)
	return shop, err
}

func (s *ShopServiceImpl) GetAll() ([]*models.Shop, error) {
	var shops []*models.Shop
	cursor, err := s.shopcollection.Find(s.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(s.ctx) {
		var shop models.Shop
		err := cursor.Decode(&shop)
		if err != nil {
			return nil, err
		}
		shops = append(shops, &shop)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(s.ctx)

	if len(shops) == 0 {
		return nil, errors.New("documents not found")
	}
	return shops, nil
}

func (s *ShopServiceImpl) UpdateShop(kode string, shop *models.Shop) error {
	filter := bson.D{bson.E{Key: "shopCode", Value: kode}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "name", Value: shop.Name}, bson.E{Key: "address", Value: shop.Address}, bson.E{Key: "city", Value: shop.City}, bson.E{Key: "owner", Value: shop.Owner}, bson.E{Key: "books", Value: shop.Books}}}}
	result, _ := s.shopcollection.UpdateOne(s.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (s *ShopServiceImpl) DeleteShop(kode string) error {
	filter := bson.D{bson.E{Key: "shopCode", Value: kode}}
	result, _ := s.shopcollection.DeleteOne(s.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}
	return nil
}

func (s *ShopServiceImpl) CreateBookByShop(kodeToko string, book *models.Book) (string, error) {
	// Query
	bookCode := uuid.New().String()
	query := bson.M{"shopCode": kodeToko}
	update := bson.M{"$push": bson.M{"books": bson.M{
		"bookCode":    bookCode,
		"title":       book.Title,
		"description": book.Description,
		"rating":      book.Rating,
		"discount":    book.Discount}}}

	// Update
	_, err := s.shopcollection.UpdateOne(s.ctx, query, update)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("Toko tidak ditemukan")
	}

	return bookCode, nil

}

func (s *ShopServiceImpl) DeleteBookByShop(kodeToko string, kodeBuku string) error {
	query := bson.M{"shopCode": kodeToko}
	change := bson.M{"$pull": bson.M{"books": bson.M{"bookCode": kodeBuku}}}
	_, err := s.shopcollection.UpdateOne(s.ctx, query, change)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}
func (s *ShopServiceImpl) UpdateBookByShop(kodeToko string, kodeBuku string, book *models.Book) error {
	query := bson.M{"$and": []interface{}{bson.M{"shopCode": kodeToko, "books.bookCode": kodeBuku}}}

	update := bson.M{
		"$set": bson.M{
			"books.$.title":       book.Title,
			"books.$.description": book.Description,
			"books.$.rating":      book.Rating,
			"books.$.discount":    book.Discount,
		},
	}

	_, err := s.shopcollection.UpdateOne(s.ctx, query, update)
	if err != nil {
		return err
	}
	return nil

}
