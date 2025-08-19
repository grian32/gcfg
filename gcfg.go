package gcfg

import (
	"errors"
	"fmt"
	"gcfg/lexer"
	"gcfg/pair"
	"gcfg/parser"
	"reflect"
	"strconv"
	"strings"
)

func Unmarshal(input []byte, v any) error {
	l := lexer.New(input)
	p := parser.New(l)
	parsed, err := p.ParseFile()
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return errors.New("value must be a ptr")
	}

	if rv.IsNil() {
		rv.Set(reflect.New(rv.Type().Elem()))
	}

	elem := rv.Elem()
	if elem.Kind() != reflect.Struct {
		return errors.New("value must be struct")
	}

	return fillStruct(elem, parsed, 0)
}

func fillStruct(elem reflect.Value, parsed map[string]any, recLevel uint32) error {
	t := elem.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		value := elem.Field(i)

		tag := field.Tag.Get("gcfg")
		if tag == "" {
			continue
		}

		switch value.Kind() {
		// so much bs duplicate code when it comes to ints here and in slices, can't rly generalize it by passing the
		// functions or something because it has diff signatures for int and uint64
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, ok := parsed[tag].(string)
			if !ok {
				return fmt.Errorf("field %s: expected string, got %T", field.Name, parsed[tag])
			}
			bits := value.Type().Bits()
			intVal, err := strconv.ParseInt(v, 10, bits)

			if err != nil {
				return err
			}
			value.SetInt(intVal)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, ok := parsed[tag].(string)
			if !ok {
				return fmt.Errorf("field %s: expected string, got %T", field.Name, parsed[tag])
			}
			bits := value.Type().Bits()

			uintVal, err := strconv.ParseUint(v, 10, bits)
			if err != nil {
				return err
			}
			value.SetUint(uintVal)
		case reflect.String:
			v, ok := parsed[tag].(string)
			if !ok {
				return fmt.Errorf("field %s: expected string, got %T", field.Name, parsed[tag])
			}
			value.SetString(v)
		case reflect.Bool:
			v, ok := parsed[tag].(bool)
			if !ok {
				return fmt.Errorf("field %s: expected bool, got %T", field.Name, parsed[tag])
			}
			value.SetBool(v)
		case reflect.Slice:
			arrType := value.Type().Elem().Kind()

			switch arrType {
			case reflect.Struct:
				v, ok := parsed[tag].([]map[string]any)
				if !ok {
					return fmt.Errorf("field %s: wanted map[string]any, got %T", field.Name, parsed[tag])
				}

				elemType := value.Type().Elem()
				arrValue := reflect.MakeSlice(value.Type(), len(v), len(v))

				if recLevel >= 1 {
					return errors.New("nesting past 1 level not allowed")
				}

				for idx := range len(v) {
					structValues := v[idx]
					newElem := reflect.New(elemType).Elem()
					err := fillStruct(newElem, structValues, recLevel+1)
					if err != nil {
						return err
					}
					arrValue.Index(idx).Set(newElem)
				}

				value.Set(arrValue)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v, ok := parsed[tag].([]any)
				if !ok {
					return fmt.Errorf("field %s: wanted []any, got %T", field.Name, parsed[tag])
				}
				bits := value.Type().Elem().Bits()
				arrValue := reflect.MakeSlice(value.Type(), len(v), len(v))

				for idx := range len(v) {
					str, ok := v[idx].(string)
					if !ok {
						return fmt.Errorf("field %s: wanted str as part of []any, got %T", field.Name, v[idx])
					}
					intVal, err := strconv.ParseInt(str, 10, bits)
					if err != nil {
						return err
					}

					arrValue.Index(idx).SetInt(intVal)
				}
				value.Set(arrValue)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				v, ok := parsed[tag].([]any)
				if !ok {
					return fmt.Errorf("field %s: wanted []any, got %T", field.Name, parsed[tag])
				}
				bits := value.Type().Elem().Bits()
				arrValue := reflect.MakeSlice(value.Type(), len(v), len(v))

				for idx := range len(v) {
					str, ok := v[idx].(string)
					if !ok {
						return fmt.Errorf("field %s: wanted str as part of []any, got %T", field.Name, v[idx])
					}
					intVal, err := strconv.ParseUint(str, 10, bits)
					if err != nil {
						return err
					}

					arrValue.Index(idx).SetUint(intVal)
				}
				value.Set(arrValue)
			default:
				v, ok := parsed[tag].([]any)
				if !ok {
					return fmt.Errorf("field %s: wanted []any, got %T", field.Name, parsed[tag])
				}

				elemType := value.Type().Elem()
				arrValue := reflect.MakeSlice(value.Type(), len(v), len(v))

				for idx, item := range v {
					itemVal := reflect.ValueOf(item)
					if !itemVal.Type().ConvertibleTo(elemType) {
						return fmt.Errorf("field %s: wanted %v as part of [], got %T", field.Name, elemType, v[idx])
					}
					arrValue.Index(idx).Set(itemVal.Convert(elemType))
				}
				value.Set(arrValue)
			}

		case reflect.Struct:
			currType := field.Type

			if currType.PkgPath() == "gcfg/pair" && strings.HasPrefix(currType.Name(), "Pair[") {
				p, ok := parsed[tag].(pair.Pair[any, any])
				if !ok {
					return fmt.Errorf("field %s: expected pair.Pair[any, any], got %T", field.Name, p)
				}

				structValues := map[string]any{
					"First":  p.First,
					"Second": p.Second,
				}

				err := fillStruct(value, structValues, recLevel+1)
				if err != nil {
					return err
				}
			} else {
				if recLevel >= 1 {
					return errors.New("nesting past 1 level not allowed")
				}

				structValues, ok := parsed[tag].(map[string]any)
				if !ok {
					return errors.New("bad input for nested struct")
				}

				err := fillStruct(value, structValues, recLevel+1)
				if err != nil {
					return err
				}
			}
		default:
			return errors.New("not accepted value")
		}
	}

	return nil
}
