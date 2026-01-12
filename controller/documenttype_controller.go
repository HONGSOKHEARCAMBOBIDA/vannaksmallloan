package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hearbong/smallloanbackend/constant/share"
	"github.com/hearbong/smallloanbackend/service"
)

type DocumentTypeController struct {
	service service.DocumentTypeService
}

func NewDocumentTypeController() DocumentTypeController {
	return DocumentTypeController{
		service: service.NewDocumentTypeService(),
	}
}

func (cr DocumentTypeController) GetAll(c *gin.Context) {
	documenttype, err := cr.service.GetAll()
	if err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, documenttype)
}
