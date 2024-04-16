package diff

import (
	"encoding/json"
	"errors"
	"reflect"
)

func UnmarshalWithDiff[T any](bytes []byte) (*T, map[string]any, error) {
	// unmarshal to a map that we then use to create a diff
	var mapped map[string]any
	err := json.Unmarshal(bytes, &mapped)
	if err != nil {
		return nil, nil, err
	}

	// create a new pointer to the struct
	structured := new(T)
	// get the value of the struct
	structValue := reflect.ValueOf(structured).Elem()
	// check if the value is a struct
	if structValue.Kind() != reflect.Struct {
		return nil, nil, errors.New("expected a struct")
	}

	// iterate over fields
	typeOfT := structValue.Type()
	for i := 0; i < typeOfT.NumField(); i++ {
		// get the field and the json tag
		field := typeOfT.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			// default to field name if no json tag is present
			jsonTag = field.Name
		}

		// check if the field is present in the map
		if value, ok := mapped[jsonTag]; ok {
			// check if the value is convertible to the field type
			val := reflect.ValueOf(value)
			if val.Type().ConvertibleTo(field.Type) {
				// set the value of the field
				fieldValue := structValue.Field(i)
				if fieldValue.CanSet() {
					fieldValue.Set(val.Convert(field.Type))
				}
			}

			// remove the field from the map if it was set
			delete(mapped, jsonTag)
		}
	}

	// return the structured value, the diff and the error
	return structured, mapped, err
}

