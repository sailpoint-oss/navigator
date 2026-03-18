package navigator

import (
	"testing"
)

func TestE2E_BarometerPattern(t *testing.T) {
	content := loadFixture(t, "single/minimal.yaml")
	idx := Parse(content)
	if idx == nil {
		t.Fatal("Parse failed")
	}

	// Iterate all operations (what barometer does)
	ops := idx.AllOperations()
	if len(ops) == 0 {
		t.Fatal("no operations found")
	}

	for _, op := range ops {
		if op.Operation.OperationID == "" {
			t.Error("operation missing operationId")
		}
	}

	// Resolve all refs
	for _, ref := range idx.AllRefs {
		val, err := idx.ResolveRef(ref.Target)
		if err != nil {
			t.Errorf("failed to resolve %q: %v", ref.Target, err)
		}
		if val == nil {
			t.Errorf("resolved nil for %q", ref.Target)
		}
	}

	// Verify schema structure
	pet := idx.Schemas["Pet"]
	if pet == nil {
		t.Fatal("Pet schema not found")
	}
	if pet.Type != "object" {
		t.Errorf("Pet.type = %q", pet.Type)
	}
	if _, ok := pet.Properties["name"]; !ok {
		t.Error("Pet missing name property")
	}
	if len(pet.Required) != 1 || pet.Required[0] != "name" {
		t.Errorf("Pet.required = %v", pet.Required)
	}
}

func TestE2E_CartographerPattern(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"api.yaml": `openapi: "3.0.3"
info:
  title: Carto Test
  version: "1.0"
paths:
  /users:
    get:
      operationId: listUsers
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "./schemas/User.yaml"
  /orders:
    get:
      operationId: listOrders
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                $ref: "./schemas/Order.yaml"`,
		"schemas/User.yaml": `type: object
properties:
  id:
    type: integer
  name:
    type: string
required:
  - name`,
		"schemas/Order.yaml": `type: object
properties:
  id:
    type: integer
  customer:
    $ref: "./User.yaml"
  total:
    type: number`,
	})

	proj, err := OpenProject(dir + "/api.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// LoadAll (cartographer needs all files for bundling)
	if err := proj.LoadAll(); err != nil {
		t.Fatal(err)
	}

	// Iterate all operations
	rootIdx := proj.Index()
	ops := rootIdx.AllOperations()
	if len(ops) != 2 {
		t.Errorf("expected 2 operations, got %d", len(ops))
	}

	// Access all cached indexes (for bundling)
	allIndexes := proj.Cache().All()
	if len(allIndexes) < 3 {
		t.Errorf("expected at least 3 cached indexes, got %d", len(allIndexes))
	}

	// Resolve cross-file refs
	result, err := proj.Resolve(proj.RootURI(), "./schemas/User.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if result.TargetIndex == nil {
		t.Error("User.yaml should be resolved")
	}

	// Resolve transitive ref (Order -> User)
	orderURI := ""
	for uri := range allIndexes {
		if idx := allIndexes[uri]; idx != nil {
			for _, ref := range idx.AllRefs {
				if ref.Target == "./User.yaml" && uri != proj.RootURI() {
					orderURI = uri
				}
			}
		}
	}
	if orderURI != "" {
		result, err := proj.Resolve(orderURI, "./User.yaml")
		if err != nil {
			t.Errorf("transitive resolve failed: %v", err)
		}
		if result == nil {
			t.Error("transitive resolve returned nil")
		}
	}
}

func TestE2E_TelescopeLSPPattern(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"root-a.yaml": `openapi: "3.0.3"
info:
  title: Root A
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
                $ref: "./shared/common.yaml"`,
		"root-b.yaml": `openapi: "3.0.3"
info:
  title: Root B
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
                $ref: "./shared/common.yaml"`,
		"shared/common.yaml": `type: object
properties:
  id:
    type: integer
  name:
    type: string`,
	})

	ws := NewWorkspace()

	// Open multiple projects sharing files
	projA, err := ws.OpenProject(dir + "/root-a.yaml")
	if err != nil {
		t.Fatal(err)
	}
	projB, err := ws.OpenProject(dir + "/root-b.yaml")
	if err != nil {
		t.Fatal(err)
	}

	projA.LoadAll()
	projB.LoadAll()

	sharedURI := PathToURI(dir + "/shared/common.yaml")

	// Verify shared file is parsed once (same pointer in cache)
	idx1 := ws.Cache.Get(sharedURI)
	idx2 := ws.Cache.Get(sharedURI)
	if idx1 != idx2 {
		t.Error("shared file should use same cached index")
	}

	// Simulate file edit: invalidate shared file
	affected := projA.OnFileChanged(sharedURI)

	// Both roots should be affected
	rootAURI := projA.RootURI()
	rootBURI := projB.RootURI()
	affectedSet := make(map[string]bool)
	for _, uri := range affected {
		affectedSet[uri] = true
	}
	if !affectedSet[rootAURI] {
		t.Error("root-a should be affected by shared file change")
	}
	if !affectedSet[rootBURI] {
		t.Error("root-b should be affected by shared file change")
	}

	// Cross-file resolution should still work
	result, err := projA.Resolve(projA.RootURI(), "./shared/common.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if result.TargetIndex == nil {
		t.Error("cross-file resolution should work after invalidation")
	}
}

func TestE2E_CycleHandling(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"root.yaml": `openapi: "3.0.3"
info:
  title: Cycle
  version: "1.0"
paths:
  /cycle:
    get:
      operationId: getCycle
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                $ref: "./a.yaml"`,
		"a.yaml": `type: object
properties:
  partner:
    $ref: "./b.yaml"`,
		"b.yaml": `type: object
properties:
  partner:
    $ref: "./a.yaml"`,
	})

	proj, err := OpenProject(dir + "/root.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// LoadAll should not hang on cycles
	if err := proj.LoadAll(); err != nil {
		t.Fatal(err)
	}

	// Graph should detect cycle
	if !proj.Graph().HasCycle() {
		t.Error("should detect cycle between a.yaml and b.yaml")
	}

	// Invalidation should not hang on cycles
	aURI := PathToURI(dir + "/a.yaml")
	affected := proj.OnFileChanged(aURI)
	_ = affected // just verify it completes
}
