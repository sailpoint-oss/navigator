package navigator

import (
	"path/filepath"
	"sort"
	"testing"
)

func mustTestdataURI(t *testing.T, rel string) string {
	t.Helper()
	abs, err := resolveAbsPath(rel)
	if err != nil {
		t.Fatal(err)
	}
	return PathToURI(abs)
}

func TestOpenProject_SingleFile(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"api.yaml": string(loadFixture(t, "single/minimal.yaml")),
	})

	proj, err := OpenProject(dir + "/api.yaml")
	if err != nil {
		t.Fatal(err)
	}

	idx := proj.Index()
	if idx == nil {
		t.Fatal("root index is nil")
	}
	if !idx.IsOpenAPI() {
		t.Error("root should be OpenAPI")
	}
	if len(idx.Operations) != 3 {
		t.Errorf("expected 3 operations, got %d", len(idx.Operations))
	}
}

func TestOpenProject_MultiFile(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"api.yaml": `openapi: "3.0.3"
info:
  title: Multi
  version: "1.0"
paths:
  /pets:
    get:
      operationId: listPets
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                $ref: "./schemas/Pet.yaml"
components:
  schemas:
    Pet:
      $ref: "./schemas/Pet.yaml"`,
		"schemas/Pet.yaml": `type: object
properties:
  name:
    type: string`,
	})

	proj, err := OpenProject(dir + "/api.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if err := proj.LoadAll(); err != nil {
		t.Fatal(err)
	}

	uris := proj.AllURIs()
	if len(uris) < 2 {
		t.Errorf("expected at least 2 URIs, got %d", len(uris))
	}
}

func TestOpenProject_MultiFileFixture(t *testing.T) {
	rootPath := filepath.Join("testdata", "multifile", "v1", "api.yaml")
	proj, err := OpenProject(rootPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := proj.LoadAll(); err != nil {
		t.Fatal(err)
	}

	uris := proj.AllURIs()
	sort.Strings(uris)
	if len(uris) != 6 {
		t.Fatalf("expected 6 URIs in multifile fixture, got %d: %v", len(uris), uris)
	}

	paramURI := mustTestdataURI(t, filepath.Join("testdata", "multifile", "v1", "components", "parameters.yaml"))
	if !proj.ContainsFile(paramURI) {
		t.Errorf("expected project to contain %q", paramURI)
	}
}

func TestProject_LazyLoading(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"api.yaml": `openapi: "3.0.3"
info:
  title: Lazy
  version: "1.0"
paths:
  /pets:
    get:
      operationId: listPets
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                $ref: "./schemas/Pet.yaml"`,
		"schemas/Pet.yaml": `type: object
properties:
  name:
    type: string`,
	})

	proj, err := OpenProject(dir + "/api.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Before LoadAll, only root should be cached
	if proj.Cache().Len() != 1 {
		t.Errorf("before LoadAll, cache len = %d, want 1", proj.Cache().Len())
	}

	// Resolving a ref should trigger lazy loading
	result, err := proj.Resolve(proj.RootURI(), "./schemas/Pet.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if result.TargetIndex == nil {
		t.Error("lazy load should produce target index")
	}

	if proj.Cache().Len() < 2 {
		t.Error("after resolve, cache should have at least 2 entries")
	}
}

func TestProject_OnFileChanged(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"api.yaml": `openapi: "3.0.3"
info:
  title: Change
  version: "1.0"
paths:
  /pets:
    get:
      operationId: listPets
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                $ref: "./schemas/Pet.yaml"`,
		"schemas/Pet.yaml": `type: object
properties:
  name:
    type: string`,
	})

	proj, err := OpenProject(dir + "/api.yaml")
	if err != nil {
		t.Fatal(err)
	}

	proj.LoadAll()

	petURI := PathToURI(dir + "/schemas/Pet.yaml")
	affected := proj.OnFileChanged(petURI)

	if len(affected) == 0 {
		t.Error("changing Pet.yaml should affect the root")
	}
}

func TestProject_Diamond(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"root.yaml": `openapi: "3.0.3"
info:
  title: Diamond
  version: "1.0"
paths:
  /a:
    get:
      operationId: getA
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                $ref: "./b.yaml"
  /c:
    get:
      operationId: getC
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                $ref: "./c.yaml"`,
		"b.yaml": `type: object
properties:
  shared:
    $ref: "./d.yaml"`,
		"c.yaml": `type: object
properties:
  shared:
    $ref: "./d.yaml"`,
		"d.yaml": `type: object
properties:
  name:
    type: string`,
	})

	proj, err := OpenProject(dir + "/root.yaml")
	if err != nil {
		t.Fatal(err)
	}

	proj.LoadAll()

	uris := proj.AllURIs()
	sort.Strings(uris)
	if len(uris) != 4 {
		t.Errorf("expected 4 URIs in diamond project, got %d: %v", len(uris), uris)
	}

	dURI := PathToURI(dir + "/d.yaml")
	affected := proj.OnFileChanged(dURI)
	if len(affected) < 2 {
		t.Errorf("changing d should affect at least b and c, got %d affected", len(affected))
	}
}

