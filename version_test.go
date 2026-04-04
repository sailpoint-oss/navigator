package navigator

import "testing"

func TestVersionFromString(t *testing.T) {
	tests := []struct {
		input string
		want  Version
	}{
		{"3.0.0", Version30},
		{"3.1.0", Version31},
		{"3.2.0", Version32},
		{"2.0", Version20},
		{"", VersionUnknown},
	}
	for _, tt := range tests {
		got := VersionFromString(tt.input)
		if got != tt.want {
			t.Errorf("VersionFromString(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestVersion_AtLeast(t *testing.T) {
	v := Version31
	if !v.AtLeast(3, 0, 0) {
		t.Error("3.1 should be >= 3.0.0")
	}
	if !v.AtLeast(3, 1, 0) {
		t.Error("3.1 should be >= 3.1.0")
	}
	if v.AtLeast(3, 2, 0) {
		t.Error("3.1 should not be >= 3.2.0")
	}
}

func TestVersion_IsZero(t *testing.T) {
	if !VersionUnknown.IsZero() {
		t.Error("unknown version should be zero")
	}
	if Version30.IsZero() {
		t.Error("3.0 should not be zero")
	}
}

func TestFormatsForVersion(t *testing.T) {
	tests := []struct {
		v    Version
		want []Format
	}{
		{Version20, []Format{FormatOAS2}},
		{Version30, []Format{FormatOAS3}},
		{Version31, []Format{FormatOAS3, FormatOAS31}},
		{Version32, []Format{FormatOAS3, FormatOAS32}},
		{VersionUnknown, nil},
	}
	for _, tt := range tests {
		got := FormatsForVersion(tt.v)
		if len(got) != len(tt.want) {
			t.Errorf("FormatsForVersion(%q) = %v, want %v", tt.v, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("FormatsForVersion(%q)[%d] = %q, want %q", tt.v, i, got[i], tt.want[i])
			}
		}
	}
}

func TestDocType_String(t *testing.T) {
	tests := []struct {
		dt   DocType
		want string
	}{
		{DocTypeRoot, "root"},
		{DocTypeFragment, "fragment"},
		{DocTypeNonOpenAPI, "non-openapi"},
		{DocTypeUnknown, "unknown"},
	}
	for _, tt := range tests {
		if got := tt.dt.String(); got != tt.want {
			t.Errorf("DocType(%d).String() = %q, want %q", tt.dt, got, tt.want)
		}
	}
}

func TestDetectDocType(t *testing.T) {
	tests := []struct {
		name          string
		keys          []string
		isGraphMember bool
		want          DocType
	}{
		{"root from openapi key", []string{"openapi", "info", "paths"}, false, DocTypeRoot},
		{"root from swagger key", []string{"swagger", "info"}, false, DocTypeRoot},
		{"fragment from schema keys", []string{"type", "properties"}, false, DocTypeFragment},
		{"fragment from path-item keys", []string{"get", "post"}, false, DocTypeFragment},
		{"fragment from operation keys", []string{"operationId", "responses"}, false, DocTypeFragment},
		{"fragment via graph membership", []string{"random", "keys"}, true, DocTypeFragment},
		{"non-openapi with unrecognized keys", []string{"random", "keys"}, false, DocTypeNonOpenAPI},
		{"non-openapi empty keys", []string{}, false, DocTypeNonOpenAPI},
		{"root overrides graph membership", []string{"openapi", "info"}, true, DocTypeRoot},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectDocType(tt.keys, tt.isGraphMember)
			if got != tt.want {
				t.Errorf("DetectDocType(%v, %v) = %v, want %v", tt.keys, tt.isGraphMember, got, tt.want)
			}
		})
	}
}
