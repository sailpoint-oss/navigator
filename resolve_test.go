package navigator

import "testing"

func TestResolveRef_Schema(t *testing.T) {
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	val, err := idx.ResolveRef("#/components/schemas/Pet")
	if err != nil {
		t.Fatal(err)
	}
	s, ok := val.(*Schema)
	if !ok {
		t.Fatalf("expected *Schema, got %T", val)
	}
	if s.Type != "object" {
		t.Errorf("schema type = %q, want %q", s.Type, "object")
	}
}

func TestResolveRef_Parameter(t *testing.T) {
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	val, err := idx.ResolveRef("#/components/parameters/PetId")
	if err != nil {
		t.Fatal(err)
	}
	p, ok := val.(*Parameter)
	if !ok {
		t.Fatalf("expected *Parameter, got %T", val)
	}
	if p.Name != "petId" {
		t.Errorf("param name = %q, want %q", p.Name, "petId")
	}
}

func TestResolveRef_Response(t *testing.T) {
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	val, err := idx.ResolveRef("#/components/responses/NotFound")
	if err != nil {
		t.Fatal(err)
	}
	r, ok := val.(*Response)
	if !ok {
		t.Fatalf("expected *Response, got %T", val)
	}
	if r.Description.Text != "Resource not found" {
		t.Errorf("response description = %q", r.Description.Text)
	}
}

func TestResolveRef_SecurityScheme(t *testing.T) {
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	val, err := idx.ResolveRef("#/components/securitySchemes/BearerAuth")
	if err != nil {
		t.Fatal(err)
	}
	ss, ok := val.(*SecurityScheme)
	if !ok {
		t.Fatalf("expected *SecurityScheme, got %T", val)
	}
	if ss.Type != "http" {
		t.Errorf("scheme type = %q, want %q", ss.Type, "http")
	}
}

func TestResolveRef_Path(t *testing.T) {
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	val, err := idx.ResolveRef("#/paths/~1pets")
	if err != nil {
		t.Fatal(err)
	}
	_, ok := val.(*PathItem)
	if !ok {
		t.Fatalf("expected *PathItem, got %T", val)
	}
}

func TestResolveRef_Operation(t *testing.T) {
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	val, err := idx.ResolveRef("#/paths/~1pets/get")
	if err != nil {
		t.Fatal(err)
	}
	op, ok := val.(*Operation)
	if !ok {
		t.Fatalf("expected *Operation, got %T", val)
	}
	if op.OperationID != "listPets" {
		t.Errorf("operationId = %q", op.OperationID)
	}
}

func TestResolveRef_Info(t *testing.T) {
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	val, err := idx.ResolveRef("#/info")
	if err != nil {
		t.Fatal(err)
	}
	info, ok := val.(*Info)
	if !ok {
		t.Fatalf("expected *Info, got %T", val)
	}
	if info.Title != "Minimal API" {
		t.Errorf("title = %q", info.Title)
	}
}

func TestResolveRef_Empty(t *testing.T) {
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	_, err := idx.ResolveRef("")
	if err == nil {
		t.Error("empty ref should return error")
	}
}

func TestResolveRef_NotFound(t *testing.T) {
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	_, err := idx.ResolveRef("#/components/schemas/NonExistent")
	if err == nil {
		t.Error("non-existent ref should return error")
	}
}

func TestEscapeJSONPointer(t *testing.T) {
	if got := EscapeJSONPointer("pets/{petId}"); got != "pets~1{petId}" {
		t.Errorf("EscapeJSONPointer = %q", got)
	}
	if got := EscapeJSONPointer("a~b"); got != "a~0b" {
		t.Errorf("EscapeJSONPointer = %q", got)
	}
}

func TestComponentRefPath(t *testing.T) {
	if got := ComponentRefPath("schemas", "Pet"); got != "#/components/schemas/Pet" {
		t.Errorf("ComponentRefPath = %q", got)
	}
}
