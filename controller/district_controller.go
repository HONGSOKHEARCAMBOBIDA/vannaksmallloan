package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hearbong/smallloanbackend/constant/share"
	"github.com/hearbong/smallloanbackend/service"
)

type DistrictController struct {
	service service.DistrictService
}

func NewDistrictController() DistrictController {
	return DistrictController{
		service: service.NewDistrictService(),
	}
}

func (cr DistrictController) GetDistrict(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.RespondError(c, http.StatusNotFound, err.Error())
		return
	}
	district, err := cr.service.GetDistrict(id)
	if err != nil {
		share.RespondError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, district)
}
