package secretmanager

import (
	"context"
	"fmt"
	"os"
	"reflect"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// Parse looks through a config (c) of type struct for any tag "secretmanager".
// It will use any struct property with this tag to generate a request to the
// GCP secretmanager API.
// It will populate the struct with the values returned by the secretmanager API.
// As a helpful default, it will use context.Background and os.Getenv("GOOGLE_CLOUD_PROJECT")
func Parse(c interface{}) error {
	return ParseWithContextAndProject(context.Background(), os.Getenv("GOOGLE_CLOUD_PROJECT"), c)
}

// ParseWithContext is the same as Parse, except you can pass in a context.
func ParseWithContext(ctx context.Context, c interface{}) error {
	return ParseWithContextAndProject(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"), c)
}

// ParseWithProject is the same as Parse, except you can pass in a project.
func ParseWithProject(project string, c interface{}) error {
	return ParseWithContextAndProject(context.Background(), project, c)
}

// ParseWithContextAndProject is the same as Parse, except you can pass in context and project.
func ParseWithContextAndProject(ctx context.Context, project string, c interface{}) error {
	t := getType(c)
	if err := validate(t); err != nil {
		return err
	}

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	v := getValue(c)
	for i := 0; i < v.NumField(); i++ {
		tag := t.Field(i).Tag.Get("secretmanager")
		if tag == "" {
			continue
		}

		f := v.FieldByName(t.Field(i).Name)
		if err := validateProp(f); err != nil {
			return fmt.Errorf("invalid field: %v", err)
		}

		version := t.Field(i).Tag.Get("version")
		if version == "" {
			version = "latest"
		}

		result, err := client.AccessSecretVersion(context.Background(), &secretmanagerpb.AccessSecretVersionRequest{
			Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", project, tag, version),
		})
		if err != nil {
			return fmt.Errorf("failed to access secret version: %v", err)
		}
		if result.Payload == nil {
			return fmt.Errorf("received invalid response from secretmanager api: %v", result.Payload)
		}

		f.SetString(string(result.Payload.Data))
	}

	return nil
}

func validate(c interface{}) error {
	t := reflect.TypeOf(c)

	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("expected a pointer to a struct, got: %s", t.Kind())
	}

	if t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected kind Struct, got %s", t.Kind())
	}

	return nil
}

func validateProp(f reflect.Value) error {
	if !f.IsValid() {
		return fmt.Errorf("field is invalid: %s", f)
	}
	if !f.CanSet() {
		return fmt.Errorf("cannot set field: %s", f)
	}
	if f.Kind() != reflect.String {
		return fmt.Errorf("secretmanager tags must only be assigned to strings: %s", f)
	}
	return nil
}

func getType(c interface{}) reflect.Type {
	t := reflect.TypeOf(c)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func getValue(c interface{}) reflect.Value {
	v := reflect.ValueOf(c)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}
