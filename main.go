package main

import (
	"labapi/module/restaurant/transport/ginrestautant"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Restaurant struct {
	Id   int    `json:"id " gorm:"column:id;"`
	Name string `json:"name" gorm:"column:name;"`
	Addr string `json:"addr" gorm:"column:addr;"`
}

func (Restaurant) TableName() string { return "restaurants" }

// sử dụng để update name thành tên rỗng.
type RestaurantUpdate struct {
	Name *string `json:"name" gorm:"column:name;"`
	Addr *string `json:"addr" gorm:"column:addr;"`
}

func (RestaurantUpdate) TableName() string { return Restaurant{}.TableName() }

func main() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	dsn := os.Getenv("MYSQL_CONNECTION")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(db, err)

	//GET
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to MySQL-API",
		})
	})

	//POST
	version1 := r.Group("/version1")
	restaurants := version1.Group("/restaurants")
	restaurants.POST("", ginrestautant.CreateRestaurant(db))

	//GET lấy thông tin theo từng ID.
	restaurants.GET("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var data Restaurant

		db.Where("id = ?", id).First(&data)
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	//GET lấy nhiều thông tin
	restaurants.GET("", func(c *gin.Context) {
		var data []Restaurant

		type Paging struct {
			Page  int `json:"page" form:"page"`
			Limit int `json:"limit" form:"limit"`
		}

		var pagingData Paging
		if err := c.ShouldBind(&pagingData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return

		}

		if pagingData.Page <= 0 {
			pagingData.Page = 2
		}

		if pagingData.Limit <= 0 {
			pagingData.Limit = 3
		}
		//c.ShouldBind(pagingData)

		db.Offset((pagingData.Page - 1) * pagingData.Limit).
			Order("id desc").Limit(pagingData.Limit).Find(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	//PUT
	restaurants.PUT("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var data RestaurantUpdate
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return

		}

		db.Where("id = ?", id).Updates(&data)
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	//DELETE
	restaurants.DELETE("/:id", ginrestautant.DeleteRestaurant(db))

	r.Run()

	//Create database. If err gàn buộc lỗi và hiển thị ID mới vừa được nhập.
	/*
		newRestaurant := Restaurant{Name: "Mặt Trời Mọc", Addr: "Lý Thường Kiệt"}
		if err := db.Create(&newRestaurant).Error; err != nil {
			log.Println(err)
		}

		//Truy vấn dữ liệu.
		var myRestaurant Restaurant
		if err := db.Where("id = ?", 6).First(&myRestaurant).Error; err != nil {
			log.Println(err)
		}
		log.Println(myRestaurant)

		//Update dữ liệu.
		newName := "GoLang Vui"
		updateData := RestaurantUpdate{Name: &newName}
		if err := db.Where("id = ?", 4).Updates(&updateData).Error; err != nil {
			log.Println(err)
		}
		log.Println(myRestaurant)

		//Delete dữ liệu. TableName nó là function của Restaurant - cách 2.
		if err := db.Table(Restaurant{}.TableName()).Where("id = ?", 6).Delete(nil).Error; err != nil {
			log.Println(err)
		}
		if err := db.Where("id = ?", 6).Delete(&myRestaurant).Error; err != nil {
			log.Println(err)
		}
		log.Println(myRestaurant)
	*/

}
