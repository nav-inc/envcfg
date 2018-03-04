package envcfg

import (
	"strconv"
	"time"
)

// The DefaultConverters are functions for converting strings into basic built-in Go types.
// A few types have been omitted: rune (use int32), byte (use uint8), uintptr, complex64, and
// complex128
var DefaultConverters = []interface{}{
	convertBool,
	convertString,
	convertInt,
	convertFloat32,
	convertFloat64,
	convertInt8,
	convertInt16,
	convertInt32,
	convertInt64,
	convertUint,
	convertUint8,
	convertUint16,
	convertUint32,
	convertUint64,
	convertDuration,
	convertTime,
}

var convertBool = strconv.ParseBool
var convertInt = strconv.Atoi

func convertString(s string) (string, error)   { return s, nil }
func convertFloat64(s string) (float64, error) { return strconv.ParseFloat(s, 64) }

func convertFloat32(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

func convertInt8(s string) (int8, error) {
	parsed, err := strconv.ParseInt(s, 0, 8)
	if err != nil {
		return 0, err
	}
	return int8(parsed), nil
}

func convertInt16(s string) (int16, error) {
	parsed, err := strconv.ParseInt(s, 0, 16)
	if err != nil {
		return 0, err
	}
	return int16(parsed), nil
}

func convertInt32(s string) (int32, error) {
	parsed, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return 0, err
	}
	return int32(parsed), nil
}

func convertInt64(s string) (int64, error) {
	parsed, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func convertUint(s string) (uint, error) {
	parsed, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}

func convertUint8(s string) (uint8, error) {
	parsed, err := strconv.ParseUint(s, 0, 8)
	if err != nil {
		return 0, err
	}
	return uint8(parsed), nil
}

func convertUint16(s string) (uint16, error) {
	parsed, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return 0, err
	}
	return uint16(parsed), nil
}

func convertUint32(s string) (uint32, error) {
	parsed, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		return 0, err
	}
	return uint32(parsed), nil
}

func convertUint64(s string) (uint64, error) {
	parsed, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return 0, err
	}
	return uint64(parsed), nil
}

var convertDuration = time.ParseDuration

func convertTime(s string) (time.Time, error) { return time.Parse(time.RFC3339, s) }
