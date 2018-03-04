// Package envcfg provides helpers for reading values from environment variables (or a
// map[string]string), converting them to Go types, and setting their values to fields on a
// user-defined struct.
package envcfg

// TODO: catch panics raised by converters and return them as errors.

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

const (
	cfgTag     = "env"
	defaultTag = "default"
)

// New returns a Loader with the default converters enabled.
func New() (*Loader, error) {
	ec := Empty()
	for _, f := range defaultConverters {
		err := ec.RegisterConverter(f)
		if err != nil {
			return nil, err
		}
	}
	return ec, nil
}

// Empty returns a Loader without any converters enabled.
func Empty() *Loader {
	ec := &Loader{}
	ec.converters = map[reflect.Type]converter{}
	return ec
}

// Our internal converter func takes a string and returns a reflect.Value and an error.  Funcs of
// this type wrap the default converters and user-provided converters that return arbitrary types.
type converter func(string) (reflect.Value, error)

// Loader is a helper for reading values from environment variables (or a map[string]string),
// converting them to Go types, and setting their values to fields on a user-provided struct.
type Loader struct {
	// a map from reflect types to functions that can take a string and return a
	// reflect value of that type.
	converters map[reflect.Type]converter
}

// RegisterConverter takes a func (string) (<anytype>, error) and registers it on the Loader as
// the converter for <anytype>
func (e *Loader) RegisterConverter(f interface{}) error {
	// alright, let's inspect this f and make sure it's a func (string) (sometype, err)
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return fmt.Errorf("envcfg: %v is not a func", f)
	}

	fname := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	// f should accept one argument
	if t.NumIn() != 1 {
		return fmt.Errorf(
			"envcfg: converter should accept 1 string argument. %v accepts %d arguments",
			fname, t.NumIn())
	}
	// it should be a string argument
	if t.In(0) != reflect.TypeOf("") {
		return fmt.Errorf(
			"envcfg: converter should accept a string argument. %s accepts a %v argument",
			fname, t.In(0))
	}
	// it should return two things
	if t.NumOut() != 2 {
		return fmt.Errorf(
			"envcfg: converter should return 2 arguments. %v returns %d arguments",
			fname, t.NumOut())
	}
	// the first can be any type. the second should be error
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if !t.Out(1).Implements(errorInterface) {
		return fmt.Errorf(
			"envcfg: converter's last return value should be error. %s's last return value is %v",
			fname, t.Out(1))
	}
	_, alreadyRegistered := e.converters[t.Out(0)]
	if alreadyRegistered {
		return fmt.Errorf("envcfg: a converter has already been registered for the %v type.  cannot also register %s",
			t.Out(0), fname,
		)
	}

	callable := reflect.ValueOf(f)
	wrapped := func(s string) (v reflect.Value, err error) {
		defer func() {
			p := recover()
			if p != nil {
				// we panicked running the inner converter function.
				err = fmt.Errorf("%s panicked: %s", fname, p)
			}
		}()
		returnvals := callable.Call([]reflect.Value{reflect.ValueOf(s)})
		if !returnvals[1].IsNil() {
			return reflect.Value{}, fmt.Errorf("%v", returnvals[1])
		}
		return returnvals[0], nil
	}
	e.converters[t.Out(0)] = wrapped
	return nil
}

// LoadFromMap loads config from the provided map into the provided struct.
func (e *Loader) LoadFromMap(vals map[string]string, c interface{}) error {
	// assert that c is a struct.
	pointerType := reflect.TypeOf(c)
	if pointerType.Kind() != reflect.Ptr {
		return fmt.Errorf("envcfg: %v is not a pointer", c)
	}

	structType := pointerType.Elem()
	if structType.Kind() != reflect.Struct {
		return fmt.Errorf("envcfg: %v is not a pointer to a struct", c)
	}
	structVal := reflect.ValueOf(c).Elem()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		envKey, ok := field.Tag.Lookup(cfgTag)
		if !ok {
			// this field doesn't have our tag.  Skip.
			continue
		}

		defaultString, defaultOK := field.Tag.Lookup(defaultTag)

		stringVal, ok := vals[envKey]
		if !ok {
			// could not find the string we're looking for in map.  is there a default?
			if defaultOK {
				stringVal = defaultString
			} else {
				return fmt.Errorf("envcfg: no %s value found, and %s.%s has no default", envKey, structType.Name(), field.Name)
			}
		}
		converterFunc, ok := e.converters[field.Type]
		if !ok {
			return fmt.Errorf("envcfg: no converter function found for type %v", field.Type)
		}

		toSet, err := converterFunc(stringVal)
		if err != nil {
			return fmt.Errorf("envcfg: cannot populate %s: %v", field.Name, err)
		}
		structVal.Field(i).Set(toSet)
	}
	return nil
}

// Load loads config from the environment into the provided struct.
func (e *Loader) Load(c interface{}) error {
	// os.Environ guarantees that it will return a list of strings in the form a=b.  It's possible for
	// b to be an empty string.
	env := map[string]string{}
	for _, pair := range os.Environ() {
		parsed := strings.Split(pair, "=")
		env[parsed[0]] = parsed[1]
	}
	return e.LoadFromMap(env, c)
}
