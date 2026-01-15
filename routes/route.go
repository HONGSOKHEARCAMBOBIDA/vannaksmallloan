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
	clientcontroller := controller.NewClientController()
	loanproductcontroller := controller.NewLoanProductController()
	documenttypecontroller := controller.NewDocumentTypeController()
	loancontroller := controller.NewLaonController()
	cashiersesseioncontroller := controller.NewCashierSessionController()
	accounttypecontroller := controller.NewAccountTypeController()
	chartaccountcontroller := controller.NewChartAccountController()
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
		auth.GET(route.ViewClient, middleware.PermissionMiddleware(permission.ViewClient), clientcontroller.GetAll)
		auth.GET(route.ListClient, middleware.PermissionMiddleware(permission.ViewClient), clientcontroller.GetList)
		auth.POST(route.AddClient, middleware.PermissionMiddleware(permission.AddClient), clientcontroller.Create)
		auth.PUT(route.EditClient, middleware.PermissionMiddleware(permission.EditClient), clientcontroller.Update)
		auth.PUT(route.ChangeStatusClient, middleware.PermissionMiddleware(permission.ChangeStatusClient), clientcontroller.ChangeStatusClient)
		auth.GET(route.ViewLoanProduct, middleware.PermissionMiddleware(permission.ViewLoanProduct), loanproductcontroller.GetAll)
		auth.GET(route.ViewDocumentType, middleware.PermissionMiddleware(permission.ViewDocumentType), documenttypecontroller.GetAll)
		auth.POST(route.AddLoan, middleware.PermissionMiddleware(permission.AddLoan), loancontroller.Create)
		auth.GET(route.ViewLoanforcheck, middleware.PermissionMiddleware(permission.ViewLoan), loancontroller.GetLoanForCheck)
		auth.PUT(route.CheckLoan, middleware.PermissionMiddleware(permission.CheckLoan), loancontroller.CheckLoan)
		auth.GET(route.ViewLoanforApprove, middleware.PermissionMiddleware(permission.ApproveLoan), loancontroller.GetLoanForApprove)
		auth.POST(route.AddCashiersSession, middleware.PermissionMiddleware(permission.AddCashiersSession), cashiersesseioncontroller.Create)
		auth.GET(route.ViewCashierSession, middleware.PermissionMiddleware(permission.ViewCashierSession), cashiersesseioncontroller.Get)
		auth.PUT(route.ApproveLoan, middleware.PermissionMiddleware(permission.ApproveLoan), loancontroller.ApproveLoan)
		auth.DELETE(route.DeleteLoan, middleware.PermissionMiddleware(permission.DeleteLoan), loancontroller.DeleteLoan)
		auth.GET(route.ViewAccountType, middleware.PermissionMiddleware(permission.ViewAccountType), accounttypecontroller.Get)
		auth.POST(route.AddChartAccount, middleware.PermissionMiddleware(permission.AddChartAccount), chartaccountcontroller.Create)
		auth.GET(route.ViewChartAccount, middleware.PermissionMiddleware(permission.ViewChartAccount), chartaccountcontroller.Get)
		auth.PUT(route.EditChartAccount, middleware.PermissionMiddleware(permission.EditChartAccount), chartaccountcontroller.Update)
		auth.PUT(route.ChangestatusChartAccount, middleware.PermissionMiddleware(permission.ChangestatusChartAccount), chartaccountcontroller.ChangeStatusChartAccount)
	}
}
