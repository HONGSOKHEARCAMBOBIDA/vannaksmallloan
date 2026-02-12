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
	journalcontroller := controller.NewJournalController()
	receiptcontroller := controller.NewReceiptController()
	paymentschedulecontroller := controller.NewPaymentScheduleController()
	provincecontroller := controller.NewProvinceController()
	districtcontroller := controller.NewDistrictController()
	communcecontroller := controller.NewCommunceController()
	villagecontroller := controller.NewVillageController()
	r.Static("/clientimage", "./public/clientimage")
	r.POST("/login", authcontroller.Login)
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET(route.ViewVillage, middleware.PermissionMiddleware(permission.Viewaddress), villagecontroller.GetVillage)

		auth.GET(route.ViewCommunce, middleware.PermissionMiddleware(permission.Viewaddress), communcecontroller.GetCommunce)

		auth.GET(route.ViewDistrict, middleware.PermissionMiddleware(permission.Viewaddress), districtcontroller.GetDistrict)

		auth.GET(route.ViewProvince, middleware.PermissionMiddleware(permission.Viewaddress), provincecontroller.GetProvince)

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

		auth.GET(route.ViewLoan, middleware.PermissionMiddleware(permission.ViewLoan), loancontroller.GetLoan)

		auth.GET(route.ViewLateLoan,middleware.PermissionMiddleware(permission.ViewLoan),loancontroller.GetLateLoan)

		auth.POST(route.AddLoan, middleware.PermissionMiddleware(permission.AddLoan), loancontroller.Create)

		auth.GET(route.ViewLoanforcheck, middleware.PermissionMiddleware(permission.ViewLoan), loancontroller.GetLoanForCheck)

		auth.PUT(route.CheckLoan, middleware.PermissionMiddleware(permission.CheckLoan), loancontroller.CheckLoan)

		auth.GET(route.ViewLoanforApprove, middleware.PermissionMiddleware(permission.ApproveLoan), loancontroller.GetLoanForApprove)

		auth.DELETE(route.DeleteLoanbeforapprove,middleware.PermissionMiddleware(permission.DeleteLoan),loancontroller.DeleteLoanbeforapprove)

		auth.POST(route.AddCashiersSession, middleware.PermissionMiddleware(permission.AddCashiersSession), cashiersesseioncontroller.Create)

		auth.GET(route.ViewCashierSession, middleware.PermissionMiddleware(permission.ViewCashierSession), cashiersesseioncontroller.Get)

		auth.GET(route.ViewCashierSessionforrollbacke,middleware.PermissionMiddleware(permission.ViewCashierSession),cashiersesseioncontroller.GetforRollback)

		auth.PUT(route.VerifyCashierSession, middleware.PermissionMiddleware(permission.VerifyCashierSession), cashiersesseioncontroller.Verify)

		auth.DELETE(route.RollbackVerify, middleware.PermissionMiddleware(permission.RollbackVerify), cashiersesseioncontroller.RollbackVerify)

		auth.PUT(route.ApproveLoan, middleware.PermissionMiddleware(permission.ApproveLoan), loancontroller.ApproveLoan)

		auth.DELETE(route.DeleteLoan, middleware.PermissionMiddleware(permission.DeleteLoan), loancontroller.DeleteLoan)

		auth.GET(route.ViewAccountType, middleware.PermissionMiddleware(permission.ViewAccountType), accounttypecontroller.Get)

		auth.POST(route.AddChartAccount, middleware.PermissionMiddleware(permission.AddChartAccount), chartaccountcontroller.Create)

		auth.GET(route.ViewChartAccount, middleware.PermissionMiddleware(permission.ViewChartAccount), chartaccountcontroller.Get)

		auth.PUT(route.EditChartAccount, middleware.PermissionMiddleware(permission.EditChartAccount), chartaccountcontroller.Update)

		auth.PUT(route.ChangestatusChartAccount, middleware.PermissionMiddleware(permission.ChangestatusChartAccount), chartaccountcontroller.ChangeStatusChartAccount)

		auth.POST(route.AddJournal, middleware.PermissionMiddleware(permission.AddJournal), journalcontroller.Create)

		auth.GET(route.ViewJournal, middleware.PermissionMiddleware(permission.ViewJournal), journalcontroller.Get)

		auth.PUT(route.EditJournal, middleware.PermissionMiddleware(permission.EditJournal), journalcontroller.Update)

		auth.DELETE(route.DeleteJournal, middleware.PermissionMiddleware(permission.DeleteLoan), journalcontroller.Delete)

		auth.GET(route.ViewReceipt, middleware.PermissionMiddleware(permission.ViewReceipt), receiptcontroller.Collectfromgoodloan)

		auth.POST(route.AddReceipt, middleware.PermissionMiddleware(permission.AddReceipt), receiptcontroller.CreateReceipt)

		auth.DELETE(route.DeleteReceipt, middleware.PermissionMiddleware(permission.DeleteReceipt), receiptcontroller.Delete)

		auth.GET(route.ViewListReceipt, middleware.PermissionMiddleware(permission.ViewReceipt), receiptcontroller.GetReceiptList)

		auth.PUT(route.RemovePenalty, middleware.PermissionMiddleware(permission.RemovePenalty), paymentschedulecontroller.RemovePenalty)

		auth.GET(route.ViewShedule, middleware.PermissionMiddleware(permission.ViewShedule), paymentschedulecontroller.GetFullPaymentSchedule)

		auth.GET(route.ViewBalancesheet, middleware.PermissionMiddleware(permission.ViewJournal), journalcontroller.GetBalanceSheet)

		auth.GET(route.ViewBalancesheetperiod, middleware.PermissionMiddleware(permission.ViewJournal), journalcontroller.GetBalanceSheetForDateRange)

		auth.GET(route.Incomestatment,middleware.PermissionMiddleware(permission.ViewJournal),journalcontroller.Incomestatement)
	}
}
