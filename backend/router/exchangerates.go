package router

import (
	"exchangeapp/controllers"

	"github.com/gin-gonic/gin"
)

func registerExchangeRatesRouter(r *gin.RouterGroup) {
	erc := controllers.NewExchangeRateController()

	r.GET("/exchangeRates", erc.GetExchangeRates)
	r.POST("/exchangeRates", erc.CreateExchangeRate)
	r.GET("/exchangeRates/:id", erc.GetExchangeRateByID)
}
