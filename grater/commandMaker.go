package grater

import (
	"fmt"
	"math/rand"
)

func MakeFileTransferCommand() string {
	randomNumber := rand.Intn(1500) + 1
	return fmt.Sprintf("file-%d.txt", randomNumber)
}
