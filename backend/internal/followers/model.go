package followers

type Follow struct {
	FollowerID int    `json:"follower_id"`
	FolloweeID int    `json:"followee_id"`
	Status     string `json:"status"`

	FirstName string `json:"first_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
}