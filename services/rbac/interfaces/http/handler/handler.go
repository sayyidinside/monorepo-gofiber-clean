package handler

type UserManagementHandler struct {
	UserHandler       UserHandler
	PermissionHandler PermissionHandler
	ModuleHandler     ModuleHandler
	RoleHandler       RoleHandler
}

type Handlers struct {
	UserManagementHandler *UserManagementHandler
	AuthHandler           AuthHandler
}
