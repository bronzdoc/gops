package util

import (
	"fmt"
)

func LogError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
