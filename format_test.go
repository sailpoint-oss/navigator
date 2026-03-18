package navigator

import "testing"

func TestFormatFromURI(t *testing.T) {
	tests := []struct {
		uri  string
		want FileFormat
	}{
		{"file:///api.yaml", FormatYAML},
		{"file:///api.yml", FormatYAML},
		{"file:///api.json", FormatJSON},
		{"file:///api.YAML", FormatYAML},
		{"file:///api.JSON", FormatJSON},
		{"file:///api.txt", FormatUnknown},
		{"", FormatUnknown},
	}
	for _, tt := range tests {
		if got := FormatFromURI(tt.uri); got != tt.want {
			t.Errorf("FormatFromURI(%q) = %v, want %v", tt.uri, got, tt.want)
		}
	}
}

func TestFormatFromContent(t *testing.T) {
	if FormatFromContent([]byte(`{"openapi":"3.0"}`)) != FormatJSON {
		t.Error("JSON content not detected")
	}
	if FormatFromContent([]byte("openapi: '3.0'")) != FormatYAML {
		t.Error("YAML content not detected")
	}
	if FormatFromContent([]byte("")) != FormatUnknown {
		t.Error("empty content should be unknown")
	}
}

func TestDetectFormat(t *testing.T) {
	if DetectFormat("api.yaml", nil) != FormatYAML {
		t.Error("should detect from URI")
	}
	if DetectFormat("", []byte(`{"a":1}`)) != FormatJSON {
		t.Error("should fallback to content sniffing")
	}
}

func TestFileFormat_String(t *testing.T) {
	if FormatYAML.String() != "yaml" {
		t.Error("YAML string wrong")
	}
	if FormatJSON.String() != "json" {
		t.Error("JSON string wrong")
	}
	if FormatUnknown.String() != "unknown" {
		t.Error("Unknown string wrong")
	}
}
