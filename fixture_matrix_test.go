package navigator

import (
	"reflect"
	"testing"
)

func TestToolchainFixtureMatrix_VersionParity(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		uri         string
		wantVersion Version
	}{
		{
			name:        "swagger 2.0 root",
			path:        "versions/swagger-2.0.yaml",
			uri:         "file:///versions/swagger-2.0.yaml",
			wantVersion: Version20,
		},
		{
			name:        "openapi 3.0 root",
			path:        "versions/openapi-3.0.yaml",
			uri:         "file:///versions/openapi-3.0.yaml",
			wantVersion: Version30,
		},
		{
			name:        "openapi 3.1 root",
			path:        "versions/openapi-3.1.yaml",
			uri:         "file:///versions/openapi-3.1.yaml",
			wantVersion: Version31,
		},
		{
			name:        "openapi 3.2 root",
			path:        "versions/openapi-3.2.yaml",
			uri:         "file:///versions/openapi-3.2.yaml",
			wantVersion: Version32,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			content := loadFixture(t, tc.path)
			standalone := ParseURI(tc.uri, content)
			treesitter := ParseContent(content, tc.uri)
			if standalone == nil || treesitter == nil {
				t.Fatalf("both parsers must succeed: standalone=%v treesitter=%v", standalone != nil, treesitter != nil)
			}

			if standalone.Version != tc.wantVersion {
				t.Fatalf("standalone version = %q, want %q", standalone.Version, tc.wantVersion)
			}
			if treesitter.Version != tc.wantVersion {
				t.Fatalf("treesitter version = %q, want %q", treesitter.Version, tc.wantVersion)
			}
			if standalone.Document == nil || standalone.Document.DocType != DocTypeRoot {
				t.Fatalf("standalone doc type = %+v, want root", standalone.Document)
			}
			if treesitter.Document == nil || treesitter.Document.DocType != DocTypeRoot {
				t.Fatalf("treesitter doc type = %+v, want root", treesitter.Document)
			}
			if !reflect.DeepEqual(snapshotForParity(standalone), snapshotForParity(treesitter)) {
				t.Fatalf("parity mismatch\nstandalone=%+v\ntreesitter=%+v", snapshotForParity(standalone), snapshotForParity(treesitter))
			}
			assertNoErrorSeverity(t, standalone.Issues)
			assertNoErrorSeverity(t, treesitter.Issues)
		})
	}
}

func TestToolchainFixtureMatrix_InvalidFragmentSchemaFixture(t *testing.T) {
	idx := ParseContent(loadFixture(t, "meta/invalid_fragment_schema.yaml"), "file:///meta/invalid_fragment_schema.yaml")
	if idx == nil {
		t.Fatal("nil index")
	}
	if idx.Document == nil {
		t.Fatal("expected parsed fragment document")
	}
	if idx.Document.DocType != DocTypeFragment {
		t.Fatalf("doc type = %q, want fragment", idx.Document.DocType)
	}

	var sawMetaError bool
	for _, iss := range idx.Issues {
		if iss.Category == CategoryMeta && iss.Severity == SeverityError {
			sawMetaError = true
			break
		}
	}
	if !sawMetaError {
		t.Fatalf("expected at least one fragment meta-schema error, got: %s", formatIssues(idx.Issues))
	}
}
