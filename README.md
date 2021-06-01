# envcfg [![Build Status](https://travis-ci.com/nav-inc/envcfg.svg?branch=master)](https://travis-ci.com/nav-inc/envcfg) 

envcfg is a Go package for loading config from environment variables into struct fields of arbitrary
types.  It's designed with a few guiding principles:

1. If there's a bug in your config values, it's far better to see it right when your application
   starts than at 3 AM Sunday morning when your app finally gets around to using that value.
   So handlers and other business logic components shouldn't read or parse env vars themselves; they
   should access config that's already been read and parsed into concrete Go types for them.
2. It's nice to have one centralized place where it's easy to see all the config that your app (or
   component) needs.
3. Loading config should require as little boilerplate as possible.

envcfg is inspired by the struct tag pattern used by Go when unmarshaling JSON or scanning database
rows.  .

## Simple Example

Imagine you had a struct that held your app's config values, that looked something like this:

```go
  type myAppConfig struct {
    Foo             string       
    Bar             int          
    RefreshInterval time.Duration
  }
```

To load those values from environment variables, you would add struct tags for `env` and optionally
`default`, and then call `envcfg.Load` with a pointer to an instance of that config struct, like
this:

```go
    type myAppConfig struct {
      Foo             string        `env:"FOO" default:"hey there"`
      Bar             int           `env:"BAR"`
      RefreshInterval time.Duration `env:"REFRESH_INTERVAL" default:"2h30m"`
    }

    var conf myAppConfig
    err := envcfg.Load(&conf)
    if err != nil {
      panic(err.Error())
    }

    // now start up your app with your nicely-populated config...
```

In the example above, our config object has a `string`, an `int`, and a `time.Duration`.  It
requires that there be environment variables set for "BAR" and "REFRESH_RATE".  If those aren't set,
then `envcfg.Load` will return an error.  The FOO environment variable may also be set, but the
default of "hey there" will be used if not.

## Built-in Supported Types

As demonstrated in the above example, envcfg already knows how to parse strings into many of the
types built in to Go and its standard library.  Here's the complete list:

    int
    bool              
    string            
    float32           
    float64           
    int8              
    int16             
    int32             
    int64             
    uint              
    uint8             
    uint16            
    uint32            
    uint64            
    time.Duration     
    time.Time         
    *url.URL          
    net.IP
    net.HardwareAddr  
    *mail.Address     
    []*mail.Address   
    *template.Template

## Parsing Other Types

If your struct has a field of some other type, you can tell envcfg how to parse a string into it by
registering your own parser function. You can register parsers for struct types, pointers to struct
types, arrays, and type aliases like `type MyInt int`.  The example below demonstrates registering a
custom parser function that will populate a `*sql.DB` field on a struct.

```go
type myAppConfig struct {
  DB              *sql.DB       `env:"DATABASE_URL"`
}

func main() {
  // to load config we need to instantiate our config struct and pass its pointer to envcfg.Load
  var conf myAppConfig
  err := envcfg.Load(&conf)
  if err != nil {
    fmt.Println(err)
  }
  // now do something useful with the DB connection we just set up...
}

func LoadDBConnection(s string) (*sql.DB, error) {
  db, err := sql.Open("postgres", s)
  if err != nil {
    return nil, err
  }
  return db, nil
}

func init() {
  // A parser func takes one or more strings and returns the type matching your struct field,
  // and an error.
  if err := envcfg.RegisterParser(LoadDBConnection); err != nil {
    panic(err.Error())
  }
}

```

## Loading a Single Field from Multiple Environment Variables

If you have a struct field that should be loaded from multiple environment variables, you can define
a parser function that takes several string arguments.  The `env` tag on your config struct field
must then provide the same number of environment variable names to be passed to the parser
(separated by commas).  Here's an example that loads an Amazon S3 client from three commonly-used
environment variables:

```go
package main

import (
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nav-inc/envcfg"
)

type myAppConfig struct {
	S3 *s3.S3 `env:"AWS_ACCESS_KEY_ID,AWS_SECRET_ACCESS_KEY,AWS_DEFAULT_REGION"`
}

func main() {
	var conf myAppConfig
	err := envcfg.Load(&conf)
	if err != nil {
		panic(err.Error())
	}
  // now the "conf" object has an S3 client on it that you can use to get/post files.
}

// LoadS3Client takes an access key id, secret, and default region, and returns an Amazon S3 client
// instance.
func LoadS3Client(key, secret, region string) (*s3.S3, error) {
	awsCreds := credentials.NewStaticCredentials(key, secret, "")
	_, err := awsCreds.Get()
	if err != nil {
		return nil, fmt.Errorf("bad AWS credentials: %s", err.Error())
	}
	awsCfg := aws.NewConfig().WithCredentials(awsCreds).WithRegion(region).WithS3ForcePathStyle(true)
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return s3.New(sess, awsCfg), nil
}

func init() {
	if err := envcfg.RegisterParser(LoadS3Client); err != nil {
		panic(err.Error())
	}
}

```

## Using a Map Instead of Environment Variables

If you want to provide your own map of values instead of reading environment variables, there's also
a envcfg.LoadFromMap function that will accept your own `map[string]string`.
```go
myVars := map[string]string{
  "FOO": "Let's configure this",
  "BAR": "1",
  "DB":  "postgres://postgres@/my_app?sslmode=disable",
}

...

err := envcfg.LoadFromMap(myVars, &conf)
```

## Instantiating Custom/Multiple Loaders

The examples above all use the default loader provided by the envcfg package.  If you want more
control you can instantiate your own loader (or multiple loaders) using envcfg.New():

```go
    ec, err := envcfg.New()
    if err != nil {
      return err
    }

    err = ec.Load(&conf)
```

If you want a loader without any of the default parsers registered, you can get one by calling
`envcfg.Empty()`:

```go
    ec := envcfg.Empty()
    err := ec.RegisterParser(myParserFunc)
    if err != nil {
      return err
    }

    err = ec.Load(&conf)
```

## Comparison to github.com/kelseyhightower/envconfig
The day after I wrote the first version of this library, a friend pointed out the similar
[envconfig](https://github.com/kelseyhightower/envconfig) library from Kelsey Hightower.  The world
is big enough for both.  There are a few differences that may make you prefer one over the other.

This library (envcfg) lets you register parsers for other people's types without needing to alias
them to add `UnmarshalText` or `Decode` methods.

envcfg supports loading one field from multiple environment variables, while envconfig does not.

envconfig's support for `Set(string) error` methods (like those in `flag.Value`) should allow you to
have struct fields that can be set from env vars or command line flags.  There is no such support in
envcfg.

If you don't designate the environment variable to use in a struct tag, the envconfig library
will use the field's capitalization to guess at the environment variable to use.  envcfg, on the
other hand, will only attempt to load fields with explicit `env` tags, and it requires that either
the environment variable or a `default` tag (or both) be set.

envconfig can load any type that implements the TextUnmarshaler or BinaryUnmarshaler interfaces.
Many stdlib types implement one of these. envcfg should do that too, but hasn't yet.  This would
allow it to delete a bunch of code that registers parser funcs for stdlib types.
