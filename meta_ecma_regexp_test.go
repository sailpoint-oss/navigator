package navigator

import "testing"

func TestCompileMetaJSONSchemaPattern_MatchesAndStringifies(t *testing.T) {
	re, err := compileMetaJSONSchemaPattern(`^user-[0-9]+$`)
	if err != nil {
		t.Fatalf("compileMetaJSONSchemaPattern: %v", err)
	}
	if !re.MatchString("user-42") {
		t.Fatal("expected regex to match valid value")
	}
	if re.MatchString("admin") {
		t.Fatal("expected regex not to match invalid value")
	}
	if re.String() != "^user-[0-9]+$" {
		t.Fatalf("String() = %q", re.String())
	}
}

func TestCompileMetaJSONSchemaPattern_InvalidPattern(t *testing.T) {
	if _, err := compileMetaJSONSchemaPattern(`(`); err == nil {
		t.Fatal("expected invalid ECMAScript regex to fail")
	}
}
