package ethparser

import (
	"strconv"
)

func hexStringToInt(hexString string) (int, error) {
	value, err := strconv.ParseInt(hexString, 16, 64)
	if err != nil {
		return 0, err
	}

	return int(value), nil
}
