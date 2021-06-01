package envcfg

import (
	"fmt"
	"html/template"
	"net"
	"net/mail"
	"net/url"
	"strconv"
	"time"
)

// The default parsers are functions for converting strings into basic built-in Go types.
// These are exported so users can choose to start with a envcfg.Empty() and then pick and choose
// which of these parsers they want to register.
// A few types have been omitted: rune (use int32 instead), byte (use uint8 instead), uintptr,
// complex64, and complex128.
var DefaultParsers = []interface{}{
	ParseBool,
	ParseString,
	ParseInt,
	ParseFloat32,
	ParseFloat64,
	ParseInt8,
	ParseInt16,
	ParseInt32,
	ParseInt64,
	ParseUint,
	ParseUint8,
	ParseUint16,
	ParseUint32,
	ParseUint64,
	ParseDuration,
	ParseTime,
	ParseURL,
	ParseMAC,
	ParseIP,
	ParseEmailAddress,
	ParseEmailAddressList,
	ParseTemplate,
	ParseBytes,
}

var (
	ParseBool             = strconv.ParseBool
	ParseInt              = strconv.Atoi
	ParseDuration         = time.ParseDuration
	ParseURL              = url.Parse
	ParseMAC              = net.ParseMAC
	ParseEmailAddress     = mail.ParseAddress
	ParseEmailAddressList = mail.ParseAddressList
	ParseTemplate         = template.New("").Parse
)

func ParseIP(s string) (net.IP, error) {
	// weirdly, net.ParseIP returns a nil IP if given invalid input, instead of an error.
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("%s is not a valid IP address", s)
	}
	return ip, nil
}
func ParseTime(s string) (time.Time, error) { return time.Parse(time.RFC3339, s) }

func ParseString(s string) (string, error)   { return s, nil }
func ParseFloat64(s string) (float64, error) { return strconv.ParseFloat(s, 64) }

func ParseFloat32(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

func ParseInt8(s string) (int8, error) {
	parsed, err := strconv.ParseInt(s, 0, 8)
	if err != nil {
		return 0, err
	}
	return int8(parsed), nil
}

func ParseInt16(s string) (int16, error) {
	parsed, err := strconv.ParseInt(s, 0, 16)
	if err != nil {
		return 0, err
	}
	return int16(parsed), nil
}

func ParseInt32(s string) (int32, error) {
	parsed, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return 0, err
	}
	return int32(parsed), nil
}

func ParseInt64(s string) (int64, error) {
	parsed, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func ParseUint(s string) (uint, error) {
	parsed, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}

func ParseUint8(s string) (uint8, error) {
	parsed, err := strconv.ParseUint(s, 0, 8)
	if err != nil {
		return 0, err
	}
	return uint8(parsed), nil
}

func ParseUint16(s string) (uint16, error) {
	parsed, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return 0, err
	}
	return uint16(parsed), nil
}

func ParseUint32(s string) (uint32, error) {
	parsed, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		return 0, err
	}
	return uint32(parsed), nil
}

func ParseUint64(s string) (uint64, error) {
	parsed, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return 0, err
	}
	return uint64(parsed), nil
}

func ParseBytes(s string) ([]byte, error) {
	return []byte(s), nil
}
