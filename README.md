# envcfg [![Build Status](https://travis-ci.org/btubbs/envcfg.svg?branch=master)](https://travis-ci.org/btubbs/envcfg)

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
rows.  Example:

    // In a real app, these would already be set by your environment.
    os.Setenv("BAR", "321")
    os.Setenv("DATABASE_URL", "postgres://postgres@/my_app?sslmode=disable")

    type myAppConfig struct {
      Foo             string        `env:"FOO" default:"hey there"`
      Bar             int           `env:"BAR"`
      DB              *sql.DB       `env:"DATABASE_URL"`
      RefreshInterval time.Duration `env:"REFRESH_INTERVAL" default:"2h30m"`
    }

    // envcfg has built in support for Go's built in types, but we need to register our own converter to
    // load other types like *sql.DB.  A converter func takes a string and returns the type matching
    // your struct field, and an error.
    envcfg.RegisterConverter(func(s string) (*sql.DB, error) {
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

Note a couple things about that example:
- envcfg already knows about Go's built-in types like `int`, `string`, etc.  You don't need to tell
  it how to turn "123" into the `int` field on your config struct.
- It also knows how to turn a duration string into a time.Duration.  envcfg doesn't have built-in 
  converters for all standard library types yet, but feel free to make a ticket if there's one you'd
  like to see added.
- If your struct has a type that envcfg doesn't know how to convert to, you can tell it.  In the example
  above we register a function that takes a `DATABASE_URL` and turns it into a connection to a
  Postgres database.  You can register converters for struct types, pointers to struct types,
  arrays, and custom types like `type MyInt int`.

If you want to provide your own map of values instead of reading environment variables, there's also
a envcfg.LoadFromMap function that will accept your own `map[string]string`.

    myVars := map[string]string{
      "FOO": "Let's configure this",
      "BAR": "1",
      "DB":  "postgres://postgres@/my_app?sslmode=disable",
    }

    ...

    err := envcfg.LoadFromMap(myVars, &conf)

The examples above all use the default loader provided by the envcfg package.  If you want more
control you can instantiate your own loader (or multiple loaders) using envcfg.New():

    ec, err := envcfg.New()
    if err != nil {
      return err
    }

    err = ec.Load(&conf)

If you want a loader without any of the default converters registered, you can get one by calling
`envcfg.Empty()`:

    ec := envcfg.Empty()
    err := ec.RegisterConverter(myConverterFunc)
    if err != nil {
      return err
    }

    err = ec.Load(&conf)
