package web

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreatePostRequest represents a request to create a post
type CreatePostRequest struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count,omitempty"`
}

// UserResponse represents a user in the response
type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// PostResponse represents a post in the response
type PostResponse struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// PostWithUserResponse represents a post with user information (for joins)
type PostWithUserResponse struct {
	PostID    int    `json:"post_id"`
	PostTitle string `json:"post_title"`
	PostBody  string `json:"post_body"`
	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}
