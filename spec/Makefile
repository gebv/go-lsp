PROTOCOL_VERSION ?= 3.16.0
JSONRPC_VERSION ?= 6.0.0

SCHEMA_DIR ?= ./schema

PROTOCOL_TYPES ?= utils/is.ts messages.ts connection.ts api.ts protocol.callHierarchy.ts \
	protocol.colorProvider.ts protocol.configuration.ts protocol.declaration.ts \
	protocol.fileOperations.ts protocol.foldingRange.ts protocol.implementation.ts \
	protocol.linkedEditingRange.ts protocol.moniker.ts protocol.progress.ts \
	protocol.selectionRange.ts protocol.semanticTokens.ts protocol.showDocument.ts \
	protocol.ts protocol.typeDefinition.ts protocol.workspaceFolders.ts

BIN_GENJSONSCHEMA ?= typescript-json-schema
BIN_GENJSONSCHEMA_ARGS ?= --aliasRefs --refs --tsNodeRegister --required

all: init download_spec gen_jsonschema

init:
	@echo '{"name":"toolchain","dependencies":{"@types/node":"^17.0.5","vscode-jsonrpc":"${JSONRPC_VERSION}","vscode-languageserver-types":"${PROTOCOL_VERSION}"}}' > package.json
	npm install

download_spec:
	echo "Delete all *.ts files"
	rm -rf ./*.ts
	rm -rf ./utils/*.ts
	mkdir -p ./utils

	echo "Download base.ts for LSP ${PROTOCOL_VERSION} version."
	curl -o base.ts https://raw.githubusercontent.com/microsoft/vscode-languageserver-node/release/protocol/${PROTOCOL_VERSION}/types/src/main.ts

	for tname in ${PROTOCOL_TYPES}; do \
		echo "Download \"$$tname\" for LSP ${PROTOCOL_VERSION} version."; \
		curl -o $$tname https://raw.githubusercontent.com/microsoft/vscode-languageserver-node/release/protocol/${PROTOCOL_VERSION}/protocol/src/common/$$tname; \
	done;

gen_jsonschema:
	mkdir -p ${SCHEMA_DIR}

	echo "Delete all ${SCHEMA_DIR}/*.json files."
	rm -rf ${SCHEMA_DIR}/type.*.json
	rm -rf ${SCHEMA_DIR}/iface.*.json

	echo "Generate json-schema files from *.ts files."
	for fname in `ls *.ts`; do \
		for tname in `sed -n -e "s/^export type \([[:alnum:]_]*\).*$\/\1/p" $$fname`; do \
			echo "Generate type \"$$tname\" from \"$$fname\" to \"${SCHEMA_DIR}/type.$$tname.json\""; \
			$(BIN_GENJSONSCHEMA) ${BIN_GENJSONSCHEMA_ARGS} --id "/$$tname" $$fname $$tname -> "${SCHEMA_DIR}/type.$$tname.json"; \
		done; \
		for tname in `sed -n -e "s/^export interface \([[:alnum:]_]*\).*$\/\1/p" $$fname`; do \
			echo "Generate interface \"$$tname\" from \"$$fname\" to \"${SCHEMA_DIR}/iface.$$tname.json\""; \
			$(BIN_GENJSONSCHEMA) ${BIN_GENJSONSCHEMA_ARGS} --id "/$$tname" $$fname $$tname -> "${SCHEMA_DIR}/iface.$$tname.json"; \
		done; \
	done
