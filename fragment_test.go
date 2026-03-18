package navigator

import "testing"

func TestDetectFragmentType(t *testing.T) {
	tests := []struct {
		name string
		keys []string
		want FragmentType
	}{
		{"root not fragment", []string{"openapi", "info", "paths"}, FragmentUnknown},
		{"schema", []string{"type", "properties"}, FragmentSchema},
		{"schema with allOf", []string{"allOf"}, FragmentSchema},
		{"schema with $ref", []string{"$ref"}, FragmentSchema},
		{"path item", []string{"get", "post"}, FragmentPathItem},
		{"operation", []string{"operationId", "responses"}, FragmentOperation},
		{"parameter", []string{"in", "name", "schema"}, FragmentParameter},
		{"components", []string{"components"}, FragmentComponents},
		{"server", []string{"url", "description"}, FragmentServer},
		{"unknown", []string{"random", "keys"}, FragmentUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectFragmentType(tt.keys)
			if got != tt.want {
				t.Errorf("DetectFragmentType(%v) = %v, want %v", tt.keys, got, tt.want)
			}
		})
	}
}

func TestFragmentType_String(t *testing.T) {
	tests := []struct {
		ft   FragmentType
		want string
	}{
		{FragmentSchema, "schema"},
		{FragmentPathItem, "path-item"},
		{FragmentOperation, "operation"},
		{FragmentParameter, "parameter"},
		{FragmentRequestBody, "request-body"},
		{FragmentResponse, "response"},
		{FragmentHeader, "header"},
		{FragmentSecurityScheme, "security-scheme"},
		{FragmentComponents, "components"},
		{FragmentServer, "server"},
		{FragmentUnknown, "unknown"},
	}
	for _, tt := range tests {
		if got := tt.ft.String(); got != tt.want {
			t.Errorf("FragmentType(%d).String() = %q, want %q", tt.ft, got, tt.want)
		}
	}
}
