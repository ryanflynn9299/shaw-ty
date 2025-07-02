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
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
}
