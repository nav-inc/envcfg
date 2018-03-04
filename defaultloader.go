package envcfg

import "fmt"

// this file ensures that a default loader is created and available on the package, so users with
// simple cases can just do envcfg.Load.

var defaultLoader *Loader

func init() {
	// we can only fail here if one of the hardcoded default converters has the wrong function
	// signature.  If that does fail, fail hard.
	var err error
	defaultLoader, err = New()
	if err != nil {
		panic(fmt.Sprintf("could not init default loader: %v", err))
	}
}

// Load loads config from the environment into the provided struct.
func Load(c interface{}) error {
	return defaultLoader.Load(c)
}

// LoadFromMap loads config from the provided map into the provided struct.
func LoadFromMap(vals map[string]string, c interface{}) error {
	return defaultLoader.LoadFromMap(vals, c)
}

// RegisterConverter takes a func (string) (<anytype>, error) and registers it on the default loader
// as the converter for <anytype>.
func RegisterConverter(f interface{}) error {
	return defaultLoader.RegisterConverter(f)
}
