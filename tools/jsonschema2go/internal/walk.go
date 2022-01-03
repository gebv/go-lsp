package internal

import (
	"errors"
	"fmt"
	"strings"

	js "github.com/santhosh-tekuri/jsonschema/v5"
)

type Kind string

var (
	Object  Kind = "object"
	Number  Kind = "number"
	Integer Kind = "integer"
	Boolean Kind = "boolean"
	Str     Kind = "string"
	Null    Kind = "null"
	Array   Kind = "array"
	Oneof   Kind = "oneof"
	Allof   Kind = "allof"
	Anyof   Kind = "anyof"
	Prop    Kind = "prop"
	Ref     Kind = "ref"
	Enum    Kind = "enum"
)

var ErrCircularDep = errors.New("circular dependence")

type WalkFn func(k Kind, parent, current *js.Schema, propKey string) error

var detectCircularDep = func(fn_ WalkFn, paths map[string]int) WalkFn {
	return func(k Kind, parent, current *js.Schema, propKey string) error {
		path := string(k) + parent.Location + current.Location + propKey + strings.Join(parent.Types, "_") + strings.Join(current.Types, "_")

		if paths[path] > 2 {
			return ErrCircularDep
		}
		defer func() {
			paths[path] = paths[path] - 1
		}()
		fmt.Println(path)

		paths[path] = paths[path] + 1

		return fn_(k, parent, current, propKey)
	}
}

func walk(in *js.Schema, fn WalkFn) error {

	panic(`делать +1 когда входим в схему и -1 когда выходим
если в очередной раз когда заходим в схему и там по пути уже равно +1 - это знак что зациклились
при чем отменять движение в глубину если обнаружили цикл
териотически должны завершиться самостоятельно успешно (потому что бесконечные ветки не обходятся)
`)

	var delayedErr error
	delayerErrFn := func(err error) error {
		if err == nil {
			return nil
		}
		if err == ErrCircularDep {

			delayedErr = err
			err = nil
			return nil
		}

		return err
	}

	for idx := range in.Types {
		// JSON Schema basic types:
		// string.
		// number.
		// integer.
		// object.
		// array.
		// boolean.
		// null.
		switch in.Types[idx] {
		case "boolean":
			err := fn(Boolean, in, in, "")
			if err != nil {
				return delayerErrFn(err)
			}
		case "null":
			err := fn(Null, in, in, "")
			if err != nil {
				return delayerErrFn(err)
			}
		case "string":
			err := fn(Str, in, in, "")
			if err != nil {
				return delayerErrFn(err)
			}
		case "number":
			err := fn(Number, in, in, "")
			if err != nil {
				return delayerErrFn(err)
			}
		case "integer":
			err := fn(Integer, in, in, "")
			if err != nil {
				return delayerErrFn(err)
			}
		case "array":
			err := fn(Array, in, in, "")
			if err != nil {
				return delayerErrFn(err)
			}
			err = walkArrayIfNeed(in, fn)
			if err != nil {
				return delayerErrFn(err)
			}
		case "object":
			err := fn(Object, in, in, "")
			if err != nil {
				return delayerErrFn(err)
			}
		}

		if len(in.Enum) > 0 {
			err := fn(Enum, in, in, in.Types[idx])
			if err != nil {
				return delayerErrFn(err)
			}
		}
	}

	for key := range in.Properties {
		err := fn(Prop, in, in.Properties[key], key)
		if err != nil {
			return delayerErrFn(err)
		}
		err = walk(in.Properties[key], fn)
		if err != nil {
			return delayerErrFn(err)
		}
	}

	if in.Ref != nil {
		err := fn(Ref, in, in.Ref, "")
		if err != nil {
			return delayerErrFn(err)
		}
		err = walk(in.Ref, fn)
		if err != nil {
			return delayerErrFn(err)
		}
	}

	if AssertXXXOf(in) {

		for idx := range in.AllOf {
			err := fn(Allof, in, in.AllOf[idx], "")
			if err != nil {
				return delayerErrFn(err)
			}
			err = walk(in.AllOf[idx], fn)
			if err != nil {
				return delayerErrFn(err)
			}
		}
		for idx := range in.AnyOf {
			err := fn(Anyof, in, in.AnyOf[idx], "")
			if err != nil {
				return delayerErrFn(err)
			}
			err = walk(in.AnyOf[idx], fn)
			if err != nil {
				return delayerErrFn(err)
			}
		}
		for idx := range in.OneOf {
			err := fn(Oneof, in, in.OneOf[idx], "")
			if err != nil {
				return delayerErrFn(err)
			}
			err = walk(in.OneOf[idx], fn)
			if err != nil {
				return delayerErrFn(err)
			}
		}
	}

	return delayedErr
}
func Walk(in *js.Schema, fn_ WalkFn) error {
	paths := map[string]int{}
	fn := detectCircularDep(fn_, paths)
	return walk(in, fn)
}

func walkArrayIfNeed(in *js.Schema, fn WalkFn) error {
	var delayedErr error
	delayerErrFn := func(err error) error {
		if err == nil {
			return nil
		}
		if err == ErrCircularDep {

			delayedErr = err
			err = nil
			return nil
		}

		return err
	}

	if in.Items == nil {
		return nil
	}
	switch v := in.Items.(type) {
	case *js.Schema:
		return walk(v, fn)
	case []*js.Schema:
		for idx := range v {
			field := v[idx]
			err := walk(field, fn)

			if err != nil {
				return delayerErrFn(err)
			}
		}
	default:
		return fmt.Errorf("invalid array? not supported type %T", v)
	}

	return delayedErr
}

// func sprintSchemaID(in *js.Schema) string {
// 	sid, err := SchemaIDFrom(in)
// 	return fmt.Sprintf("schema ID = %q, err = %v", sid, err)
// }
