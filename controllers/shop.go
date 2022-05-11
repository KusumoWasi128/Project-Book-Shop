package controllers

import (
	"BOOK-SHOP/models"
	"BOOK-SHOP/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShopController struct {
	ShopServices services.ShopServices
}

func NewShopController(shopservices services.ShopServices) ShopController {
	return ShopController{
		ShopServices: shopservices}

}
func (sc ShopController) CreateShop(c *gin.Context) {
	s := models.Shop{}
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	s.Books = []models.Book{}
	s.ShopCode = uuid.New().String()
	err := sc.ShopServices.CreateShop(&s)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":   "berhasil menanmbah data",
		"status":    "success",
		"kode_toko": s.ShopCode,
	})

}

func (sc *ShopController) GetShop(c *gin.Context) {
	kodeToko := c.Param("kode")
	shop, err := sc.ShopServices.GetShopByKode((kodeToko))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": shop,
	})
}

func (sc *ShopController) GetAll(c *gin.Context) {
	shops, err := sc.ShopServices.GetAll()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": shops,
	})

}

func (sc *ShopController) UpdateShop(c *gin.Context) {
	kodeToko := c.Param("kode")
	var shop models.Shop
	oldShop, _ := sc.ShopServices.GetShopByKode(kodeToko)
	if err := c.ShouldBindJSON(&shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	shop.Books = oldShop.Books
	err := sc.ShopServices.UpdateShop(kodeToko, &shop)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (sc *ShopController) DeleteShop(c *gin.Context) {
	kodeToko := c.Param("kode")
	err := sc.ShopServices.DeleteShop(kodeToko)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (sc *ShopController) CreateBook(c *gin.Context) {
	kodeToko := c.Param("kode")
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	bookCode, err := sc.ShopServices.CreateBookByShop(kodeToko, &book)
	fmt.Println(book)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"message":   "sukses menambah buku",
			"kode_buku": bookCode,
		})
	}

}

func (sc *ShopController) GetBookByShop(c *gin.Context) {
	kodeToko := c.Param("kodeToko")
	kodeBuku := c.Param("kodeBuku")
	shop, err := sc.ShopServices.GetShopByKode(kodeToko)
	books := shop.Books
	var book models.Book
	for _, val := range books {
		if val.BookCode == kodeBuku {
			book = val
			break
		}
	}
	if len(book.BookCode) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Buku tidak ditemukan",
		})
		return

	}
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": book,
	})

}
func (sc *ShopController) GetAllBooksByShop(c *gin.Context) {
	kodeToko := c.Param("kodeToko")
	shop, err := sc.ShopServices.GetShopByKode(kodeToko)
	books := shop.Books
	if len(books) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Buku tidak ditemukan",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": books,
	})

}

func (sc *ShopController) DeleteBook(c *gin.Context) {
	kodeToko := c.Param("kodeToko")
	kodeBuku := c.Param("kodeBuku")
	err := sc.ShopServices.DeleteBookByShop(kodeToko, kodeBuku)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}

func (sc *ShopController) UpdateBook(c *gin.Context) {
	kodeToko := c.Param("kodeToko")
	kodeBuku := c.Param("kodeBuku")

	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := sc.ShopServices.UpdateBookByShop(kodeToko, kodeBuku, &book)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}
