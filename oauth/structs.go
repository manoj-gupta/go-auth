package oauth

// User contains common information about user
type User struct {
	Any         map[string]interface{}
	Provider    string
	UserID      string
	Email       string
	Name        string
	Location    string
	AvatarURL   string
	Description string
}
