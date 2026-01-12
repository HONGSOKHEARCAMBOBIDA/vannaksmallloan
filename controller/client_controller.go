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

type ClientController struct {
	service service.ClientService
}

func NewClientController() ClientController {
	return ClientController{
		service: service.NewClientService(),
	}
}

func (cr ClientController) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	filter := map[string]string{
		"name": c.Query("name"),
	}
	clients, metadata, err := cr.service.GetList(filter, request.Pagination{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       clients,
		"pagination": metadata,
	})

}

func (cr ClientController) GetAll(c *gin.Context) {
	client, err := cr.service.GetAll()
	if err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, client)

}

func (cr ClientController) Create(c *gin.Context) {
	var input request.ClientRequestCreate
	if err := c.ShouldBind(&input); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.RespondError(c, http.StatusUnauthorized, "Please Login")
		return
	}
	err := cr.service.Create(input, c, userID)
	if err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "Client has Created")
}

func (cr ClientController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	var input request.ClientRequestUpdate
	if err := c.ShouldBind(&input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.RespondError(c, http.StatusUnauthorized, "please login")
		return
	}
	if err := cr.service.Update(id, input, c, userID); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "Client has updated")
}

func (cr ClientController) ChangeStatusClient(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := cr.service.ChangeStatusClient(id); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "Status has update")
}
