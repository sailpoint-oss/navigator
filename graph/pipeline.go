package graph

import (
	"context"
	"time"

	navigator "github.com/sailpoint-oss/navigator"
)

// Stage defines a processing stage in the pipeline.
type Stage interface {
	Name() StageName
	DependsOn() []StageName
	Run(ctx context.Context, g *WorkspaceGraph, uri string) (*StageResult, error)
}

// PipelineRunner executes stages in topological order.
type PipelineRunner struct {
	stages []Stage
	order  []StageName
	byName map[StageName]Stage
}

// NewPipelineRunner creates a runner with the given stages.
func NewPipelineRunner(stages ...Stage) *PipelineRunner {
	pr := &PipelineRunner{
		stages: stages,
		byName: make(map[StageName]Stage),
	}
	for _, s := range stages {
		pr.byName[s.Name()] = s
	}
	pr.order = topoSort(stages)
	return pr
}

// AddStage adds a custom stage to the pipeline.
func (pr *PipelineRunner) AddStage(s Stage) {
	pr.stages = append(pr.stages, s)
	pr.byName[s.Name()] = s
	pr.order = topoSort(pr.stages)
}

// Stages returns the ordered list of stage names.
func (pr *PipelineRunner) Stages() []StageName {
	return pr.order
}

// RunStage runs a single named stage for a URI.
func (pr *PipelineRunner) RunStage(ctx context.Context, g *WorkspaceGraph, uri string, name StageName) error {
	return pr.runStage(ctx, g, uri, name)
}

// RunAll runs all dirty stages for a URI in topological order.
func (pr *PipelineRunner) RunAll(ctx context.Context, g *WorkspaceGraph, uri string) error {
	for _, name := range pr.order {
		if g.IsStageDirty(uri, name) {
			if err := pr.runStage(ctx, g, uri, name); err != nil {
				return err
			}
		}
	}
	return nil
}

// RunThrough runs all stages up to and including the given stage.
func (pr *PipelineRunner) RunThrough(ctx context.Context, g *WorkspaceGraph, uri string, through StageName) error {
	for _, name := range pr.order {
		if g.IsStageDirty(uri, name) {
			if err := pr.runStage(ctx, g, uri, name); err != nil {
				return err
			}
		}
		if name == through {
			break
		}
	}
	return nil
}

func (pr *PipelineRunner) runStage(ctx context.Context, g *WorkspaceGraph, uri string, name StageName) error {
	stage, ok := pr.byName[name]
	if !ok {
		return nil
	}
	result, err := stage.Run(ctx, g, uri)
	if err != nil {
		return err
	}
	if result != nil {
		g.SetStageResult(uri, result)
	}
	return nil
}

// DefaultStages returns Navigator's shared three-stage substrate.
// Downstream lint/validate/analyze stages are owned by higher-level
// orchestrators such as Telescope and Barrelman.
func DefaultStages() []Stage {
	return []Stage{
		&rawStage{},
		&parseStage{},
		&bindStage{},
	}
}

// --- Built-in stages ---

type rawStage struct{}

func (s *rawStage) Name() StageName        { return StageRaw }
func (s *rawStage) DependsOn() []StageName { return nil }
func (s *rawStage) Run(ctx context.Context, g *WorkspaceGraph, uri string) (*StageResult, error) {
	start := time.Now()
	node := g.Node(uri)
	if node == nil || node.Source == nil {
		return nil, nil
	}
	content, version, err := node.Source.Read(ctx)
	if err != nil {
		return nil, err
	}
	g.mu.Lock()
	node.Raw = content
	node.Version = version
	g.mu.Unlock()
	return &StageResult{Stage: StageRaw, Data: content, Version: version, Duration: time.Since(start)}, nil
}

type parseStage struct{}

func (s *parseStage) Name() StageName        { return StageParse }
func (s *parseStage) DependsOn() []StageName { return []StageName{StageRaw} }
func (s *parseStage) Run(ctx context.Context, g *WorkspaceGraph, uri string) (*StageResult, error) {
	start := time.Now()
	node := g.Node(uri)
	if node == nil || node.Raw == nil {
		return nil, nil
	}
	idx := navigator.ParseURI(uri, node.Raw)
	if idx == nil {
		return nil, nil
	}
	g.cache.Set(uri, idx)
	g.mu.Lock()
	node.Index = idx
	g.mu.Unlock()
	return &StageResult{Stage: StageParse, Data: idx, Version: node.Version, Duration: time.Since(start)}, nil
}

type bindStage struct{}

func (s *bindStage) Name() StageName        { return StageBind }
func (s *bindStage) DependsOn() []StageName { return []StageName{StageParse} }
func (s *bindStage) Run(_ context.Context, g *WorkspaceGraph, uri string) (*StageResult, error) {
	start := time.Now()
	node := g.Node(uri)
	if node == nil || node.Index == nil {
		return nil, nil
	}
	navigator.UpdateGraphFromIndex(g.graph, uri, node.Index)
	if node.Index.IsOpenAPI() {
		g.SetRoot(uri, true)
	}
	return &StageResult{Stage: StageBind, Version: node.Version, Duration: time.Since(start)}, nil
}

// --- topological sort ---

func topoSort(stages []Stage) []StageName {
	inDegree := make(map[StageName]int)
	deps := make(map[StageName][]StageName)
	all := make(map[StageName]bool)

	for _, s := range stages {
		name := s.Name()
		all[name] = true
		for _, d := range s.DependsOn() {
			deps[d] = append(deps[d], name)
			inDegree[name]++
		}
		if _, ok := inDegree[name]; !ok {
			inDegree[name] = 0
		}
	}

	var queue []StageName
	for name := range all {
		if inDegree[name] == 0 {
			queue = append(queue, name)
		}
	}

	var order []StageName
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		order = append(order, name)
		for _, dependent := range deps[name] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}
	return order
}
