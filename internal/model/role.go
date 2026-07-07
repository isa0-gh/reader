package model

type Role string

const (
	RoleReader    Role = "reader"    // read-only access
	RoleUploader  Role = "uploader"  // can upload chapters
	RoleModerator Role = "moderator" // can moderate/edit/delete chapters and series
	RoleAdmin     Role = "admin"     // full access
)

// Permissions each role grants
var RolePermissions = map[Role][]string{
	RoleReader: {
		"comment:create",
	},
	RoleUploader: {
		"chapter:create",
		"chapter:update:own",
		"chapter:delete:own",
		"comment:create",
	},
	RoleModerator: {
		"chapter:create",
		"chapter:update",
		"chapter:delete",
		"series:create",
		"series:update",
		"series:delete",
		"comment:create",
		"comment:delete",
		"comment:suspend",
	},
	RoleAdmin: {
		"chapter:create",
		"chapter:update",
		"chapter:delete",
		"series:create",
		"series:update",
		"series:delete",
		"user:list",
		"user:update",
		"user:delete",
		"comment:create",
		"comment:delete",
		"comment:suspend",
	},
}

func (r Role) HasPermission(perm string) bool {
	for _, p := range RolePermissions[r] {
		if p == perm {
			return true
		}
	}
	return false
}
