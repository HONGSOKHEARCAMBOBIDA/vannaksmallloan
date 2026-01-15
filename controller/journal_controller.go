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

type JournalController struct {
	service service.JournalService
}

func NewJournalController() JournalController {
	return JournalController{
		service: service.NewJournalService(),
	}
}

func (cr JournalController) Create(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.RespondError(c, http.StatusUnauthorized, "Please Login")
		return
	}
	var input request.JournalRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.Create(userID, input); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "Journal Created")
}

func (cr JournalController) Get(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	filters := map[string]string{
		"reference_code": c.Query("reference_code"),
		"between":        c.Query("between"),
	}
	journals, metadata, err := cr.service.Get(filters, request.Pagination{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       journals,
		"pagination": metadata,
	})
}

func (cr JournalController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.RespondError(c, http.StatusUnauthorized, "please login")
		return
	}
	var input request.JournalRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.Update(id, input); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "journal updated")
}
