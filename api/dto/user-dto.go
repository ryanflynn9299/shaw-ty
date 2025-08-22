package dto

type UserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email,omitempty"`
	CreatedAt string `json:"created_at"`
}

type FullUserResponse struct {
	ID        int64          `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	CreatedAt string         `json:"created_at"`
	Links     []LinkResponse `json:"links,omitempty"`
}

type UpdateUserRequest struct {
	UUID               int64  `json:"uuid,pk,autoincrement"`
	FirstName          string `json:"first_name,omitempty"`
	LastName           string `json:"last_name,omitempty"`
	Email              string `json:"email,omitempty"`
	Username           string `json:"username,omitempty"`
	Password           string `json:"password,omitempty"`
	Salt               string `json:"salt,omitempty"`
	DateJoined         int64  `json:"date_joined,omitempty"`
	DateModified       int64  `json:"date_modified,omitempty"`
	LastPasswordUpdate int64  `json:"last_password_update,omitempty"`
	IsActive           bool   `json:"is_active,omitempty"`
}
