package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hearbong/smallloanbackend/constant/share"
	"github.com/hearbong/smallloanbackend/helper"
	"github.com/hearbong/smallloanbackend/request"
	"github.com/hearbong/smallloanbackend/service"
)

type ReceiptController struct {
	service service.ReceiptService
}

func NewReceiptController() ReceiptController {
	return ReceiptController{
		service: service.NewReceiptService(),
	}
}

func (cr ReceiptController) Collectfromgoodloan(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	filters := map[string]string{
		"client_name":  c.Query("client_name"),
		"village_name": c.Query("village_name"),
	}
	collectfromgoodloan, metadata, err := cr.service.Collectfromgoodloan(filters, request.Pagination{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       collectfromgoodloan,
		"pagination": metadata,
	})
}

func (cr ReceiptController) CreateReceipt(c *gin.Context) {
	var input request.ReceiptRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.RespondError(c, http.StatusUnauthorized, "please login")
		return
	}
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.RespondError(c, http.StatusBadRequest, "Invalid id")
		return
	}
	if err := cr.service.CreateReceipt(id, userID, input); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "receipt created")
}
