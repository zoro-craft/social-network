package users

type UserResult struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
}