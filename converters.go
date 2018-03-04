package envcfg

import (
	"strconv"
	"time"
)

var defaultConverters = []interface{}{
	ConvertBool,
	ConvertString,
	ConvertInt,
	ConvertFloat32,
	ConvertFloat64,
	ConvertInt8,
	ConvertInt16,
	ConvertInt32,
	ConvertInt64,
	ConvertUint,
	ConvertUint8,
	ConvertUint16,
	ConvertUint32,
	ConvertUint64,
	ConvertDuration,
	ConvertTime,
}

// The default converters are functions for converting strings into basic built-in Go types.
// These are exported so users can choose to start with a envcfg.Empty() and then pick and choose
// which of these converters they want to register.
// A few types have been omitted: rune (use int32 instead), byte (use uint8 instead), uintptr,
// complex64, and complex128.
var (
	ConvertBool = strconv.ParseBool
	ConvertInt  = strconv.Atoi

	ConvertString  = func(s string) (string, error) { return s, nil }
	ConvertFloat64 = func(s string) (float64, error) { return strconv.ParseFloat(s, 64) }

	ConvertFloat32 = func(s string) (float32, error) {
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return 0, err
		}
		return float32(f), nil
	}

	ConvertInt8 = func(s string) (int8, error) {
		parsed, err := strconv.ParseInt(s, 0, 8)
		if err != nil {
			return 0, err
		}
		return int8(parsed), nil
	}

	ConvertInt16 = func(s string) (int16, error) {
		parsed, err := strconv.ParseInt(s, 0, 16)
		if err != nil {
			return 0, err
		}
		return int16(parsed), nil
	}

	ConvertInt32 = func(s string) (int32, error) {
		parsed, err := strconv.ParseInt(s, 0, 32)
		if err != nil {
			return 0, err
		}
		return int32(parsed), nil
	}

	ConvertInt64 = func(s string) (int64, error) {
		parsed, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return 0, err
		}
		return parsed, nil
	}

	ConvertUint = func(s string) (uint, error) {
		parsed, err := strconv.ParseUint(s, 0, 64)
		if err != nil {
			return 0, err
		}
		return uint(parsed), nil
	}

	ConvertUint8 = func(s string) (uint8, error) {
		parsed, err := strconv.ParseUint(s, 0, 8)
		if err != nil {
			return 0, err
		}
		return uint8(parsed), nil
	}

	ConvertUint16 = func(s string) (uint16, error) {
		parsed, err := strconv.ParseUint(s, 0, 16)
		if err != nil {
			return 0, err
		}
		return uint16(parsed), nil
	}

	ConvertUint32 = func(s string) (uint32, error) {
		parsed, err := strconv.ParseUint(s, 0, 32)
		if err != nil {
			return 0, err
		}
		return uint32(parsed), nil
	}

	ConvertUint64 = func(s string) (uint64, error) {
		parsed, err := strconv.ParseUint(s, 0, 64)
		if err != nil {
			return 0, err
		}
		return uint64(parsed), nil
	}

	ConvertTime     = func(s string) (time.Time, error) { return time.Parse(time.RFC3339, s) }
	ConvertDuration = time.ParseDuration
)
