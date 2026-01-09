package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hearbong/smallloanbackend/constant/share"
	"github.com/hearbong/smallloanbackend/request"
	"github.com/hearbong/smallloanbackend/service"
)

type RoleController struct {
	service service.RoleService
}

func NewRoleController() RoleController {
	return RoleController{
		service: service.NewRoleService(),
	}
}

func (cr RoleController) GetRole(c *gin.Context) {
	data, err := cr.service.GetRole()
	if err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr RoleController) CreateRole(c *gin.Context) {
	var input request.RoleRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateRole(input); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "Role Created")
}

func (cr RoleController) UpdateRole(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	var input request.RoleRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateRole(id, input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "Role Updated")
}

func (cr RoleController) ChangeStatusRole(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := cr.service.ChangeStatusRole(id); err != nil {
		share.RespondError(c, http.StatusNotFound, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "Role Status Success")
}
