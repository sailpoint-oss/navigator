package navigator

// Components holds reusable objects for different aspects of the OAS.
type Components struct {
	Schemas         map[string]*Schema
	Responses       map[string]*Response
	Parameters      map[string]*Parameter
	Examples        map[string]*Example
	RequestBodies   map[string]*RequestBody
	Headers         map[string]*Header
	SecuritySchemes map[string]*SecurityScheme
	Links           map[string]*Link
	Callbacks       map[string]*Callback
	PathItems       map[string]*PathItem
	Loc             Loc
}

// Example describes an example value.
type Example struct {
	Summary       string
	Description   DescriptionValue
	Value         *Node
	ExternalValue string
	Ref           string
	Loc           Loc
	NameLoc       Loc
}

// SecurityScheme defines a security mechanism.
type SecurityScheme struct {
	Type             string
	Description      DescriptionValue
	Name             string
	In               string
	Scheme           string
	BearerFormat     string
	Flows            *OAuthFlows
	OpenIDConnectURL string
	Ref              string
	Extensions       map[string]*Node
	Loc              Loc
	NameLoc          Loc
}

// OAuthFlows allows configuration of the supported OAuth Flows.
type OAuthFlows struct {
	Implicit          *OAuthFlow
	Password          *OAuthFlow
	ClientCredentials *OAuthFlow
	AuthorizationCode *OAuthFlow
	Loc               Loc
}

// OAuthFlow describes an OAuth 2.0 flow configuration.
type OAuthFlow struct {
	AuthorizationURL    string
	TokenURL            string
	RefreshURL          string
	Scopes              map[string]string
	Loc                 Loc
	AuthorizationURLLoc Loc
	TokenURLLoc         Loc
}

// Callback is a map of runtime expressions to PathItem objects.
type Callback = map[string]*PathItem
