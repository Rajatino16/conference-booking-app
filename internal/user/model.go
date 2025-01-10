package user

type User struct {
	ID string
}

type AddUserRequest struct {
	ID string `json:"id"`
}
