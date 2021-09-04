package helpers

import (
	"fmt"
	"strings"
)

func MergeErrors(array []error) error {

	errorString := make([]string, len(array), cap(array))

	for index, err := range array {
		errorString[index] = err.Error()
	}

	return fmt.Errorf(strings.Join(errorString, "\n"))
}
