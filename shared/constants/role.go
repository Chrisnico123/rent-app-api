// shared/constants/roles.go
package constants

const (
	CustomerRole        = 1
	SellerRole          = 2
	AdminRole           = 3
	RoleAdminString     = "admin"
	RoleModeratorString = "moderator"
	RolePublicString    = "public"
	RoleCustomerString  = "customer"
	RoleSellerString    = "seller"
)

var RoleNames = map[int]string{
	CustomerRole: "customer",
	SellerRole:   "seller",
	AdminRole:    "admin",
}

func GetRoleName(roleID int) string {
	if name, exists := RoleNames[roleID]; exists {
		return name
	}
	return "unknown"
}
