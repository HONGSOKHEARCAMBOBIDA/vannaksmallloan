package route

const (
	ViewRole         = "/viewrole"
	AddRole          = "/addrole"
	EditRole         = "/editrole/:id"
	ChangeStatusRole = "/changestatusrole/:id"

	ViewUser         = "/viewuser"
	AddUser          = "/adduser"
	EditUser         = "/edituser/:id"
	ChangeStatusUser = "/changestatususer/:id"
	ResetPassword    = "/resetpassword/:id"

	AddPermissionToRole      = "/addpermissiontorole"
	RemovePermissionFromRole = "/removepermissionfromrole"
	ViewRolePermission       = "/viewrolepermission/:id"

	ViewClient         = "viewclient"
	ListClient         = "listclient"
	AddClient          = "addclient"
	EditClient         = "editclient/:id"
	ChangeStatusClient = "changestatusclient/:id"

	ViewLoanProduct = "viewloanproduct"

	ViewDocumentType = "viewdocumenttype"

	ViewLoan           = "viewloan"
	ViewLoanforcheck   = "viewloanforcheck"
	ViewLoanforApprove = "viewloanforapprove"
	AddLoan            = "addloan"
	EditLoan           = "editloan"
	CheckLoan          = "checkloan/:id"
	ApproveLoan        = "approveloan/:id"
	DeleteLoan         = "deleteloan/:id"

	ViewCashierSession = "viewcashiersession"
	AddCashiersSession = "addcashiersession"

	ViewAccountType = "viewaccounttype"

	ViewChartAccount         = "viewchartaccount"
	AddChartAccount          = "addchartaccount"
	EditChartAccount         = "editchartaccount/:id"
	ChangestatusChartAccount = "changestatuschartaccount/:id"
)
