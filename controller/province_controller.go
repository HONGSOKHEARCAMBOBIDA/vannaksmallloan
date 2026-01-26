package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hearbong/smallloanbackend/constant/share"

	"github.com/hearbong/smallloanbackend/service"
)

type ProvinceController struct {
	service service.ProvinceService
}

func NewProvinceController() ProvinceController {
	return ProvinceController{
		service: service.NewProvinceService(),
	}
}

func (cr ProvinceController) GetProvince(c *gin.Context) {

	province, err := cr.service.GetProvince()
	if err != nil {
		share.RespondError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, province)
}
