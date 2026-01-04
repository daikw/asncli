package asana

import (
	"fmt"
	"reflect"
	"strconv"
)

// buildQuery converts a params struct to a query map using struct tags.
// It handles string, int, *bool, and map[string]string fields automatically.
//
// Fields must have explicit `query:"key"` tags. Fields without tags are ignored.
// Empty strings, zero ints, and nil pointers are omitted.
//
// Example:
//
//	type Params struct {
//	    Text        string `query:"text"`
//	    AssigneeAny string `query:"assignee.any"`
//	    Limit       int    `query:"limit"`
//	    Completed   *bool  `query:"completed"`
//	    Extra       map[string]string `query:"-"`
//	}
func buildQuery(params any) map[string]string {
	query := make(map[string]string)
	val := reflect.ValueOf(params)

	// Handle pointer to struct
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return query
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeField := typ.Field(i)

		// Get query key from tag
		queryKey := typeField.Tag.Get("query")

		// Special handling for maps with query:"-" - merge them
		if field.Kind() == reflect.Map && queryKey == "-" {
			if !field.IsNil() && field.Type().Key().Kind() == reflect.String && field.Type().Elem().Kind() == reflect.String {
				iter := field.MapRange()
				for iter.Next() {
					k := iter.Key().String()
					v := iter.Value().String()
					if v != "" {
						query[k] = v
					}
				}
			}
			continue
		}

		// Skip fields without query tag or marked with "-"
		if queryKey == "" || queryKey == "-" {
			continue
		}

		// Add to query based on field type
		switch field.Kind() {
		case reflect.String:
			if str := field.String(); str != "" {
				query[queryKey] = str
			}
		case reflect.Int:
			if i := field.Int(); i > 0 {
				query[queryKey] = strconv.FormatInt(i, 10)
			}
		case reflect.Pointer:
			if !field.IsNil() && field.Elem().Kind() == reflect.Bool {
				query[queryKey] = fmt.Sprintf("%t", field.Elem().Bool())
			}
		}
	}

	return query
}
