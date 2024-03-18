package pkg

import (
	"fmt"
	"time"
)

func DateValidation(s string) error {
	_, err := time.Parse("02.01.2006", s)
	if err != nil {
		return fmt.Errorf("%s is not a valid date\n", s)
	}
	return nil
}
