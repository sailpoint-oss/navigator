package navigator

//go:generate bun run schemas:build

import (
	"embed"
	"encoding/json"
	"sync"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed schemas/meta/*.json
var metaSchemaEmbedded embed.FS

const (
	metaRootSchemaURL     = "https://navigator.local/meta/openapi-root.json"
	metaFragmentSchemaURL = "https://navigator.local/meta/openapi-fragment.json"
)

var (
	metaCompileOnce sync.Once
	metaRootSchema  *jsonschema.Schema
	metaFragSchema  *jsonschema.Schema
	metaCompileErr  error
)

func compiledMetaSchemas() (*jsonschema.Schema, *jsonschema.Schema, error) {
	metaCompileOnce.Do(func() {
		c := jsonschema.NewCompiler()
		c.DefaultDraft(jsonschema.Draft2020)
		c.UseRegexpEngine(compileMetaJSONSchemaPattern)

		files := []struct {
			file string
			url  string
		}{
			{"schemas/meta/openapi-root.schema.json", metaRootSchemaURL},
			{"schemas/meta/openapi-fragment.schema.json", metaFragmentSchemaURL},
		}
		for _, f := range files {
			raw, err := metaSchemaEmbedded.ReadFile(f.file)
			if err != nil {
				metaCompileErr = err
				return
			}
			var doc any
			if err := json.Unmarshal(raw, &doc); err != nil {
				metaCompileErr = err
				return
			}
			if err := c.AddResource(f.url, doc); err != nil {
				metaCompileErr = err
				return
			}
		}

		metaRootSchema, metaCompileErr = c.Compile(metaRootSchemaURL)
		if metaCompileErr != nil {
			return
		}
		metaFragSchema, metaCompileErr = c.Compile(metaFragmentSchemaURL)
	})
	return metaRootSchema, metaFragSchema, metaCompileErr
}
