[![forthebadge](http://forthebadge.com/images/badges/compatibility-club-penguin.svg)](http://forthebadge.com) [![forthebadge](http://forthebadge.com/images/badges/no-ragrets.svg)](http://forthebadge.com)

# goroadie [![](https://godoc.org/github.com/utrack/goroadie?status.svg)](http://godoc.org/github.com/utrack/goroadie)
Yet another envvar parser for Go.

## Features
This parser can load structs and nested structs of any complexity along with `map[string]string` and basic types like `uint`,`int`,`string` and `bool`.

It is possible to pass initialized maps without breaking them, i.e. existing keys of a map in the struct won't be deleted.

Members' names can be overridden using the `env` tag: `type Foo struct { Bar string ``env:"baz"`` } `

## Installation
Pretty standard, use the `go get` tool:
````
go get github.com/utrack/goroadie
````

## Usage
Just make the struct and pass it! Check [examples](https://godoc.org/github.com/utrack/goroadie#ex-Process--Map) for more.
````
type Database struct {
    Active bool `env:"enabled"`
    Type string
    URI string
    Config map[string]string `env:"opts"`
}

type Config struct {
    Primary Database
    Secondary Database
}

conf := Config{}

err := goroadie.Process("YOURAPP",&conf)
````
Variables scanned in the example:

````
YOURAPP_PRIMARY_ENABLED
YOURAPP_PRIMARY_TYPE
YOURAPP_PRIMARY_URI
YOURAPP_PRIMARY_OPTS_*

YOURAPP_SECONDARY_ENABLED
YOURAPP_SECONDARY_TYPE
YOURAPP_SECONDARY_URI
YOURAPP_SECONDARY_OPTS_*
````
If any config map was initialized and loaded before it won't lose any keys (only replaced): check the [examples](https://godoc.org/github.com/utrack/goroadie#ex-Process--Map) for more.
