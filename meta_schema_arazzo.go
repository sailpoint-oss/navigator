package navigator

import (
	"encoding/json"
	"sync"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

const metaArazzoRootSchemaURL = "https://navigator.local/meta/arazzo-1.0-root.json"

var (
	metaArazzoCompileOnce sync.Once
	metaArazzoRootSchema  *jsonschema.Schema
	metaArazzoCompileErr  error
)

func compiledArazzoMetaSchema() (*jsonschema.Schema, error) {
	metaArazzoCompileOnce.Do(func() {
		c := jsonschema.NewCompiler()
		c.DefaultDraft(jsonschema.Draft2020)
		c.UseRegexpEngine(compileMetaJSONSchemaPattern)

		raw, err := metaSchemaEmbedded.ReadFile("schemas/meta/arazzo-1.0-root.schema.json")
		if err != nil {
			metaArazzoCompileErr = err
			return
		}

		var doc any
		if err := json.Unmarshal(raw, &doc); err != nil {
			metaArazzoCompileErr = err
			return
		}
		if err := c.AddResource(metaArazzoRootSchemaURL, doc); err != nil {
			metaArazzoCompileErr = err
			return
		}

		metaArazzoRootSchema, metaArazzoCompileErr = c.Compile(metaArazzoRootSchemaURL)
	})
	return metaArazzoRootSchema, metaArazzoCompileErr
}
