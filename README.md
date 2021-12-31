# go-lsp

Toolchain for Language Server Protocol in Go.

- package `protocol` DTO `go`-struct from TypeScript files from official spec. The main goal is to automate the generation of `go`-struct from the official specification.

Official spec for `3.16` version.
https://microsoft.github.io/language-server-protocol/specification

Models `3.16` version
* https://github.com/microsoft/vscode-languageserver-node/blob/release/protocol/3.16.0/types/src/main.ts
* https://github.com/microsoft/vscode-languageserver-node/tree/release/protocol/3.16.0/protocol/src/common

<details>
<summary>How is go-struct generated from TypeScript? Short answer: via json-schema.</summary>
<p><code>go</code>-struct is generated from the json-schema. Json-schema is generated from the TypeScript.</p>
<br/>
<p>Is parsed and extracted the names of all types and interfaces TypeScript files.</p>
<p>The <code>typescript-json-schema</code> https://github.com/YousefED/typescript-json-schema is called to generate the json-schema for each is type or interface.</p>
<p>json-schema is parsed from the https://github.com/santhosh-tekuri/jsonschema library</p>
<p>To generate <code>go</code>-struct via ... WIP</p>
</details>

# For developers

```bash
# go to spec folder
cd spec

# "do all" is meaning
# - install toolchain (required node)
# - download ts files (required curl)
# - generate json-schema files
# - (TODO) generate go-struct files
make all
```
