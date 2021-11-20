package util

import "strconv"

func Convert(v string) interface{} {
	if i, err := strconv.Atoi(v); err == nil {
		// No error, this is an int
		return i
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		// no error, this is a float
		return f
	}
	if v == "false" || v == "true" {
		// only parse true or false into bools
		return v == "true"
	}
	return v
}
