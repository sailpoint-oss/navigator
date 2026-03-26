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
	metaCompileOnce          sync.Once
	metaRootSchema           *jsonschema.Schema
	metaFragSchema           *jsonschema.Schema
	metaRootSchemasByVersion map[Version]*jsonschema.Schema
	metaFragSchemasByVersion map[Version]map[FragmentType]*jsonschema.Schema
	metaCompileErr           error
)

func compiledMetaSchemas() (*jsonschema.Schema, *jsonschema.Schema, error) {
	compileMetaSchemas()
	return metaRootSchema, metaFragSchema, metaCompileErr
}

func compiledVersionedMetaSchemas() (map[Version]*jsonschema.Schema, map[Version]map[FragmentType]*jsonschema.Schema, error) {
	compileMetaSchemas()
	return metaRootSchemasByVersion, metaFragSchemasByVersion, metaCompileErr
}

func compileMetaSchemas() {
	metaCompileOnce.Do(func() {
		c := jsonschema.NewCompiler()
		c.DefaultDraft(jsonschema.Draft2020)
		c.UseRegexpEngine(compileMetaJSONSchemaPattern)

		metaRootSchemasByVersion = make(map[Version]*jsonschema.Schema, 4)
		metaFragSchemasByVersion = make(map[Version]map[FragmentType]*jsonschema.Schema, 4)

		for _, res := range metaSchemaResources() {
			raw, err := metaSchemaEmbedded.ReadFile(res.file)
			if err != nil {
				metaCompileErr = err
				return
			}
			var doc any
			if err := json.Unmarshal(raw, &doc); err != nil {
				metaCompileErr = err
				return
			}
			if err := c.AddResource(res.url, doc); err != nil {
				metaCompileErr = err
				return
			}
		}

		metaRootSchema, metaCompileErr = c.Compile(metaRootSchemaURL)
		if metaCompileErr != nil {
			return
		}
		metaFragSchema, metaCompileErr = c.Compile(metaFragmentSchemaURL)
		if metaCompileErr != nil {
			return
		}

		for _, ver := range []Version{Version20, Version30, Version31, Version32} {
			rootURL := metaRootSchemaURLForVersion(ver)
			metaRootSchemasByVersion[ver], metaCompileErr = c.Compile(rootURL)
			if metaCompileErr != nil {
				return
			}
			metaFragSchemasByVersion[ver] = make(map[FragmentType]*jsonschema.Schema)
			for _, frag := range fragmentTypesForVersion(ver) {
				fragURL, ok := metaFragmentSchemaURLForVersion(ver, frag)
				if !ok {
					continue
				}
				metaFragSchemasByVersion[ver][frag], metaCompileErr = c.Compile(fragURL)
				if metaCompileErr != nil {
					return
				}
			}
		}
	})
}

type metaSchemaResource struct {
	file string
	url  string
}

func metaSchemaResources() []metaSchemaResource {
	resources := []metaSchemaResource{
		{file: "schemas/meta/openapi-root.schema.json", url: metaRootSchemaURL},
		{file: "schemas/meta/openapi-fragment.schema.json", url: metaFragmentSchemaURL},
	}
	for _, ver := range []Version{Version20, Version30, Version31, Version32} {
		resources = append(resources, metaSchemaResource{
			file: metaRootSchemaFileForVersion(ver),
			url:  metaRootSchemaURLForVersion(ver),
		})
		for _, frag := range fragmentTypesForVersion(ver) {
			file, ok := metaFragmentSchemaFileForVersion(ver, frag)
			if !ok {
				continue
			}
			url, ok := metaFragmentSchemaURLForVersion(ver, frag)
			if !ok {
				continue
			}
			resources = append(resources, metaSchemaResource{file: file, url: url})
		}
	}
	return resources
}

func metaRootSchemaFileForVersion(ver Version) string {
	return "schemas/meta/openapi-" + string(ver) + "-root.schema.json"
}

func metaRootSchemaURLForVersion(ver Version) string {
	return "https://navigator.local/meta/openapi-" + string(ver) + "-root.json"
}

func metaFragmentSchemaFileForVersion(ver Version, frag FragmentType) (string, bool) {
	name, ok := metaFragmentName(frag)
	if !ok {
		return "", false
	}
	return "schemas/meta/openapi-" + string(ver) + "-" + name + ".schema.json", true
}

func metaFragmentSchemaURLForVersion(ver Version, frag FragmentType) (string, bool) {
	name, ok := metaFragmentName(frag)
	if !ok {
		return "", false
	}
	return "https://navigator.local/meta/openapi-" + string(ver) + "-" + name + ".json", true
}

func fragmentTypesForVersion(ver Version) []FragmentType {
	switch ver {
	case Version20:
		return []FragmentType{
			FragmentSchema,
			FragmentPathItem,
			FragmentOperation,
			FragmentParameter,
			FragmentResponse,
			FragmentHeader,
		}
	case Version30, Version31, Version32:
		return []FragmentType{
			FragmentSchema,
			FragmentPathItem,
			FragmentOperation,
			FragmentParameter,
			FragmentRequestBody,
			FragmentResponse,
			FragmentHeader,
			FragmentSecurityScheme,
			FragmentComponents,
			FragmentServer,
		}
	default:
		return nil
	}
}

func metaFragmentName(frag FragmentType) (string, bool) {
	switch frag {
	case FragmentSchema:
		return "schema", true
	case FragmentPathItem:
		return "path-item", true
	case FragmentOperation:
		return "operation", true
	case FragmentParameter:
		return "parameter", true
	case FragmentRequestBody:
		return "request-body", true
	case FragmentResponse:
		return "response", true
	case FragmentHeader:
		return "header", true
	case FragmentSecurityScheme:
		return "security-scheme", true
	case FragmentComponents:
		return "components", true
	case FragmentServer:
		return "server", true
	default:
		return "", false
	}
}
