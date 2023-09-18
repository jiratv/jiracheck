package ginrestautant

import (
	resbusiness "labapi/module/restaurant/business"
	resmodel "labapi/module/restaurant/model"
	resstorage "labapi/module/restaurant/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Func CreateRestaurant sẽ thay thế vào func trong code API POST.
// sau khi đưa CreateRestaurant vào thì code POST như sau: restaurants.POST("", ginrestautant.CreateRestaurant(db))
func CreateRestaurant(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var data resmodel.RestaurantCreate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return

		}

		store := resstorage.NewSQLStore(db)
		bus := resbusiness.NewCreateRestaurantBus(store)

		if err := bus.CreateRestaurant(c.Request.Context(), &data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": "Data input completed successfully",
		})

	}
}
