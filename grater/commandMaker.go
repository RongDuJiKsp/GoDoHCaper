package grater

import (
	"fmt"
	"math/rand"
)

func MakeFileTransferCommand() string {
	randomNumber := rand.Intn(1500) + 1
	return fmt.Sprintf("download ./random_files/file-%d.txt", randomNumber)
}
