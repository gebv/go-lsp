package internal

import (
	"log"
	"strings"

	js "github.com/santhosh-tekuri/jsonschema/v5"
)

func AssertXXXOf(in *js.Schema) bool {
	// asserts
	xxxof := 0
	if len(in.AllOf) > 0 {
		xxxof++
	}
	if len(in.AnyOf) > 0 {
		xxxof++
	}
	if len(in.OneOf) > 0 {
		xxxof++
	}
	if xxxof > 1 {
		log.Printf("multiple xxxof %q\n", in)
		return false
	}

	return xxxof == 1
}

func RootName(in string) string {
	// file:///ClientCapabilities#/definitions/FileOperationClientCapabilities
	// file:///ClientCapabilities#/definitions/CodeActionClientCapabilities/properties/codeActionLiteralSupport
	// file:///ClientCapabilities#/definitions/SemanticTokensClientCapabilities/properties/requests/properties/full/anyOf/0

	in = cleanName(in)
	args := strings.Split(in, "/")
	if len(args) >= 1 {
		return args[0]
	}
	return ""
}

func DefName(in string) string {
	// file:///ClientCapabilities#/definitions/FileOperationClientCapabilities
	// file:///ClientCapabilities#/definitions/CodeActionClientCapabilities/properties/codeActionLiteralSupport
	// file:///ClientCapabilities#/definitions/SemanticTokensClientCapabilities/properties/requests/properties/full/anyOf/0

	in = cleanName(in)
	args := strings.Split(in, "/")
	if len(args) == 3 {
		return args[2]
	}
	return ""
}

func cleanName(in string) string {
	// file:///ClientCapabilities#/definitions/FileOperationClientCapabilities
	// file:///ClientCapabilities#/definitions/CodeActionClientCapabilities/properties/codeActionLiteralSupport
	// file:///ClientCapabilities#/definitions/SemanticTokensClientCapabilities/properties/requests/properties/full/anyOf/0

	in = strings.ReplaceAll(in, "file:///", "")
	in = strings.ReplaceAll(in, "#", "")
	in = strings.ReplaceAll(in, "<T>", "Generic")
	return in
}
