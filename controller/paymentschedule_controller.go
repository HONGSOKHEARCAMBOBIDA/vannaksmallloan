package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hearbong/smallloanbackend/constant/share"
	"github.com/hearbong/smallloanbackend/service"
)

type PaymentScheduleController struct {
	service service.PaymentScheduleService
}

func NewPaymentScheduleController() PaymentScheduleController {
	return PaymentScheduleController{
		service: service.NewPaymentScheduleService(),
	}
}

func (cr PaymentScheduleController) RemovePenalty(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.RemovePenalty(id); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "remove penalty success")
}
