package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var Users = []*User{
	{
		Username: "user1",
		Password: "password1",
	},
	{
		Username: "user2",
		Password: "password2",
	},
	{
		Username: "user3",
		Password: "password3",
	},
}

func InitUsers() error {
	// hash users password
	err := HashUsersPassword()
	if err != nil {
		return err
	}

	return nil
}

// change users password to hash
func HashUsersPassword() error {
	for i := 0; i < len(Users); i++ {
		// hash password with bcrypt and error handling
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Users[i].Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// assign hashed password to user password
		Users[i].Password = string(hashedPassword)
	}
	return nil
}
