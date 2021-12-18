package key

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetSecretKey() string {

	err := godotenv.Load("./env/safe.env")
	if err != nil {
		fmt.Println("Error Loading env file")
		fmt.Println(err)
	}
	return os.Getenv("SECRET_KEY")
}
