package controller

import (
	"net/http"
	// "strconv"

	"github.com/gin-gonic/gin"
	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/constant/share"
	"github.com/hearbong/smallloanbackend/request"
	"github.com/hearbong/smallloanbackend/response"
	"github.com/hearbong/smallloanbackend/service"
)

type RoleHasPermissionController struct {
	service service.RolePermissionService
}

func NewRoleHasPermissionController() RoleHasPermissionController {
	return RoleHasPermissionController{
		service: service.NewRolePermissionService(),
	}
}

func (cr RoleHasPermissionController) CreateRolePermissions(c *gin.Context) {
	var input request.CreateRolePermissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateRolePermissions(input); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "Permission Assigned to role success")
}

func (cr RoleHasPermissionController) DeleteRolePermission(c *gin.Context) {
	var input request.DeleteRolePermissionsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		share.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.DeleteRolePermission(input); err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponeSuccess(c, http.StatusOK, "Permission Remove from role success")
}

// func (cr RoleHasPermissionController) GetRolePermission(c *gin.Context) {
// 	roleidparam := c.Param("id")
// 	roleid, err := strconv.Atoi(roleidparam)
// 	if err != nil {
// 		share.RespondError(c, http.StatusBadRequest, "Ivalid role id provide")
// 		return
// 	}
// 	data, err := cr.service.GetRolePermission(roleid)
// 	if err != nil {
// 		share.RespondError(c, http.StatusNotFound, err.Error())
// 		return
// 	}
// 	share.RespondDate(c, http.StatusOK, data)
// }
func (cr RoleHasPermissionController) GetRolePermission(c *gin.Context) {
	id := c.Param("id")
	var permissions []response.PermissionWithAssignedRole

	err := config.DB.Table("permissions").
		Select(`
            permissions.id AS id,
            permissions.name AS name,
            permissions.display_name AS display_name,
            CASE 
                WHEN role_has_permissions.permission_id IS NULL THEN false 
                ELSE true 
            END AS assigned
        `).
		Joins(`
            LEFT JOIN role_has_permissions 
            ON permissions.id = role_has_permissions.permission_id 
            AND role_has_permissions.role_id = ?
        `, id).
		Scan(&permissions).Error

	if err != nil {
		share.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	share.RespondDate(c, http.StatusOK, permissions)
}