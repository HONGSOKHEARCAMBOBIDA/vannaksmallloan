package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hearbong/smallloanbackend/constant/share"
	"github.com/hearbong/smallloanbackend/request"
	"github.com/hearbong/smallloanbackend/service"
)

type ChartAccountController struct {
	service service.ChartAccountService
}

func NewChartAccountController() ChartAccountController {
	return ChartAccountController{
		service: service.NewChartAccountService(),
	}
}

func (cr ChartAccountController) Create(c *gin.Context) {
	var input request.ChartAccountRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.Create(input); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "chartaccount created")
}

func (cr ChartAccountController) Get(c *gin.Context) {
	data, err := cr.service.Get()
	if err != nil {
		share.RespondError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr ChartAccountController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var input request.ChartAccountRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.Update(id, input); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "Chartaccount Update")
}

func (cr ChartAccountController) ChangeStatusChartAccount(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.RespondError(c, http.StatusBadRequest, "Invalid id")
		return
	}
	if err := cr.service.ChangeStatusChartAccount(id); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "Chartaccount Update status")
}
