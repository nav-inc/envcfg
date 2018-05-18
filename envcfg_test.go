package envcfg

import (
	"errors"
	"html/template"
	"net"
	"net/mail"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLoader(t *testing.T) {
	type bigConfig struct {
		Untagged       int
		B              bool               `env:"B"`
		S              string             `env:"S"`
		SDefault       string             `env:"SDEFAULT" default:"AWW YEAH"`
		I              int                `env:"I"`
		IDefault       int                `env:"SDEFAULT" default:"7"`
		F32            float32            `env:"F32"`
		F64            float64            `env:"F64"`
		I8             int8               `env:"I8"`
		I16            int16              `env:"I16"`
		I32            int32              `env:"I32"`
		I64            int64              `env:"I64"`
		UI             uint               `env:"UI"`
		UI8            uint8              `env:"UI8"`
		UI16           uint16             `env:"UI16"`
		UI32           uint32             `env:"UI32"`
		UI64           uint64             `env:"UI64"`
		Dur            time.Duration      `env:"DUR"`
		When           time.Time          `env:"TIME"`
		URL            *url.URL           `env:"PUBLIC_URL"`
		MAC            net.HardwareAddr   `env:"MAC_ADDRESS"`
		Email          *mail.Address      `env:"EMAIL_ADDRESS"`
		EmailAddresses []*mail.Address    `env:"EMAIL_ADDRESSES"`
		Template       *template.Template `env:"GREETING_TEMPLATE"`
		IP             net.IP             `env:"IP"`
	}

	vals := map[string]string{
		"B":                 "true",
		"S":                 "hi",
		"I":                 "-2",
		"F32":               "1.23",
		"F64":               "-3.21",
		"I8":                "8",
		"I16":               "16",
		"I32":               "32",
		"I64":               "64",
		"UI":                "3",
		"UI8":               "4",
		"UI16":              "5",
		"UI32":              "6",
		"UI64":              "7",
		"DUR":               "2h30m",
		"TIME":              "2017-12-25T00:00:00Z",
		"PUBLIC_URL":        "https://www.example.com/",
		"MAC_ADDRESS":       "01:23:45:67:89:ab:cd:ef:00:00:01:23:45:67:89:ab:cd:ef:00:00",
		"EMAIL_ADDRESS":     "Brent Tubbs <brent.tubbs@gmail.com>",
		"EMAIL_ADDRESSES":   "Alice <alice@example.com>, Bob <bob@example.com>, Eve <eve@example.com>",
		"GREETING_TEMPLATE": "Hello {{.Name}}!",
		"IP":                "8.8.8.8",
	}

	var conf bigConfig
	err := LoadFromMap(vals, &conf)
	assert.Nil(t, err)
	dur, _ := time.ParseDuration("2h30m")
	mac, _ := net.ParseMAC("01:23:45:67:89:ab:cd:ef:00:00:01:23:45:67:89:ab:cd:ef:00:00")
	tmpl, _ := template.New("").Parse("Hello {{.Name}}!")
	expected := bigConfig{
		Untagged: 0,
		B:        true,
		S:        "hi",
		SDefault: "AWW YEAH",
		I:        -2,
		IDefault: 7,
		F32:      1.23,
		F64:      -3.21,
		I8:       8,
		I16:      16,
		I32:      32,
		I64:      64,
		UI:       3,
		UI8:      4,
		UI16:     5,
		UI32:     6,
		UI64:     7,
		Dur:      dur,
		When:     time.Date(2017, time.December, 25, 0, 0, 0, 0, time.UTC),
		URL:      &url.URL{Scheme: "https", Host: "www.example.com", Path: "/"},
		MAC:      mac,
		Email:    &mail.Address{Name: "Brent Tubbs", Address: "brent.tubbs@gmail.com"},
		EmailAddresses: []*mail.Address{
			&mail.Address{Name: "Alice", Address: "alice@example.com"},
			&mail.Address{Name: "Bob", Address: "bob@example.com"},
			&mail.Address{Name: "Eve", Address: "eve@example.com"},
		},
		Template: tmpl,
		IP:       net.IPv4(8, 8, 8, 8),
	}
	assert.Equal(t, expected, conf)
}

func TestParserShape(t *testing.T) {
	type foo struct{}
	type bar foo

	tt := []struct {
		desc   string
		parser interface{}
		err    string
	}{
		{
			desc:   "not even a function",
			parser: "I can't even",
			err:    "envcfg: I can't even is not a func",
		},
		{
			desc:   "wrong number of inputs",
			parser: func(s, t string) (foo, error) { return foo{}, nil },
			err:    "envcfg: parser should accept 1 string argument. github.com/btubbs/envcfg.TestParserShape.func1 accepts 2 arguments",
		},
		{
			desc:   "non-string input",
			parser: func(i int) (foo, error) { return foo{}, nil },
			err:    "envcfg: parser should accept a string argument. github.com/btubbs/envcfg.TestParserShape.func2 accepts a int argument",
		},
		{
			desc:   "wrong number of outputs",
			parser: func(s string) (foo, string, error) { return foo{}, "", nil },
			err:    "envcfg: parser should return 2 arguments. github.com/btubbs/envcfg.TestParserShape.func3 returns 3 arguments",
		},
		{
			desc:   "second output not error",
			parser: func(s string) (foo, string) { return foo{}, "" },
			err:    "envcfg: parser's last return value should be error. github.com/btubbs/envcfg.TestParserShape.func4's last return value is string",
		},
		{
			desc:   "success",
			parser: func(s string) (foo, error) { return foo{}, nil },
			err:    "",
		},
		{
			desc:   "overwriting parser forbidden",
			parser: func(s string) (foo, error) { return foo{}, nil },
			err:    "envcfg: a parser has already been registered for the envcfg.foo type.  cannot also register github.com/btubbs/envcfg.TestParserShape.func6",
		},
		{
			desc:   "success with type alias",
			parser: func(s string) (bar, error) { return bar{}, nil },
			err:    "",
		},
	}

	for _, tc := range tt {
		err := RegisterParser(tc.parser)
		if tc.err != "" {
			assert.Equal(t, errors.New(tc.err), err, tc.desc)
		} else {
			assert.Equal(t, nil, err, tc.desc)
		}
	}
}

func TestBuggyParsers(t *testing.T) {
	type foo struct{}
	type myConfig struct {
		B foo `env:"BAR"`
	}
	vals := map[string]string{
		"BAR": "it doesn't really matter",
	}
	tt := []struct {
		desc   string
		parser interface{}
		err    string
	}{
		{
			desc:   "parser that errors",
			parser: func(s string) (foo, error) { return foo{}, errors.New("oops") },
			err:    "1 error occurred:\n\n* envcfg: cannot populate B: oops",
		},
		{
			desc:   "parser that panics",
			parser: func(s string) (foo, error) { panic("I panicked"); return foo{}, nil },
			err:    "1 error occurred:\n\n* envcfg: cannot populate B: github.com/btubbs/envcfg.TestBuggyParsers.func2 panicked: I panicked",
		},
	}

	for _, tc := range tt {
		var conf myConfig
		ec, _ := New()
		ec.RegisterParser(tc.parser)
		err := ec.LoadFromMap(vals, &conf)
		assert.Equal(t, tc.err, err.Error(), tc.desc)
	}
}

func TestLoadFromEnv(t *testing.T) {
	type myConfig struct {
		F string `env:"FOO2"`
	}

	os.Setenv("FOO2", "This is foo")

	var conf myConfig
	err := Load(&conf)
	assert.Nil(t, err)
	assert.Equal(t, "This is foo", conf.F)

}

func TestMissingValue(t *testing.T) {
	type myConfig struct {
		F string `env:"FOO3"`
	}

	var conf myConfig
	err := LoadFromMap(map[string]string{}, &conf)
	assert.Equal(t, "1 error occurred:\n\n* no FOO3 value found, and myConfig.F has no default", err.Error())
}

func TestBadStructs(t *testing.T) {
	type baz struct{}
	type quux struct {
		B baz `env:"FOO" default:"BAR"`
	}

	tt := []struct {
		desc  string
		strct interface{}
		err   string
	}{
		{
			desc:  "not a pointer",
			strct: "I'm a string, dummy",
			err:   "envcfg: I'm a string, dummy is not a pointer",
		},
		{
			desc:  "a pointer, but not to a struct",
			strct: &[]string{},
			err:   "envcfg: &[] is not a pointer to a struct",
		},
		{
			desc:  "no parser for this type",
			strct: &quux{},
			err:   "1 error occurred:\n\n* no parser function found for type envcfg.baz",
		},
	}

	for _, tc := range tt {
		err := LoadFromMap(map[string]string{}, tc.strct)
		assert.Equal(t, tc.err, err.Error(), tc.desc)
	}
}

func TestEnvListToMap(t *testing.T) {
	ss := []string{
		"FOO==2",
		"BAR=2=",
		"BAZ==",
		"QUUX=",
	}

	expected := map[string]string{
		"FOO":  "=2",
		"BAR":  "2=",
		"BAZ":  "=",
		"QUUX": "",
	}
	assert.Equal(t, expected, envListToMap(ss))
}
