package cmd

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

var m = map[string]int{"one": 1, "two": 2, "three": 3}

var a = make(map[string]interface{})

func main() {
	b := new(bytes.Buffer)

	e := gob.NewEncoder(b)

	// Encoding the map
	err := e.Encode(m)
	if err != nil {
		panic(err)
	}

	var decodedMap map[string]int
	d := gob.NewDecoder(b)

	// Decoding the serialized data
	err = d.Decode(&decodedMap)
	if err != nil {
		panic(err)
	}

	// Ta da! It is a map!
	fmt.Printf("%#v\n", decodedMap)
}
