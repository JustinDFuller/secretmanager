package secretmanager

import (
	"reflect"
	"testing"
)

type testStruct struct {
	ValidString      string
	unexportedString string
	ExportedBytes    []byte
	ExportedInt      int
}

func TestValidate(t *testing.T) {
	var ts testStruct
	if err := validate(&ts); err != nil {
		t.Fatal(err)
	}

	if err := validate(ts); err == nil {
		t.Fatal("Wanted err, got nil")
	}

	e := reflect.ValueOf(&ts).Elem()
	if err := validateProp(e.FieldByName("ValidString")); err != nil {
		t.Fatal(err)
	}

	if err := validateProp(e.FieldByName("unexportedString")); err == nil {
		t.Fatal("Wanted err, got nil")
	}

	if err := validateProp(e.FieldByName("ExportedBytes")); err == nil {
		t.Fatal("Wanted err, got nil")
	}

	if err := validateProp(e.FieldByName("ExportedInt")); err == nil {
		t.Fatal("Wanted err, got nil")
	}
}
