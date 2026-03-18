package navigator

import "testing"

func TestPathToURI(t *testing.T) {
	uri := PathToURI("/home/user/api.yaml")
	if uri != "file:///home/user/api.yaml" {
		t.Errorf("PathToURI = %q", uri)
	}
}

func TestURIToPath(t *testing.T) {
	p := URIToPath("file:///home/user/api.yaml")
	if p != "/home/user/api.yaml" {
		t.Errorf("URIToPath = %q", p)
	}
}

func TestURIToPath_NonFileURI(t *testing.T) {
	p := URIToPath("/home/user/api.yaml")
	if p != "/home/user/api.yaml" {
		t.Errorf("URIToPath of plain path = %q", p)
	}
}

func TestPathToURI_Roundtrip(t *testing.T) {
	original := "/home/user/project/api.yaml"
	uri := PathToURI(original)
	back := URIToPath(uri)
	if back != original {
		t.Errorf("roundtrip failed: %q -> %q -> %q", original, uri, back)
	}
}

func TestNormalizeURI(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"file:///home//user/api.yaml", "file:///home/user/api.yaml"},
		{"file:///home/user/./api.yaml", "file:///home/user/api.yaml"},
		{"", ""},
		{"not-a-file-uri", "not-a-file-uri"},
	}
	for _, tt := range tests {
		got := NormalizeURI(tt.input)
		if got != tt.want {
			t.Errorf("NormalizeURI(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestResolveRelativeURI(t *testing.T) {
	base := "file:///workspace/api.yaml"

	got := ResolveRelativeURI(base, "./schemas/Pet.yaml")
	want := "file:///workspace/schemas/Pet.yaml"
	if got != want {
		t.Errorf("ResolveRelativeURI(%q, %q) = %q, want %q", base, "./schemas/Pet.yaml", got, want)
	}
}

func TestResolveRelativeURI_WithFragment(t *testing.T) {
	base := "file:///workspace/api.yaml"
	got := ResolveRelativeURI(base, "./schemas/Pet.yaml#/components/schemas/Pet")
	if got == "" {
		t.Fatal("should not be empty")
	}
	if got != "file:///workspace/schemas/Pet.yaml#/components/schemas/Pet" {
		t.Errorf("got %q", got)
	}
}

func TestResolveRelativeURI_EmptyRel(t *testing.T) {
	base := "file:///workspace/api.yaml"
	got := ResolveRelativeURI(base, "#/components/schemas/Pet")
	if got != "file:///workspace/api.yaml#/components/schemas/Pet" {
		t.Errorf("got %q", got)
	}
}

func TestSplitRefURI(t *testing.T) {
	tests := []struct {
		ref          string
		wantFile     string
		wantFragment string
	}{
		{"#/components/schemas/Pet", "", "#/components/schemas/Pet"},
		{"./schemas/Pet.yaml#/components/schemas/Pet", "./schemas/Pet.yaml", "#/components/schemas/Pet"},
		{"./schemas/Pet.yaml", "./schemas/Pet.yaml", ""},
	}
	for _, tt := range tests {
		file, frag := SplitRefURI(tt.ref)
		if file != tt.wantFile || frag != tt.wantFragment {
			t.Errorf("SplitRefURI(%q) = (%q, %q), want (%q, %q)", tt.ref, file, frag, tt.wantFile, tt.wantFragment)
		}
	}
}
