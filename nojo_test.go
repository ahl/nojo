// Copyright 2016 Adam H. Leventhal. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/attic-labs/noms/go/marshal"
	"github.com/attic-labs/noms/go/types"
)

type testy struct {
	b bool
}

func check(t *testing.T, v types.Value) {

	buf := new(bytes.Buffer)
	writeValue(v, buf)

	t.Log(buf.String())

	var js map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &js); err != nil {
		t.Fatalf("json Unmarshal() failed %s", err)
	}
	t.Log(js)

	var no map[string]interface{}
	if err := marshal.Unmarshal(v, &no); err != nil {
		t.Fatalf("noms Unmarshal() failed %s", err)
	}
	t.Log(no)

	if !reflect.DeepEqual(js, no) {
		t.Fail()
	}
}

func TestBasic(t *testing.T) {
	check(t, types.NewMap(types.String("b"), types.Bool(true), types.String("n"), types.Number(31337), types.String("s"), types.String("\"yes'")))
}

func TestList(t *testing.T) {
	check(t, types.NewMap(types.String("l"), types.NewList(types.String("a"), types.String("b"), types.String("c"), types.Number(1))))
}

func TestNestedLists(t *testing.T) {
	check(t, types.NewMap(types.String("l"), types.NewList(types.NewList(types.String("bool"), types.Bool(false)),
		types.NewList(types.String("number"), types.Number(7)),
		types.NewList(types.NewList(types.Number(1), types.String("a")), types.String("listy")))))
}

func TestNestedMaps(t *testing.T) {
	check(t, types.NewMap(
		types.String("m1"), types.NewMap(
			types.String("a"), types.Bool(false),
			types.String("b"), types.String("true")),
		types.String("m2"), types.NewMap(
			types.String("a"), types.String("false"),
			types.String("b"), types.Bool(true))))
}
