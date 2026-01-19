package users

type userRepo struct{}

func NewUserRepo() *userRepo {
	return &userRepo{}
}

func (r *userRepo) CreateUser(user *User) error {
	// Implementation for creating a user in the database
	return nil
}

func (r *userRepo) GetUserByEmail(email string) (*User, error) {
	// Implementation for retrieving a user by email from the database
	return nil, nil
}
