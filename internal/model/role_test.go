package model

import "testing"

func TestHasPermission(t *testing.T) {
	tests := []struct {
		role  Role
		perm  string
		allow bool
	}{
		{RoleReader, "chapter:create", false},
		{RoleReader, "series:create", false},

		{RoleUploader, "chapter:create", true},
		{RoleUploader, "chapter:update:own", true},
		{RoleUploader, "chapter:delete:own", true},
		{RoleUploader, "chapter:update", false},
		{RoleUploader, "chapter:delete", false},
		{RoleUploader, "series:create", false},
		{RoleUploader, "series:update", false},

		{RoleModerator, "chapter:create", true},
		{RoleModerator, "chapter:update", true},
		{RoleModerator, "chapter:delete", true},
		{RoleModerator, "series:create", true},
		{RoleModerator, "series:update", true},
		{RoleModerator, "series:delete", true},
		{RoleModerator, "user:list", false},
		{RoleModerator, "user:delete", false},

		{RoleAdmin, "chapter:create", true},
		{RoleAdmin, "chapter:update", true},
		{RoleAdmin, "chapter:delete", true},
		{RoleAdmin, "series:create", true},
		{RoleAdmin, "series:update", true},
		{RoleAdmin, "series:delete", true},
		{RoleAdmin, "user:list", true},
		{RoleAdmin, "user:update", true},
		{RoleAdmin, "user:delete", true},

		{Role("nonexistent"), "chapter:create", false},
	}

	for _, tt := range tests {
		if got := tt.role.HasPermission(tt.perm); got != tt.allow {
			t.Errorf("%s.HasPermission(%q) = %v, want %v", tt.role, tt.perm, got, tt.allow)
		}
	}
}

func TestUserHasPermissionDelegatesToRole(t *testing.T) {
	u := &User{Role: RoleUploader}
	if !u.HasPermission("chapter:create") {
		t.Error("expected uploader user to have chapter:create")
	}
	if u.HasPermission("chapter:delete") {
		t.Error("expected uploader user not to have blanket chapter:delete")
	}
}
