package ginrestautant

import (
	resbusiness "labapi/module/restaurant/business"
	resstorage "labapi/module/restaurant/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Func CreateRestaurant sẽ thay thế vào func trong code API POST.
// sau khi đưa CreateRestaurant vào thì code POST như sau: restaurants.POST("", ginrestautant.CreateRestaurant(db))
func DeleteRestaurant(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		store := resstorage.NewSQLStore(db)
		bus := resbusiness.NewDeleteRestaurantBus(store)

		if err := bus.DeleteRestaurant(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": "Data delete successfully",
		})

	}
}
