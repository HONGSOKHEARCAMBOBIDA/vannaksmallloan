package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hearbong/smallloanbackend/constant/share"
	"github.com/hearbong/smallloanbackend/service"
)

type LoanProductController struct {
	service service.LoanProductService
}

func NewLoanProductController() LoanProductController {
	return LoanProductController{
		service: service.NewLoanProductService(),
	}
}

func (cr LoanProductController) GetAll(c *gin.Context) {
	loanproduct, err := cr.service.GetAll()
	if err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, loanproduct)
}
