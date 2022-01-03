package internal

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	js "github.com/santhosh-tekuri/jsonschema/v5"
)

type SchemaID struct {
	Hash string
	Kind string
}

var NopeSchemaID = SchemaID{Kind: "nope", Hash: "null"}

func MustSchemaID(in *js.Schema) SchemaID {
	sid, err := SchemaIDFrom(in)
	if err != nil {
		panic(fmt.Sprintf("failed calc schema ID for %q: %v", in, err))
	}
	return sid
}

func SchemaIDFrom(in *js.Schema) (SchemaID, error) {
	if len(in.Types) > 1 {
		tran := &js.Schema{Location: in.Location, AnyOf: []*js.Schema{}}
		for idx := range in.Types {
			tran.AnyOf = append(tran.AnyOf, &js.Schema{
				Location: in.Location + "?type=" + in.Types[idx],
				Types:    []string{in.Types[idx]}},
			)
		}
		return SchemaIDFrom(tran)
	}

	if in.Ref != nil {
		return SchemaIDFrom(in.Ref)
	}

	if AssertXXXOf(in) {
		return schemaIDFromXXXOf(in)
	}

	if len(in.Types) == 0 {
		fmt.Printf("without type definition - define as ANY for %q\n", in)
		return NopeSchemaID, nil
	}

	switch in.Types[0] {
	case "boolean", "null":
		return SchemaID{
			Hash: DefName(in.Location),
			Kind: in.Types[0],
		}, nil
	case "number", "string", "integer":
		sid, err := schemIDFromEnum(in)
		if sid != NopeSchemaID {
			return sid, nil
		}
		if err != nil {
			return NopeSchemaID, errors.Wrapf(err, "failed calc schema ID from %q of type %q", in, in.Types[0])
		}

		return SchemaID{
			Hash: DefName(in.Location),
			Kind: in.Types[0],
		}, nil
	case "array":
		sid, err := schemaIDFromArrayID(in)
		if sid != NopeSchemaID {
			return sid, nil
		}
		if err != nil {
			return NopeSchemaID, errors.Wrapf(err, "Failed calc shcema ID from %q of type %q", in, in.Types[0])
		}
		panic(fmt.Sprintf("array type without schema ID"))
	case "object":
		sid, err := schemaIDFromObjectID(in)
		if sid != NopeSchemaID {
			return sid, nil
		}
		if err != nil {
			return NopeSchemaID, errors.Wrapf(err, "Failed calc shcema ID from %q of type %q", in, in.Types[0])
		}
		panic(fmt.Sprintf("object type without schema ID"))
	}

	return NopeSchemaID, fmt.Errorf("not supported type %T", in.Types[0])
}

func schemaIDFromArrayID(in *js.Schema) (SchemaID, error) {
	if len(in.Types) == 0 || in.Types[0] != "array" {
		return NopeSchemaID, nil
	}
	if in.Items == nil {
		return NopeSchemaID, fmt.Errorf("invalid array? no items")
	}
	payload := []string{}
	switch v := in.Items.(type) {
	case *js.Schema:
		return SchemaIDFrom(v)
	case []*js.Schema:
		for idx := range v {
			field := v[idx]
			sid, err := SchemaIDFrom(field)
			if sid != NopeSchemaID {
				payload = append(payload, sid.Hash)
			}
			if err != nil {
				return NopeSchemaID, errors.Wrapf(err, "failed calc schema ID from item %q of array %q ", field, in)
			}

		}
	default:
		return NopeSchemaID, fmt.Errorf("invalid array? not supported type %T", v)
	}

	return SchemaID{
		Hash: hashFromStrSlice(payload),
		Kind: "array",
	}, nil
}

func schemIDFromEnum(in *js.Schema) (SchemaID, error) {
	if len(in.Types) == 0 {
		return NopeSchemaID, nil
	}

	switch in.Types[0] {
	case "number", "string", "integer":
		// pass
	default:
		return NopeSchemaID, nil
	}

	if len(in.Enum) == 0 {
		return NopeSchemaID, nil
	}

	// NOTE: подлагаем что значения всегда в одном и том же порядке и одного и того же типа
	payload, err := json.Marshal(in.Enum)
	if err != nil {
		return NopeSchemaID, errors.Wrap(err, "failed marshal enum")
	}
	hash := hashFromStrSlice([]string{string(payload)})
	return SchemaID{
		Hash: hash,
		Kind: "enum",
	}, nil
}

func schemaIDFromXXXOf(in *js.Schema) (SchemaID, error) {
	if !AssertXXXOf(in) {
		return NopeSchemaID, nil
	}

	hashFn := func(ns string, in []*js.Schema) (SchemaID, error) {
		payload := []string{}
		for idx := range in {
			field := in[idx]
			sid, err := SchemaIDFrom(field)
			if sid != NopeSchemaID {
				payload = append(payload, "opts.xxxof-"+ns+"="+sid.Hash)
			}
			if err != nil {
				return NopeSchemaID, errors.Wrapf(err, "failed calc schema ID for %q", field)
			}
		}
		return SchemaID{
			Hash: hashFromStrSlice(payload),
			Kind: ns,
		}, nil
	}

	if len(in.AllOf) > 0 {
		return hashFn("allof", in.AllOf)
	}
	if len(in.AnyOf) > 0 {
		return hashFn("anyof", in.AnyOf)
	}
	if len(in.OneOf) > 0 {
		return hashFn("oneof", in.OneOf)
	}

	return NopeSchemaID, fmt.Errorf("nothing xxxof %q", in)
}

func schemaIDFromObjectID(in *js.Schema) (SchemaID, error) {
	if in.Ref != nil {
		return SchemaIDFrom(in.Ref)
	}

	if len(in.Types) == 0 || in.Types[0] != "object" {
		return NopeSchemaID, nil
	}

	payload := []string{}
	for key := range in.Properties {
		prop := in.Properties[key]

		if prop.Ref != nil {
			var err error
			sid, err := SchemaIDFrom(prop.Ref)
			if err != nil {
				return NopeSchemaID, errors.Wrapf(err, "failed calc schema ID from %q of props %q of %q", prop.Ref, key, in)
			}
			payload = append(payload, key+"="+sid.Hash)
			continue
		}

		if AssertXXXOf(prop) {
			sid, err := schemaIDFromXXXOf(prop)
			if err != nil {
				return NopeSchemaID, errors.Wrapf(err, "failed calc schema ID from %q (as xxxof) of props %q of %q", prop, key, in)
			}
			payload = append(payload, key+"="+sid.Hash)
			continue
		}

		if len(prop.Types) > 0 {
			payload = append(payload, key+"="+strings.Join(prop.Types, "_"))
			continue
		}

	}

	for idx := range in.Required {
		payload = append(payload, "opts.required="+in.Required[idx])
	}

	return SchemaID{
		Hash: hashFromStrSlice(payload),
		Kind: "object",
	}, nil
}
