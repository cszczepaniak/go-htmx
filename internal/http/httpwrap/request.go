package httpwrap

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Request struct {
	Request  *http.Request
	Response http.ResponseWriter
}

// Unmarshal hydrates the given value with values from the request. It does so using struct tags.
// The given target value must be a pointer to a struct type.
//
// There are three main sources for data: form, path, and query. The name used to retrieve the value
// is the name given in the struct tag following the source specifier (so the thing on the right
// side of the colon). Options can be provided after the source:name segment separated by commas.
//   - form gets its value using r.Request.FormValue("name")
//   - path gets its value using r.Request.PathValue("name")
//   - query gets its value using r.Request.URL.Query().Get("name")
//
// Valid options are listed below along with the types with which they can be used.
//   - required: validates that the value is non-empty (i.e. not the zero value of the type). Can be
//     used with string.
//
// Example:
//
//	type myStruct struct {
//		UserID    string `req:"form:userID,required"`
//		UserEmail string `req:"path:userEmail,required"`
//		UserName  string `req:"query:userName"`
//	}
func (r Request) Unmarshal(target any) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Pointer {
		return errors.New("target must be a pointer type")
	}

	if val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct type, but the type was %v", val.Elem().Kind())
	}

	for i := range val.Elem().NumField() {
		field := val.Elem().Type().Field(i)
		tag, ok, err := parseTag(field)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}

		err = r.hydrateField(val.Elem().Field(i), tag)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r Request) hydrateField(field reflect.Value, tag reqTag) error {
	var rawVal string
	switch tag.source {
	case "form":
		rawVal = r.Request.FormValue(tag.name)
	case "path":
		rawVal = r.Request.PathValue(tag.name)
	case "query":
		rawVal = r.Request.URL.Query().Get(tag.name)
	default:
		panic("dev error: should already have a valid source")
	}

	var val any
	switch field.Kind() {
	case reflect.String:
		val = rawVal
	case reflect.Int:
		parsed, err := strconv.ParseInt(rawVal, 10, 64)
		if err != nil {
			return StatusCodeError(err, http.StatusBadRequest)
		}
		val = int(parsed)
	default:
		return fmt.Errorf("unsupported field type for req tag: %v", field.Kind())
	}

	for _, opt := range tag.opts {
		err := handleOpt(tag, opt, val)
		if err != nil {
			return err
		}
	}

	field.Set(reflect.ValueOf(val))
	return nil
}

func handleOpt(tag reqTag, opt string, val any) error {
	switch opt {
	case "required":
		return validateRequired(tag, val)
	default:
		return fmt.Errorf("unknown req field option: %q", opt)
	}
}

func validateRequired(tag reqTag, val any) error {
	switch val := val.(type) {
	case string:
		if val == "" {
			return ErrBadRequestf("field %q is required but was empty", tag.name)
		}
		return nil
	default:
		return StatusCodeError(
			fmt.Errorf("required option not supported for type %T", val),
			http.StatusInternalServerError,
		)
	}
}

type reqTag struct {
	source string
	name   string
	opts   []string
}

func parseTag(field reflect.StructField) (reqTag, bool, error) {
	tagVal := field.Tag.Get("req")
	if tagVal == "" {
		return reqTag{}, false, nil
	}

	// We expect something like source:val,opt1,opt2 (where all options are optional)
	sourceAndName, opts, _ := strings.Cut(tagVal, ",")

	source, name, ok := strings.Cut(sourceAndName, ":")
	if !ok {
		return reqTag{}, false, errors.New("malformed req tag in struct: expected the source and name to be colon-separated")
	}

	switch source {
	case "form", "path", "query":
	default:
		return reqTag{}, false, errors.New("malformed req tag in struct: source must be form, path, or query")
	}

	var optSlice []string
	if opts != "" {
		optSlice = strings.Split(opts, ",")
	}

	return reqTag{
		source: source,
		name:   name,
		opts:   optSlice,
	}, true, nil
}
