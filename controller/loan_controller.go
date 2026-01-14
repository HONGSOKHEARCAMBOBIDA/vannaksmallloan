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

type LoanController struct {
	service service.LoanService
}

func NewLaonController() LoanController {
	return LoanController{
		service: service.NewLoanService(),
	}
}

func (cr LoanController) Create(c *gin.Context) {
	var input request.LoanRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.RespondError(c, http.StatusUnauthorized, "please login")
		return
	}
	err := cr.service.Create(userID, input)
	if err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "Loan request success")
}

func (cr LoanController) GetLoanForCheck(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.RespondError(c, http.StatusUnauthorized, "Please Login")
	}
	loan, err := cr.service.GetLoanForCheck(userID)
	if err != nil {
		share.RespondError(c, http.StatusNoContent, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, loan)
}

func (cr LoanController) GetLoanForApprove(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.RespondError(c, http.StatusUnauthorized, "Please Login")
	}
	loan, err := cr.service.GetLoanForApprove(userID)
	if err != nil {
		share.RespondError(c, http.StatusNoContent, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, loan)
}

func (cr LoanController) CheckLoan(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.RespondError(c, http.StatusNotFound, err.Error())
		return
	}
	if err := cr.service.CheckLoan(id); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "checked success")

}

func (cr LoanController) ApproveLoan(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.RespondError(c, http.StatusNotFound, err.Error())
		return
	}
	if err := cr.service.ApproveLoan(id); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "approved success")

}

func (cr LoanController) DeleteLoan(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.RespondError(c, http.StatusNotFound, err.Error())
		return
	}
	if err := cr.service.DeleteLoan(id); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "deleted success")

}
