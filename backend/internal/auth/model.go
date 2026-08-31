package auth

type User struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"-"`
	DateOfBirth string `json:"date_of_birth"`

	Nickname string `json:"nickname,omitempty"`
	AboutMe  string `json:"about_me,omitempty"`
	Avatar   string `json:"avatar,omitempty"`

	CreatedAt string `json:"created_at,omitempty"`
}
