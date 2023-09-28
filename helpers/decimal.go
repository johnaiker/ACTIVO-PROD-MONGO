package helpers

import (
	"gopkg.in/inf.v0"
)

func ToDecimal(s string) (*inf.Dec, error) {
	d := new(inf.Dec)

	res, _ := d.SetString(s)
	return res, nil
}
