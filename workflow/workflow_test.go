package workflow

import (
	"context"
	"log"
	"testing"
)

func TestWorkflow(t *testing.T) {
	workflow := NewWorkflow()

	workflow.AppendLambdaNodeStart(func(ctx context.Context, input Object) (Object, error) {
		return input, nil
	})

	workflow.AppendNode(&Node{
		Name: "read_project_structure",
		InputMapper: ObjectMapper{
			{Name: "owner", Node: "start", Path: "owner"},
			{Name: "repo", Node: "start", Path: "repo"},
			{Name: "path", IsConstant: true, Constant: "/"},
			{Name: "recursion", IsConstant: true, Constant: true},
		},
		Runner: NewFunctionNodeRunner(func(ctx context.Context, input Object) (Object, error) {
			return Object{"children": []Object{{"path": "README.md", "type": "file"}}}, nil
		}),
	})

	workflow.AppendNode(&Node{
		Name: "plugin_call",
		InputMapper: ObjectMapper{
			{Name: "owner", Node: "start", Path: "owner"},
			{Name: "repo", Node: "start", Path: "repo"},
			{Name: "path", IsConstant: true, Constant: "/"},
			{Name: "recursion", IsConstant: true, Constant: true},
		},
		Runner: NewPluginNodeRunner(0, 0, map[string]string{"GITHUB_TOKEN": "xxx"}),
	})

	workflow.AppendLambdaNode("extract_target_directory",
		ObjectMapper{
			{Name: "owner", Node: "start", Path: "owner"},
			{Name: "repo", Node: "start", Path: "repo"},
			{Name: "query", Node: "start", Path: "GITHUB_TOKEN"},
			{Name: "tree", Node: "read_project_structure", Path: "children"},
		},
		func(ctx context.Context, input Object) (Object, error) {
			return Object{
				"dirs": []string{"/", "/apps/"},
			}, nil
		})

	workflow.AppendLambdaNodeBatch(
		"read_target_trees",
		ObjectMapper{
			{Name: "dirs", Node: "extract_target_directory", Path: "dirs"},
			{Name: "owner", Node: "start", Path: "owner"},
			{Name: "repo", Node: "start", Path: "repo"},
			{Name: "recursion", IsConstant: true, Constant: true},
		},
		ArraySplitter(&ObjectField{
			Name: "dir",
			Node: "extract_target_directory",
			Path: "dirs",
		}),
		func(ctx context.Context, input Object) (Object, error) {
			return Object{
				"children": []Object{{"path": "README.md", "type": "file"}},
			}, nil
		},
	)

	workflow.AppendLambdaNode("extract_files_by_trees",
		ObjectMapper{
			{Name: "trees", Node: "read_target_trees", IsBatchOutput: true},
		},
		func(ctx context.Context, input Object) (Object, error) {
			trees := input.ObjectArray("trees")

			var files []string
			for _, tree := range trees {
				for _, child := range tree.ObjectArray("children") {
					if child.String("type") == "file" {
						files = append(files, child.String("path"))
					}
				}
			}

			return Object{"files": files}, nil
		})

	workflow.AppendLambdaNode("read_files_content",
		ObjectMapper{
			{Name: "owner", Node: "start", Path: "owner"},
			{Name: "repo", Node: "start", Path: "repo"},
			{Name: "files", Node: "extract_files_by_trees", Path: "files"},
		},
		func(ctx context.Context, input Object) (Object, error) {
			return Object{
				"files": Object{
					"path":    "README.md",
					"content": "this is README.",
				},
			}, nil
		})

	workflow.AppendLambdaNodeEnd(
		ObjectMapper{
			{Name: "files", Node: "read_files_content", Path: "files"},
		},
		func(ctx context.Context, input Object) (Object, error) {
			return input.Object("files"), nil
		})

	output, err := workflow.Run(context.TODO(), Object{
		"owner": "aiagt",
		"repo":  "aiagt",
		"query": "any question",
	})
	if err != nil {
		panic(err)
	}

	log.Println("output:", output.Pretty())
}
