package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var JWT_SECRET string

func InitConfig(path string) error {

	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	JWT_SECRET = os.Getenv("JWT_SECRET")

	fmt.Println(JWT_SECRET)

	return nil

}
