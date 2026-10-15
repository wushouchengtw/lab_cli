package models

// Define response struct for check permission CLI.
type CheckPermissionResponse struct {
	Board         string `json:"board"`
	Model         string `json:"model"`
	HasPermission bool   `json:"has_permission"`
}
