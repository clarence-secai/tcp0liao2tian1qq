package utils

import (
	"fmt"
	"os"
)

func ErrExit(where string, err error) {
	if err != nil {
		fmt.Println(where, err)
		os.Exit(1)
	}
}
func ErrContinue(where string, err error) {
	if err != nil {
		fmt.Println(where, err)
	}
}
