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
)
