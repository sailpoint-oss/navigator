package navigator

// Schema describes a data type used in OpenAPI.
type Schema struct {
	Type                       string
	Format                     string
	Title                      string
	Description                DescriptionValue
	Default                    *Node
	Enum                       []string
	Required                   []string
	Properties                 map[string]*Schema
	AdditionalProperties       *Schema
	AdditionalPropertiesFalse  bool // true when additionalProperties is explicitly false
	UnevaluatedProperties      *Schema
	UnevaluatedPropertiesFalse bool // true when unevaluatedProperties is explicitly false
	AllOf                      []*Schema
	AnyOf                      []*Schema
	OneOf                      []*Schema
	Not                        *Schema
	Items                      *Schema
	MinLength                  *int
	MaxLength                  *int
	Minimum                    *float64
	Maximum                    *float64
	ExclusiveMinimum           *float64 // OAS 3.1+: numeric value
	ExclusiveMaximum           *float64 // OAS 3.1+: numeric value
	MinItems                   *int
	MaxItems                   *int
	MaxProperties              *int
	Pattern                    string
	HasConst                   bool // true when const is present
	Nullable                   bool
	ReadOnly                   bool
	WriteOnly                  bool
	Deprecated                 bool
	Example                    *Node
	ExternalDocs               *ExternalDocs
	Discriminator              *Discriminator
	Ref                        string
	Extensions                 map[string]*Node
	Loc                        Loc
	TypeLoc                    Loc
	NameLoc                    Loc
}

// Discriminator aids in serialization, deserialization, and validation of
// request/response bodies when using oneOf or anyOf.
type Discriminator struct {
	PropertyName string
	Mapping      map[string]string
	Loc          Loc
}
