package dto

type LinkResponse struct {
	ID        int64  `json:"id"`
	ShortCode string `json:"short_code"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	Active    bool   `json:"active"`
}

type FullLinkResponse struct {
	ID        int64  `json:"id"`
	ShortCode string `json:"short_code"`
	URL       string `json:"url"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
	Clicks    int    `json:"clicks"`
	Active    bool   `json:"active"`
}

type CreateLinkRequest struct {
	URL       string  `json:"url" binding:"required"`
	ShortCode *string `json:"short_code,omitempty"`
}

type UpdateLinkRequest struct {
	URL       *string `json:"url,omitempty"`
	ShortCode *string `json:"short_code,omitempty"`
	Active    *bool   `json:"active,omitempty"`
}
