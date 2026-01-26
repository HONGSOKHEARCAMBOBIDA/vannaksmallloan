package controller

import (
	"net/http"
	"strconv"
	"time"

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

func (cr JournalController) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.RespondError(c, http.StatusUnauthorized, "please login")
		return
	}
	if err := cr.service.Delete(id); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "Journal deleted")
}

func (cr JournalController) GetBalanceSheet(c *gin.Context) {
	data := c.Query("date")
	if data == "" {
		data = time.Now().Format("2006-01-02")
	}
	if _, err := time.Parse("2006-01-02", data); err != nil {
		share.RespondError(c, http.StatusBadRequest, "invalid data format")
		return
	}
	balancesheet, err := cr.service.GenerateBalanceSheet(data)
	if err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !balancesheet.Totals.IsBalanced {
		balancesheet.Message = "warning: Balance sheet not balance!"
	}
	share.RespondDate(c, http.StatusOK, balancesheet)
}

func (cr JournalController) GetBalanceSheetForDateRange(c *gin.Context) {
	startDate := c.Query("start")
	endDate := c.Query("end")
	if startDate == "" || endDate == "" {
		share.RespondError(c, http.StatusBadRequest, "start and end required")
		return
	}
	balancesheet, err := cr.service.GenerateBalanceSheet(endDate)
	if err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	balancesheet.ReportTitle = "BALANCE SHEET - PERIOD ENDING"
	share.RespondDate(c, http.StatusOK, balancesheet)
}
