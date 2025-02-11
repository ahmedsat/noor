package noor

import (
	"fmt"
)

func Assert(constrain bool, message string) {

	if !constrain {
		panic(fmt.Sprintf("Assertion failed: %s", message))
	}
}
