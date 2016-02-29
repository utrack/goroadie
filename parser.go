// Package goroadie provides the decoder of environment variables
// to Go structs.
//
// It can be used on config loading stage to override
// previously loaded settings (or load the full config altogether).
package goroadie

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/ansel1/merry"
)

// Process populates the struct with data from environment
// variables.
// Accepts this struct's common prefix and the struct itself.
func Process(prefix string, c interface{}) error {
	if reflect.TypeOf(c).Kind() != reflect.Ptr {
		panic(merry.New("Received a struct instead of ptr!"))
	}

	cfg := reflect.ValueOf(c).Elem()
	return structFromEnv(cfg, prefix)
}

// structFromEnv loads the reflect.Value from envvars.
// The prefix is added to the name of every property, so for the struct
//  struct foo {
// 	 A string
//  }
// and prefix ololo the variable OLOLO_A is loaded.
//
// Internal structs are loaded recursively.
func structFromEnv(cfg reflect.Value, prefix string) error {
	// Get the value type
	cTyp := cfg.Type()

	for i := 0; i < cTyp.NumField(); i++ {
		field := cTyp.Field(i)

		// Form this field's envname.
		var ename string
		// Try to read tag 'env' first, or form the name
		// automatically if failed.
		if ename = field.Tag.Get("env"); ename == "" {
			ename = field.Name
		}

		// Append the prefix if supplied
		if prefix != "" {
			ename = prefix + "_" + ename
		}
		ename = strings.ToUpper(ename)

		env := os.Getenv(ename)
		if env == "" && field.Type.Kind() != reflect.Struct && field.Type.Kind() != reflect.Map {
			continue
		}

		var val interface{}
		var err error
		switch field.Type.Kind() {
		case reflect.String:
			val = env
			err = nil
		case reflect.Bool:
			val, err = strconv.ParseBool(env)
		case reflect.Uint:
			val, err = strconv.ParseUint(env, 10, 0)
			val = uint(val.(uint64))
		case reflect.Int:
			val, err = strconv.Atoi(env)
		case reflect.Struct:
			// Recurse if struct
			err := structFromEnv(cfg.FieldByName(field.Name), ename)
			if err != nil {
				return err
			}
			continue
		case reflect.Map:
			// Parse this map
			switch field.Type.Elem().Kind() {
			case reflect.String:
				// Fetch previous values
				vals := cfg.FieldByName(field.Name).Interface().(map[string]string)
				cfg.FieldByName(field.Name).Set(reflect.ValueOf(getMapString(vals, ename)))
			default:
				panic(fmt.Errorf("Can't unmarshal Map with value type %v (%v)", field.Type.Elem().String(), ename))
			}
			continue
		default:
			panic(fmt.Errorf("Unknown config type: %v (field %v)", field.Type.Kind().String(), field.Name))
		}

		if err != nil {
			return merry.New(fmt.Sprintf("Could not parse value %v from envvars!",ename)).WithValue("var", ename)
		}
		cfg.FieldByName(field.Name).Set(reflect.ValueOf(val))

	}
	return nil
}

func getMapString(ret map[string]string, env string) map[string]string {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], env) {
			ret[pair[0][len(env)+1:]] = pair[1]
		}
	}
	return ret
}
