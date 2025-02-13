package workflow

import (
	"context"
	"github.com/ahaostudy/eino-learn/plugin"
	"github.com/ahaostudy/eino-learn/schema"
	"github.com/cloudwego/eino/compose"
)

type Node struct {
	Name         string
	InputMapper  ObjectMapper
	OutputSchema *schema.Definition

	Batch      bool
	BatchField *ObjectField

	Runner NodeRunner
}

func (node *Node) Lambda() *compose.Lambda {
	if node.Batch {
		return NodeLambdaBatch(node.Name, node.InputMapper, ArraySplitter(node.BatchField), node.Runner.Run)
	}

	return NodeLambda(node.Name, node.InputMapper, node.Runner.Run)
}

type NodeRunner interface {
	Run(ctx context.Context, input Object) (Object, error)
}

type FunctionNodeRunner struct {
	runner func(ctx context.Context, input Object) (Object, error)
}

func NewFunctionNodeRunner(runner func(ctx context.Context, input Object) (Object, error)) *FunctionNodeRunner {
	return &FunctionNodeRunner{runner: runner}
}

func (r *FunctionNodeRunner) Run(ctx context.Context, input Object) (Object, error) {
	return r.runner(ctx, input)
}

type PluginNodeRunner struct {
	pluginID int64
	tooID    int64
	secrets  map[string]string
}

func NewPluginNodeRunner(pluginID, tooID int64, secrets map[string]string) *PluginNodeRunner {
	return &PluginNodeRunner{pluginID, tooID, secrets}
}

func (r *PluginNodeRunner) Run(ctx context.Context, input Object) (Object, error) {
	req, err := input.JSON()
	if err != nil {
		return nil, err
	}

	resp, err := plugin.MockPluginCall(ctx, &plugin.CallPluginToolReq{
		PluginId: r.pluginID,
		ToolId:   r.tooID,
		Secrets:  r.secrets,
		Request:  req,
	})
	if err != nil {
		return nil, err
	}

	output, err := NewJSONObject(resp.Response)
	if err != nil {
		return nil, err
	}

	return output, nil
}
