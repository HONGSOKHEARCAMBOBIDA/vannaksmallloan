package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hearbong/smallloanbackend/constant/permission"
	"github.com/hearbong/smallloanbackend/constant/route"
	"github.com/hearbong/smallloanbackend/controller"
	"github.com/hearbong/smallloanbackend/middleware"
)

func SetupRoutes(r *gin.Engine) {
	authcontroller := controller.NewAuthController()
	rolecontroller := controller.NewRoleController()
	usercontroller := controller.NewUserController()
	rolepermissioncontroller := controller.NewRoleHasPermissionController()
	r.Static("/clientimage", "./public/clientimage")
	r.POST("/login", authcontroller.Login)
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET(route.ViewRole, middleware.PermissionMiddleware(permission.ViewRole), rolecontroller.GetRole)
		auth.POST(route.AddRole, middleware.PermissionMiddleware(permission.AddRole), rolecontroller.CreateRole)
		auth.PUT(route.EditRole, middleware.PermissionMiddleware(permission.EditRole), rolecontroller.UpdateRole)
		auth.PUT(route.ChangeStatusRole, middleware.PermissionMiddleware(permission.ChangeStatusRole), rolecontroller.ChangeStatusRole)
		auth.GET(route.ViewUser, middleware.PermissionMiddleware(permission.ViewUser), usercontroller.GetUser)
		auth.POST(route.AddUser, middleware.PermissionMiddleware(permission.AddUser), usercontroller.Register)
		auth.PUT(route.EditUser, middleware.PermissionMiddleware(permission.EditUser), usercontroller.Update)
		auth.PUT(route.ChangeStatusUser, middleware.PermissionMiddleware(permission.EditUser), usercontroller.ChangeStatusUser)
		auth.PUT(route.ResetPassword, middleware.PermissionMiddleware(permission.ResetPassword), usercontroller.ResetPassword)
		auth.POST(route.AddPermissionToRole, middleware.PermissionMiddleware(permission.AddPermissionToRole), rolepermissioncontroller.CreateRolePermissions)
		auth.DELETE(route.RemovePermissionFromRole, middleware.PermissionMiddleware(permission.RemovePermissionFromRole), rolepermissioncontroller.DeleteRolePermission)
		auth.GET(route.ViewRolePermission, middleware.PermissionMiddleware(permission.ViewRolePermission), rolepermissioncontroller.GetRolePermission)
	}
}
