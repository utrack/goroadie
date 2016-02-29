package goyaec_test

import (
	"fmt"
	"os"

	"github.com/utrack/goyaec"
)

// This example shows the basic usage.
//
// It shows how to parse basic structs and nested structs.
func ExampleProcess() {

	type BarType struct {
		C uint `env:"baz"`
	}

	type FooType struct {
		A string
		B int

		Bar BarType

		BazField string
	}

	// Wireup variables
	// You can use prefixes for structs
	os.Setenv("PREF_A", "VarA")
	os.Setenv("PREF_B", "2")

	// Use tag "env" to change the envvar name.
	// Previous prefixes are applied, so for BarType's C it becomes
	// pref (common prefix) -> bar (field name in FooType) -> baz (field name from tag).
	os.Setenv("PREF_BAR_BAZ", "2")

	os.Setenv("PREF_BAZFIELD","BazField")

	conf := FooType{}

	// Remember to pass the pointer!
	err := goyaec.Process("pref", &conf)

	if err != nil {
		panic(err)
	}

	fmt.Println(conf)
	// Output: {VarA 2 {2} BazField}
}

// This example shows the usage with map[string]string included
// in the struct.
func ExampleProcess_map() {
	type T struct {
		A bool `env:"somebool"`
		B map[string]string `env:"dict"`
	}

	// Wire the vars up
	os.Setenv("SOMEBOOL","true")
	// Case of the map keys is saved, capitalize as you want.
	os.Setenv("DICT_Foo","bar")
	os.Setenv("DICT_BaR","baz")
	os.Setenv("DICT_qux","whatever")


	t := T{}
	// Init the map and load some values in it - existing values
	// will be saved!
	t.B = map[string]string{"Foo":"prevFoo","SomeVar":"wow"}

	err := goyaec.Process("",&t)
	if err != nil {
		panic(err)
	}

	// Existing values are saved
	fmt.Printf("Key SomeVar: %v\n",t.B["SomeVar"])

	// ... but replaced if there's an envvar matching it
	fmt.Printf("Key Foo: %v\n",t.B["Foo"])

	// New values are loaded as they should.
	fmt.Printf("Key BaR: %v\n",t.B["BaR"])
	fmt.Printf("Key qux: %v\n",t.B["qux"])
	
	fmt.Printf("Member A: %v\n",t.A)
	// Output: Key SomeVar: wow
	// Key Foo: bar
	// Key BaR: baz
	// Key qux: whatever
	// Member A: true
}
