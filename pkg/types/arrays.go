package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Int64Array []int64

func (a *Int64Array) Scan(value interface{}) error {
	bs, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	if string(bs) == "{}" {
		*a = []int64{}
		return nil
	}

	if err := json.Unmarshal(bs, a); err != nil {
		return err
	}
	return nil
}

func (a Int64Array) Value() (driver.Value, error) {
	return json.Marshal(a)
}
