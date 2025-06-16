package siplib

import "fmt"
import "math/rand"

//Generates a random string of determinate length.
//Don't forget to call rand.Seed(time.Now().UnixNano()) before using it.

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numbers = []rune("0123456789")

func random_string(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func random_number_string(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(b)
}

//Generate a UUID
//Don't forget to call rand.Seed(time.Now().UnixNano()) before using it.

func GenerateUUID(length int) []byte {
	token := make([]byte, 16)
	rand.Read(token)
	return token
}

func GenerateHexUUID(length int) string {
	uuid := GenerateUUID(length)
	return fmt.Sprintf("%x", uuid)
}
