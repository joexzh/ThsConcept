package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

func Scan[T any](dest T, src interface{}) error {
	var source []byte
	switch src.(type) {
	case []uint8:
		source = src.([]uint8)
	case nil:
		return nil
	default:
		return fmt.Errorf("incompatible type for %t", dest)
	}
	err := json.Unmarshal(source, &dest)
	if err != nil {
		return err
	}
	return nil
}

func Value[T any](src T) (driver.Value, error) {
	j, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	return driver.Value(j), nil
}
