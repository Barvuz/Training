package structers

// User struct to hold user information
type User struct {
	ID      int    `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"` //should be at least 2 letters
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone" validate:"required"`
	Address string `json:"address" validate:"required"`
}
