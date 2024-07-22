package logger

import "fmt"

func Log(v ...any) {
	fmt.Print("[LOG] ")
	fmt.Println(v)
}
func Output(v ...any) {
	fmt.Print("[STDOUT] ")
	fmt.Println(v)
}
