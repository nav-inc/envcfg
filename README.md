# envcfg [![Build Status](https://travis-ci.org/btubbs/envcfg.svg?branch=master)](https://travis-ci.org/btubbs/envcfg) [![Coverage Status](https://coveralls.io/repos/github/btubbs/envcfg/badge.svg?branch=master)](https://coveralls.io/github/btubbs/envcfg?branch=master)

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
rows.  In the example below, our config object has a `string`, an `int`, a database connection, and
a `time.Duration`.  All of the values in this example are set by envcfg from environment variables
or from defaults set in the struct tags.

## Example

    // In a real app, these would already be set by your environment.
    os.Setenv("BAR", "321")
    os.Setenv("DATABASE_URL", "postgres://postgres@/my_app?sslmode=disable")

    type myAppConfig struct {
      Foo             string        `env:"FOO" default:"hey there"`
      Bar             int           `env:"BAR"`
      DB              *sql.DB       `env:"DATABASE_URL"`
      RefreshInterval time.Duration `env:"REFRESH_INTERVAL" default:"2h30m"`
    }

    // envcfg has built in support for many of Go's built in types, but not *sql.DB, so we'll have to
    // register our own parser.  A parser func takes a string and returns the type matching your
    // struct field, and an error.
    envcfg.RegisterParser(func(s string) (*sql.DB, error) {
      db, err := sql.Open("postgres", s)
      if err != nil {
        return nil, err
      }
      return db, nil
    })

    // to load config we need to instantiate our config struct and pass its pointer to envcfg.Load
    var conf myAppConfig
    err := envcfg.Load(&conf)
    if err != nil {
      fmt.Println(err)
    }
    fmt.Println("Foo", conf.Foo)
    fmt.Println("Bar", conf.Bar)
    fmt.Println("Refresh Interval", conf.RefreshInterval)
    // Output: Foo hey there
    // Bar 321
    // Refresh Interval 2h30m0s

## Built-in Supported Types

As noted in the comments in the above example, envcfg already knows how to parse strings into many of the
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
registering your own parser function (as done for `*sql.DB` in the example above).  You can register
parsers for struct types, pointers to struct types, arrays, and custom types like `type MyInt int`.

## Using a Map Instead of Environment Variables

If you want to provide your own map of values instead of reading environment variables, there's also
a envcfg.LoadFromMap function that will accept your own `map[string]string`.

    myVars := map[string]string{
      "FOO": "Let's configure this",
      "BAR": "1",
      "DB":  "postgres://postgres@/my_app?sslmode=disable",
    }

    ...

    err := envcfg.LoadFromMap(myVars, &conf)

## Instantiating Custom/Multiple Loaders

The examples above all use the default loader provided by the envcfg package.  If you want more
control you can instantiate your own loader (or multiple loaders) using envcfg.New():

    ec, err := envcfg.New()
    if err != nil {
      return err
    }

    err = ec.Load(&conf)

If you want a loader without any of the default parsers registered, you can get one by calling
`envcfg.Empty()`:

    ec := envcfg.Empty()
    err := ec.RegisterParser(myParserFunc)
    if err != nil {
      return err
    }

    err = ec.Load(&conf)

