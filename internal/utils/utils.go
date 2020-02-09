package utils

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func Check(err error){
if err != nil {
	fmt.Println(err)
	panic(err)
}
}
func LoadENV(){
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}