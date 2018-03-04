package envcfg

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLoader(t *testing.T) {
	type bigConfig struct {
		Untagged int
		B        bool    `env:"B"`
		S        string  `env:"S"`
		SDefault string  `env:"SDEFAULT" default:"AWW YEAH"`
		I        int     `env:"I"`
		IDefault int     `env:"SDEFAULT" default:"7"`
		F32      float32 `env:"F32"`
		F64      float64 `env:"F64"`
		I8       int8    `env:"I8"`
		I16      int16   `env:"I16"`
		I32      int32   `env:"I32"`
		I64      int64   `env:"I64"`
		UI       uint    `env:"UI"`
		UI8      uint8   `env:"UI8"`
		UI16     uint16  `env:"UI16"`
		UI32     uint32  `env:"UI32"`
		UI64     uint64  `env:"UI64"`
	}

	vals := map[string]string{
		"B":    "true",
		"S":    "hi",
		"I":    "-2",
		"F32":  "1.23",
		"F64":  "-3.21",
		"I8":   "8",
		"I16":  "16",
		"I32":  "32",
		"I64":  "64",
		"UI":   "3",
		"UI8":  "4",
		"UI16": "5",
		"UI32": "6",
		"UI64": "7",
	}

	var conf bigConfig
	err := LoadFromMap(vals, &conf)
	assert.Nil(t, err)
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
	}
	assert.Equal(t, expected, conf)
}

// registering converters of the wrong shape
func TestConverterShape(t *testing.T) {
	type foo struct{}
	type bar foo

	tt := []struct {
		desc      string
		converter interface{}
		err       string
	}{
		{
			desc:      "not even a function",
			converter: "I can't even",
			err:       "envcfg: I can't even is not a func",
		},
		{
			desc:      "wrong number of inputs",
			converter: func(s, t string) (foo, error) { return foo{}, nil },
			err:       "envcfg: converter should accept 1 string argument. github.com/btubbs/envcfg.TestConverterShape.func1 accepts 2 arguments",
		},
		{
			desc:      "non-string input",
			converter: func(i int) (foo, error) { return foo{}, nil },
			err:       "envcfg: converter should accept a string argument. github.com/btubbs/envcfg.TestConverterShape.func2 accepts a int argument",
		},
		{
			desc:      "wrong number of outputs",
			converter: func(s string) (foo, string, error) { return foo{}, "", nil },
			err:       "envcfg: converter should return 2 arguments. github.com/btubbs/envcfg.TestConverterShape.func3 returns 3 arguments",
		},
		{
			desc:      "second output not error",
			converter: func(s string) (foo, string) { return foo{}, "" },
			err:       "envcfg: converter's last return value should be error. github.com/btubbs/envcfg.TestConverterShape.func4's last return value is string",
		},
		{
			desc:      "success",
			converter: func(s string) (foo, error) { return foo{}, nil },
			err:       "",
		},
		{
			desc:      "overwriting converter forbidden",
			converter: func(s string) (foo, error) { return foo{}, nil },
			err:       "envcfg: a converter has already been registered for the envcfg.foo type.  cannot also register github.com/btubbs/envcfg.TestConverterShape.func6",
		},
		{
			desc:      "success with type alias",
			converter: func(s string) (bar, error) { return bar{}, nil },
			err:       "",
		},
	}

	for _, tc := range tt {
		err := RegisterConverter(tc.converter)
		if tc.err != "" {
			assert.Equal(t, errors.New(tc.err), err, tc.desc)
		} else {
			assert.Equal(t, nil, err, tc.desc)
		}
	}
}

func TestBuggyConverters(t *testing.T) {
	type foo struct{}
	type myConfig struct {
		B foo `env:"BAR"`
	}
	vals := map[string]string{
		"BAR": "it doesn't really matter",
	}
	tt := []struct {
		desc      string
		converter interface{}
		err       string
	}{
		{
			desc:      "converter that errors",
			converter: func(s string) (foo, error) { return foo{}, errors.New("oops") },
			err:       "envcfg: cannot populate B: oops",
		},
		{
			desc:      "converter that panics",
			converter: func(s string) (foo, error) { panic("I panicked"); return foo{}, nil },
			err:       "envcfg: cannot populate B: github.com/btubbs/envcfg.TestBuggyConverters.func2 panicked: I panicked",
		},
	}

	for _, tc := range tt {
		var conf myConfig
		ec, _ := New()
		ec.RegisterConverter(tc.converter)
		err := ec.LoadFromMap(vals, &conf)
		assert.Equal(t, errors.New(tc.err), err, tc.desc)
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
	assert.Equal(t, "envcfg: no FOO3 value found, and myConfig.F has no default", err.Error())
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
			desc:  "no converter for this type",
			strct: &quux{},
			err:   "envcfg: no converter function found for type envcfg.baz",
		},
	}

	for _, tc := range tt {
		err := LoadFromMap(map[string]string{}, tc.strct)
		assert.Equal(t, tc.err, err.Error(), tc.desc)
	}
}
