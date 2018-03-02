package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
)

func main() {

	// register some basic converters
	// strings
	registerConverter(func(s string) (string, error) { return s, nil })
	// integers
	registerConverter(strconv.Atoi)

	// float32

	// float64

	vals := map[string]string{
		"NUMBER": "3",
		"BLAH":   "asdfasdfasfd",
	}
	conf := myConfig{}
	fromMap(vals, conf)
	err := fromMap(vals, &conf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(conf)
}

type myConfig struct {
	//DB              *sql.DB       `config:"DATABASE_URL"`
	//RefreshInterval time.Duration `cfg:"REFRESH_INTERVAL_SECS" default:"30"`
	SomeString    string `cfg:"BLAH"`
	Number        int    `cfg:"NUMBER"`
	DefaultNumber int    `cfg:"DEFAULT_NUMBER" default:"2"`
	MissingNumber int    `cfg:"MISSING_NUMBER"`
	MissingTag    int
}

// a map from reflect types to functions that can take a string and return a
// reflect value of that type.
type converter func(string) reflect.Value

var converters = map[reflect.Type]converter{}

func registerConverter(f interface{}) error {
	// alright, let's inspect this f and make sure it's a func (string) (sometype, err)
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return fmt.Errorf("%v is not a func", f)
	}

	fname := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	// f should accept one argument
	if t.NumIn() != 1 {
		return fmt.Errorf(
			"converter should accept 1 string argument. %v accepts %d arguments",
			fname, t.NumIn())
	}
	// it should be a string argument
	if t.In(0) != reflect.TypeOf("") {
		return fmt.Errorf(
			"converter should accept a string argument. %s accepts a %v argument",
			fname, t.In(0))
	}
	// it should return two things
	if t.NumOut() != 2 {
		return fmt.Errorf(
			"converter should return 2 arguments. %v returns %d arguments",
			fname, t.NumOut())
	}
	// the first can be any type. the second should be error
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if !t.Out(1).Implements(errorInterface) {
		return fmt.Errorf(
			"converter's last return value should be error. %s's last return value is %v",
			fname, t.Out(1))
	}
	callable := reflect.ValueOf(f)
	wrapped := func(s string) reflect.Value {
		returnvals := callable.Call([]reflect.Value{reflect.ValueOf(s)})
		// is err is nil, return the real result
		if returnvals[1].IsNil() {
			return returnvals[0]
		}
		// but if we have an err, return the zero value for the type
		return reflect.Zero(t.Out(0))
	}
	converters[t.Out(0)] = wrapped
	return nil
}

// ok now I want a function that will take a map[string]string
func fromMap(vals map[string]string, c interface{}) error {
	// assert that c is a struct.
	pointerType := reflect.TypeOf(c)
	if pointerType.Kind() != reflect.Ptr {
		return fmt.Errorf("%v is not a pointer", c)
	}

	structType := pointerType.Elem()
	if structType.Kind() != reflect.Struct {
		return fmt.Errorf("%v is not a pointer to a struct", c)
	}
	structVal := reflect.ValueOf(c).Elem()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		envKey := field.Tag.Get("cfg")
		defaultString := field.Tag.Get("default")

		toSet, ok := vals[envKey]
		if !ok {
			toSet = defaultString
		}
		converterFunc, ok := converters[field.Type]
		if !ok {
			return fmt.Errorf("no converter function found for type %v", field.Type)
		}
		structVal.Field(i).Set(converterFunc(toSet))
	}
	return nil
}