func TestProject_Resolve_MultiFileFixtureWholeFileAndFragmentRefs(t *testing.T) {
	rootPath := filepath.Join("testdata", "multifile", "v1", "api.yaml")
	proj, err := OpenProject(rootPath)
	if err != nil {
		t.Fatal(err)
	}

	pathItemResult, err := proj.Resolve(proj.RootURI(), "./paths/pets.yaml")
	if err != nil {
		t.Fatal(err)
	}
	pathItem, ok := pathItemResult.Value.(*PathItem)
	if !ok {
		t.Fatalf("expected *PathItem, got %T", pathItemResult.Value)
	}
	if pathItem.Get == nil || pathItem.Get.OperationID != "listPets" {
		t.Fatalf("expected pets path item get operation, got %#v", pathItem.Get)
	}

	paramResult, err := proj.Resolve(proj.RootURI(), "./components/parameters.yaml#/Limit")
	if err != nil {
		t.Fatal(err)
	}
	param, ok := paramResult.Value.(*Parameter)
	if !ok {
		t.Fatalf("expected *Parameter, got %T", paramResult.Value)
	}
	if param.Name != "limit" {
		t.Errorf("parameter name = %q, want %q", param.Name, "limit")
	}
}

func TestProject_LoadAll_CycleFixtureTerminates(t *testing.T) {
	rootPath := filepath.Join("testdata", "cycles", "cycle-a.yaml")
	proj, err := OpenProject(rootPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := proj.LoadAll(); err != nil {
		t.Fatal(err)
	}

	uris := proj.AllURIs()
	sort.Strings(uris)
	if len(uris) != 2 {
		t.Fatalf("expected 2 URIs in cycle fixture, got %d: %v", len(uris), uris)
	}
	if !proj.Graph().HasCycle() {
		t.Fatal("expected cycle fixture graph to contain a cycle")
	}
}

func TestProject_DiamondFixtureSharedSchema(t *testing.T) {
	rootPath := filepath.Join("testdata", "diamond", "root.yaml")
	proj, err := OpenProject(rootPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := proj.LoadAll(); err != nil {
		t.Fatal(err)
	}

	uris := proj.AllURIs()
	sort.Strings(uris)
	if len(uris) != 4 {
		t.Fatalf("expected 4 URIs in diamond fixture, got %d: %v", len(uris), uris)
	}

	orderURI := mustTestdataURI(t, filepath.Join("testdata", "diamond", "schemas", "Order.yaml"))
	invoiceURI := mustTestdataURI(t, filepath.Join("testdata", "diamond", "schemas", "Invoice.yaml"))
	sharedURI := mustTestdataURI(t, filepath.Join("testdata", "diamond", "shared", "Customer.yaml"))

	fromOrder, err := proj.Resolve(orderURI, "../shared/Customer.yaml")
	if err != nil {
		t.Fatal(err)
	}
	fromInvoice, err := proj.Resolve(invoiceURI, "../shared/Customer.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if fromOrder.TargetURI != sharedURI {
		t.Errorf("order shared target = %q, want %q", fromOrder.TargetURI, sharedURI)
	}
	if fromInvoice.TargetURI != sharedURI {
		t.Errorf("invoice shared target = %q, want %q", fromInvoice.TargetURI, sharedURI)
	}
	if _, ok := fromOrder.Value.(*Schema); !ok {
		t.Fatalf("expected shared schema from order, got %T", fromOrder.Value)
	}
	if _, ok := fromInvoice.Value.(*Schema); !ok {
		t.Fatalf("expected shared schema from invoice, got %T", fromInvoice.Value)
	}
}

func TestOpenProjectFromIndex(t *testing.T) {
	content := loadFixture(t, "single/minimal.yaml")
	idx := ParseURI("file:///test.yaml", content)

	proj := OpenProjectFromIndex("file:///test.yaml", idx)
	if proj.Index() == nil {
		t.Fatal("index should not be nil")
	}
	if !proj.ContainsFile("file:///test.yaml") {
		t.Error("should contain root file")
	}
}

func TestProject_ContainsFile(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"api.yaml": `openapi: "3.0.3"
info:
  title: Test
  version: "1.0"
paths:
  /pets:
    get:
      operationId: listPets
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                $ref: "./schemas/Pet.yaml"`,
		"schemas/Pet.yaml": `type: object
properties:
  name:
    type: string`,
	})

	proj, err := OpenProject(dir + "/api.yaml")
	if err != nil {
		t.Fatal(err)
	}
	proj.LoadAll()

	if !proj.ContainsFile(proj.RootURI()) {
		t.Error("should contain root")
	}
}

func TestWorkspace_MultiRoot(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"a.yaml": `openapi: "3.0.3"
info:
  title: A
  version: "1.0"
paths:
  /a:
    get:
      operationId: getA
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                $ref: "./shared.yaml"`,
		"b.yaml": `openapi: "3.0.3"
info:
  title: B
  version: "1.0"
paths:
  /b:
    get:
      operationId: getB
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                $ref: "./shared.yaml"`,
		"shared.yaml": `type: object
properties:
  name:
    type: string`,
	})

	ws := NewWorkspace()
	projA, err := ws.OpenProject(dir + "/a.yaml")
	if err != nil {
		t.Fatal(err)
	}
	projB, err := ws.OpenProject(dir + "/b.yaml")
	if err != nil {
		t.Fatal(err)
	}

	projA.LoadAll()
	projB.LoadAll()

	// Shared file should be in the same cache
	sharedURI := PathToURI(dir + "/shared.yaml")
	if ws.Cache.Get(sharedURI) == nil {
		t.Error("shared.yaml should be cached")
	}

	// Both projects should see the shared file
	if !projA.ContainsFile(sharedURI) {
		t.Error("projA should contain shared.yaml")
	}
	if !projB.ContainsFile(sharedURI) {
		t.Error("projB should contain shared.yaml")
	}
}
