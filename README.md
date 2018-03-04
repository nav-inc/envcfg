# envcfg [![Build Status](https://travis-ci.org/btubbs/envcfg.svg?branch=master)](https://travis-ci.org/btubbs/envcfg)

envcfg is a Go package for loading config from environment variables into struct fields of arbitrary
types.  It's inspired by the struct tag pattern used by the Go stdlib when unmarshaling JSON or
scanning database rows.  Example: 

    // In a real app, these would already be set by your environment.
    os.Setenv("BAR", "321")
    os.Setenv("DATABASE_URL", "postgres://postgres@/my_app?sslmode=disable")

    type myAppConfig struct {
      Foo string  `env:"FOO" default:"hey there"`
      Bar int     `env:"BAR"`
      DB  *sql.DB `env:"DATABASE_URL"`
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
    // Output: Foo hey there
    // Bar 321

