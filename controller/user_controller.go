package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hearbong/smallloanbackend/constant/share"
	"github.com/hearbong/smallloanbackend/request"
	"github.com/hearbong/smallloanbackend/service"
)

type UserController struct {
	service service.UserService
}

func NewUserController() UserController {
	return UserController{
		service: service.NewUserService(),
	}
}

func (cr UserController) Register(c *gin.Context) {
	var input request.UserRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.Register(input); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "User Register Success")
}

func (cr UserController) GetUser(c *gin.Context) {
	filter := map[string]string{
		"name":      c.Query("name"),
		"role_id":   c.Query("role_id"),
		"is_active": c.Query("is_active"),
	}
	user, err := cr.service.GetUser(filter)
	if err != nil {
		share.RespondError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, user)
}

func (cr UserController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.RespondError(c, http.StatusBadRequest, "Invalid Id")
		return
	}
	var input request.UserRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.Update(id, input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "User Update")
}

func (cr UserController) ChangeStatusUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.RespondError(c, http.StatusBadRequest, "Invalid Id")
		return
	}
	if err := cr.service.ChangeStatusUser(id); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "User Staus Success")
}

func (cr UserController) ResetPassword(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var input request.UserRequestResetPassword
	if err := c.ShouldBindJSON(&input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := cr.service.ResetPassword(id, input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	share.ResponeSuccess(c, http.StatusOK, "Password reset successfully")
}
