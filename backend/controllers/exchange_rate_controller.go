package controllers

import (
	"exchangeapp/dao"
	"exchangeapp/models"
	"exchangeapp/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ExchangeRateController struct {
	exchangeRateService service.ExchangeRateService
}

func NewExchangeRateController() *ExchangeRateController {
	exchangeRateDAO := dao.NewExchangeRateDAO()
	exchangeRateService := service.NewExchangeRateService(exchangeRateDAO)
	return &ExchangeRateController{
		exchangeRateService: exchangeRateService,
	}
}

func (erc *ExchangeRateController) CreateExchangeRate(ctx *gin.Context) {
	var exchangeRate models.ExchangeRate

	if err := ctx.ShouldBindJSON(&exchangeRate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	exchangeRate.Date = time.Now()

	if err := erc.exchangeRateService.CreateExchangeRate(&exchangeRate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, exchangeRate)
}

func (erc *ExchangeRateController) GetExchangeRates(ctx *gin.Context) {
	exchangeRates, err := erc.exchangeRateService.GetExchangeRates()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, exchangeRates)
}

func (erc *ExchangeRateController) GetExchangeRateByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	exchangeRate, err := erc.exchangeRateService.GetExchangeRateByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, exchangeRate)
}
