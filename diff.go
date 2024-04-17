package diff

import (
	"encoding/json"
	"errors"
	"reflect"
)

func Unmarshal(bytes []byte, v any) (map[string]any, error) {
	var diff map[string]any
	err := json.Unmarshal(bytes, &diff)
	if err != nil {
		return nil, err
	}

	structValue := reflect.ValueOf(v).Elem()
	if structValue.Kind() != reflect.Struct {
		return nil, errors.New("expected a struct")
	}

	err = populateStruct(structValue, diff)
	if err != nil {
		return nil, err
	}

	return diff, nil
}

func populateStruct(structValue reflect.Value, diff map[string]any) error {
	typeOfT := structValue.Type()
	for i := 0; i < typeOfT.NumField(); i++ {
		field := typeOfT.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = field.Name // Use field name if no json tag
		}

		value, ok := diff[jsonTag]
		if !ok {
			continue // Guard: skip if field is not present in the diff
		}

		if field.Type.Kind() == reflect.Slice {
			sliceValue, ok := value.([]any)
			if !ok {
				continue // Guard: skip if diff entry is not a slice as expected
			}
			err := handleSliceField(structValue.Field(i), sliceValue)
			if err != nil {
				return err
			}
			delete(diff, jsonTag)
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			structMap, ok := value.(map[string]any)
			if !ok {
				continue // Guard: skip if diff entry is not a map as expected
			}
			err := populateStruct(structValue.Field(i), structMap)
			if err != nil {
				return err
			}
			delete(diff, jsonTag)
			continue
		}

		val := reflect.ValueOf(value)
		if !val.Type().ConvertibleTo(field.Type) {
			continue // Guard: skip if the value is not convertible to the field type
		}

		fieldValue := structValue.Field(i)
		if fieldValue.CanSet() {
			fieldValue.Set(val.Convert(field.Type))
		}
		delete(diff, jsonTag)
	}

	return nil
}

func handleSliceField(sliceValue reflect.Value, sliceData []any) error {
	itemType := sliceValue.Type().Elem()
	for _, item := range sliceData {
		newItem := reflect.New(itemType).Elem()
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue // Guard: skip non-map elements
		}
		err := populateStruct(newItem, itemMap)
		if err != nil {
			return err
		}
		sliceValue.Set(reflect.Append(sliceValue, newItem))
	}
	return nil
}
