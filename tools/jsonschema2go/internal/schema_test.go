package internal

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	js "github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/assert"
)

func Test_SchemaID(t *testing.T) {
	dat := openFile("iface.DocumentSymbol.json")
	s := mustJsonSchema(t, dat)
	sid, err := SchemaIDFrom(s)
	assert.NoError(t, err)
	t.Log(sid)
	walkFn := func(k Kind, prev, in *js.Schema, propName string) error {
		if k == Ref {
			fmt.Printf("%s %s\n", k, in)
		}
		return nil
	}
	Walk(s, walkFn)
}

func mustJsonSchema(t *testing.T, in string) *js.Schema {
	t.Helper()
	c := js.NewCompiler()
	err := c.AddResource("root", strings.NewReader(in))
	assert.NoError(t, err, "add to storage manage")
	schema, err := c.Compile("root")
	assert.NoError(t, err, "compile json-schema")
	return schema
}

func openFile(in string) string {
	dat, _ := ioutil.ReadFile(filepath.Join("testdata", in))
	return string(dat)
}
