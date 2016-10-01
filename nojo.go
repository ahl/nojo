// Copyright 2016 Adam H. Leventhal. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"encoding/json"

	"github.com/attic-labs/noms/go/spec"
	"github.com/attic-labs/noms/go/types"
)

func writeValue(v types.Value, w io.Writer) {

	switch v.Type().Kind() {
	case types.RefKind:
		// Convert the hash to a string and treat it like one
		v = types.String(v.(types.Ref).Hash().String())
		fallthrough
	case types.BoolKind, types.StringKind:
		b, _ := json.Marshal(v)
		w.Write(b)
	case types.NumberKind:
		// While this could also be handled in the previous case, I hate the exponential notation it uses.
		io.WriteString(w, strconv.FormatFloat(float64(v.(types.Number)), 'f', -1, 64))
	case types.StructKind:
		writeStruct(v.(types.Struct), w)
	case types.ListKind:
		writeList(v.(types.List), w)
	case types.SetKind:
		writeSet(v.(types.Set), w)
	case types.MapKind:
		writeMap(v.(types.Map), w)
	case types.BlobKind:
		writeBlob(v.(types.Blob), w)
	default:
		panic(fmt.Sprintf("unexpected type: %s", types.KindToString[v.Type().Kind()]))
	}
}

func writeStruct(v types.Struct, w io.Writer) {
	io.WriteString(w, "{")
	var first = true
	v.Type().Desc.(types.StructDesc).IterFields(func(name string, _ *types.Type) {
		if !first {
			io.WriteString(w, ",")
		}
		first = false

		io.WriteString(w, "\""+name+"\":")
		writeValue(v.Get(name), w)
	})
	io.WriteString(w, "}")
}

func writeIter(iter func(func(types.Value)), w io.Writer) {
	io.WriteString(w, "[")
	var first = true
	iter(func(v types.Value) {
		if !first {
			io.WriteString(w, ",")
		}
		first = false

		writeValue(v, w)
	})
	io.WriteString(w, "]")
}

func writeList(v types.List, w io.Writer) {
	writeIter(func(cb func(types.Value)) {
		v.IterAll(func(value types.Value, _ uint64) {
			cb(value)
		})
	}, w)
}

func writeSet(v types.Set, w io.Writer) {
	writeIter(func(cb func(types.Value)) {
		v.IterAll(func(value types.Value) {
			cb(value)
		})
	}, w)
}

func writeBlob(v types.Blob, w io.Writer) {
	writeIter(func(cb func(types.Value)) {
		br := bufio.NewReader(v.Reader())
		for i := uint64(0); i < v.Len(); i++ {
			b, _ := br.ReadByte()
			cb(types.Number(b))
		}
	}, w)
}

func writeMap(v types.Map, w io.Writer) {
	if v.Len() == 0 {
		io.WriteString(w, "{}")
	} else {
		keyType := v.Type().Desc.(types.CompoundDesc).ElemTypes[0]

		if keyType != types.StringType {
			panic(fmt.Sprintf("json can only represent maps whose keys are of type String (found %s)", types.KindToString[keyType.Kind()]))
		}

		io.WriteString(w, "{")
		var first = true
		v.IterAll(func(key, value types.Value) {
			if !first {
				io.WriteString(w, ",")
			}
			first = false
			writeValue(key, w)
			io.WriteString(w, ":")
			writeValue(value, w)
		})
		io.WriteString(w, "}")
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no path specified")
		os.Exit(1)
	}
	database, value, err := spec.GetPath(os.Args[1])
	if err != nil {
		fmt.Printf("path '%s' is invalid\n", os.Args[1])
		os.Exit(1)
	}

	defer database.Close()
	writeValue(value, os.Stdout)
}
