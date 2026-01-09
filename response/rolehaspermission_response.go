package response

type PermissionWithAssignedRole struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	GroupName   string `json:"permission_group" gorm:"column:permission_group"`
	Sort        int    `json:"sort_order" gorm:"column:sort_order"`
	Assigned    bool   `json:"assigned"`
}
