package main

import (
	"BOOK-SHOP/controllers"
	"BOOK-SHOP/services"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connection URI

const uri = "mongodb+srv://wasi:gampangajalah@cluster0.gubot.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

var (
	server         *gin.Engine
	ctx            context.Context
	shopservice    services.ShopServices
	ShopController controllers.ShopController
	shopcollection *mongo.Collection
	mongoclient    *mongo.Client
	err            error
)

func init() {
	ctx = context.TODO()
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)

	mongoconn := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)
	mongoclient, err := mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
	}
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("mongo connection established")

	shopcollection = mongoclient.Database("book-shop").Collection("shops")
	shopservice = services.NewShopService(shopcollection, ctx)
	ShopController = controllers.NewShopController(shopservice)
	server = gin.Default()

}
func main() {
	defer mongoclient.Disconnect(ctx)

	shoproute := server.Group("/book-shop")
	shoproute.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server Berjalan dengan baik",
			"status":  "OK",
		})
	})
	shoproute.POST("/shops", ShopController.CreateShop)
	shoproute.GET("/shops", ShopController.GetAll)
	shoproute.GET("/shops/:kode", ShopController.GetShop)
	shoproute.PUT("/shops/:kode", ShopController.UpdateShop)
	shoproute.DELETE("/shops/:kode", ShopController.DeleteShop)

	shoproute.POST("/shops/book/:kode", ShopController.CreateBook)
	shoproute.GET("/shops/book/:kodeToko", ShopController.GetAllBooksByShop)
	shoproute.GET("shops/book/:kodeToko/:kodeBuku", ShopController.GetBookByShop)
	shoproute.DELETE("shops/book/:kodeToko/:kodeBuku", ShopController.DeleteBook)
	shoproute.PUT("shops/book/:kodeToko/:kodeBuku", ShopController.UpdateBook)

	log.Fatal(server.Run(":8000"))
}
