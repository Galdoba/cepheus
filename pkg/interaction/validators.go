package interaction

import (
	"errors"
	"strconv"
)

var Number = func(str string) error {
	_, err := strconv.Atoi(str)
	if err != nil {
		return errors.New("input must be integer")
	}
	return nil
}
