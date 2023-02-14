package permissions

type ListPermissionsRequest struct {
	UserID string `json:"user_id"`
}

type ListPermissionsResponse struct {
	Permissions []string `json:"permissions"`
}
