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

func TestGetType(t *testing.T) {
	var ts testStruct
	if ty := getType(ts); ty.Kind() != reflect.Struct {
		t.Fatalf("Wanted struct, got %s", ty.Kind())
	}
	if ty := getType(&ts); ty.Kind() != reflect.Struct {
		t.Fatalf("Wanted struct, got %s", ty.Kind())
	}
	ptr := &ts
	if ty := getType(&ptr); ty.Kind() != reflect.Struct {
		t.Fatalf("Wanted struct, got %s", ty.Kind())
	}
}

func TestGetValue(t *testing.T) {
	var ts testStruct
	if v := getValue(ts); v.Kind() != reflect.Struct {
		t.Fatalf("Wanted struct, got %s", v.Kind())
	}
	if v := getValue(&ts); v.Kind() != reflect.Struct {
		t.Fatalf("Wanted struct, got %s", v.Kind())
	}
	ptr := &ts
	if v := getValue(&ptr); v.Kind() != reflect.Struct {
		t.Fatalf("Wanted struct, got %s", v.Kind())
	}
}
