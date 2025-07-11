package friendship

const (
	statusPending  = "pending"
	statusAccepted = "accepted"
	statusDeclined = "declined"
)

type Friendship struct {
	Username1 string `json:"username1"`
	Username2 string `json:"username2"`
	Status    string `json:"status"`
	Recipient string `json:"recipient"`
}

type FriendList struct {
	Friends []FriendListEntity `json:"usernames"`
	Total   int                `json:"total"`
}

type FriendListEntity struct {
	Username  string `json:"username"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

type FriendshipRequestDTO struct {
	RequesterUsername string
	RecipientUsername string `json:"recipient_username"`
	Status            string `json:"status"`
	Recipient         string `json:"recipient"`
}

func isValidStatus(status string) bool {
	switch status {
	case statusPending, statusDeclined, statusAccepted:
		return true
	}
	return false
}

// user1 must be alphabetically lower than user2
func (req *FriendshipRequestDTO) normalize() {
	if req.RequesterUsername > req.RecipientUsername {
		req.RequesterUsername, req.RecipientUsername = req.RecipientUsername, req.RequesterUsername
	}
}
