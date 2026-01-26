package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hearbong/smallloanbackend/constant/share"
	"github.com/hearbong/smallloanbackend/service"
)

type CommunceController struct {
	service service.CommunceService
}

func NewCommunceController() CommunceController {
	return CommunceController{
		service: service.NewCommunceService(),
	}
}

func (cr CommunceController) GetCommunce(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.RespondError(c, http.StatusNotFound, err.Error())
		return
	}
	communce, err := cr.service.GetCommunce(id)
	if err != nil {
		share.RespondError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, communce)
}
