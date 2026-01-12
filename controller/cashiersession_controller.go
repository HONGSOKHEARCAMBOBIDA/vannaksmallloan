package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hearbong/smallloanbackend/constant/share"
	"github.com/hearbong/smallloanbackend/helper"
	"github.com/hearbong/smallloanbackend/service"
)

type CashierSessionController struct {
	service service.CashierSessionService
}

func NewCashierSessionController() CashierSessionController {
	return CashierSessionController{
		service: service.NewCashierSessionService(),
	}
}

func (cr CashierSessionController) Create(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.RespondError(c, http.StatusUnauthorized, "Please Login")
		return
	}
	if err := cr.service.Create(userID); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "cashier session created")
}
