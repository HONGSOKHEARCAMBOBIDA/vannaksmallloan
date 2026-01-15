package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hearbong/smallloanbackend/constant/share"
	"github.com/hearbong/smallloanbackend/service"
)

type AccountTypeController struct {
	service service.AccountTypeService
}

func NewAccountTypeController() AccountTypeController {
	return AccountTypeController{
		service: service.NewAccountTypeService(),
	}
}

func (cr AccountTypeController) Get(c *gin.Context) {
	accounttype, err := cr.service.Get()
	if err != nil {
		share.RespondError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, accounttype)
}
