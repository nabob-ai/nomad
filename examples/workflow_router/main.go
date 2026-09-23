// WorkflowRouter 演示 Router 步骤的流式执行，支持基于输入内容动态选择
// 执行路径。
package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/nabob-ai/nomad/pkg/workflow"
)

func main() {
	fmt.Println("=== Nomad Router 流式执行示例 ===")

	ctx := context.Background()

	// 创建一些步骤
	step1 := workflow.NewFunctionStep("analyze", func(ctx context.Context, input *workflow.StepInput) (*workflow.StepOutput, error) {
		fmt.Println("  🔍 Analyzing input...")
		return &workflow.StepOutput{
			Content:  map[string]any{"analysis": "complex", "priority": "high"},
			Metadata: make(map[string]any),
		}, nil
	})

	step2 := workflow.NewFunctionStep("process_complex", func(ctx context.Context, input *workflow.StepInput) (*workflow.StepOutput, error) {
		fmt.Println("  ⚙️  Processing complex case...")
		return &workflow.StepOutput{
			Content:  "Processed with advanced algorithm",
			Metadata: make(map[string]any),
		}, nil
	})

	step3 := workflow.NewFunctionStep("process_simple", func(ctx context.Context, input *workflow.StepInput) (*workflow.StepOutput, error) {
		fmt.Println("  ⚡ Processing simple case...")
		return &workflow.StepOutput{
			Content:  "Processed with basic algorithm",
			Metadata: make(map[string]any),
		}, nil
	})

	step4 := workflow.NewFunctionStep("finalize", func(ctx context.Context, input *workflow.StepInput) (*workflow.StepOutput, error) {
		fmt.Println("  ✅ Finalizing...")
		return &workflow.StepOutput{
			Content:  fmt.Sprintf("Final result: %v", input.PreviousStepContent),
			Metadata: make(map[string]any),
		}, nil
	})

	// 创建链式路由器 - 根据分析结果选择不同的处理链
	router := workflow.ChainRouter("smart_processor",
		func(input *workflow.StepInput) []string {
			// 根据前一步的分析结果决定执行路径
			if input.PreviousStepContent != nil {
				if analysis, ok := input.PreviousStepContent.(map[string]any); ok {
					if analysis["analysis"] == "complex" {
						fmt.Println("\n📍 Router 选择: complex 路径 (2步)")
						return []string{"process_complex", "finalize"}
					}
				}
			}
			fmt.Println("\n📍 Router 选择: simple 路径 (2步)")
			return []string{"process_simple", "finalize"}
		},
		map[string]workflow.Step{
			"process_complex": step2,
			"process_simple":  step3,
			"finalize":        step4,
		},
	)

	// 创建 Workflow
	wf := workflow.New("RouterDemo").
		WithStream().
		AddStep(step1).
		AddStep(router)

	if err := wf.Validate(); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
		return
	}

	fmt.Println("=== 开始流式执行 ===")

	// 执行并接收流式事件
	input := &workflow.WorkflowInput{
		Input: "Process this data",
	}

	eventCount := 0
	reader := wf.Execute(ctx, input)
	for {
		event, err := reader.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}

		eventCount++

		switch event.Type {
		case workflow.EventWorkflowStarted:
			fmt.Printf("\n[Event %d] 🚀 Workflow Started\n", eventCount)

		case workflow.EventStepStarted:
			fmt.Printf("\n[Event %d] ▶️  Step Started: %s\n", eventCount, event.StepName)

		case workflow.EventStepProgress:
			fmt.Printf("[Event %d] 📊 Step Progress: %s\n", eventCount, event.StepName)

		case workflow.EventStepCompleted:
			fmt.Printf("[Event %d] ✅ Step Completed: %s\n", eventCount, event.StepName)
			if data, ok := event.Data.(map[string]any); ok {
				if output, ok := data["output"].(*workflow.StepOutput); ok {
					fmt.Printf("   Output: %v\n", output.Content)
					if len(output.NestedSteps) > 0 {
						fmt.Printf("   Nested Steps: %d\n", len(output.NestedSteps))
						for i, nested := range output.NestedSteps {
							fmt.Printf("     %d. %s: %v\n", i+1, nested.StepName, nested.Content)
						}
					}
				}
			}

		case workflow.EventWorkflowCompleted:
			fmt.Printf("\n[Event %d] 🎉 Workflow Completed\n", eventCount)
			if data, ok := event.Data.(map[string]any); ok {
				if output, ok := data["output"]; ok {
					fmt.Printf("   Final Output: %v\n", output)
				}
				if metrics, ok := data["metrics"].(*workflow.RunMetrics); ok {
					fmt.Printf("   Total Time: %.3fs\n", metrics.TotalExecutionTime)
					fmt.Printf("   Steps: %d total, %d succeeded\n",
						metrics.TotalSteps, metrics.SuccessfulSteps)
				}
			}
		}
	}

	fmt.Printf("\n=== 完成 ===\n共处理 %d 个事件\n", eventCount)

	// 演示简单路由
	fmt.Println("\n\n=== 演示简单路由 ===")

	simpleRouter := workflow.SimpleRouter("simple_route",
		func(input *workflow.StepInput) string {
			// 简单的条件判断
			if inputStr, ok := input.Input.(string); ok {
				if len(inputStr) > 15 {
					return "route_complex"
				}
			}
			return "route_simple"
		},
		map[string]workflow.Step{
			"route_complex": step2,
			"route_simple":  step3,
		},
	)

	wf2 := workflow.New("SimpleRouterDemo").
		WithStream().
		AddStep(simpleRouter)

	input2 := &workflow.WorkflowInput{
		Input: "short",
	}

	fmt.Println("输入: 'short'")
	reader2 := wf2.Execute(ctx, input2)
	for {
		event, err := reader2.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			continue
		}
		if event.Type == workflow.EventWorkflowCompleted {
			if data, ok := event.Data.(map[string]any); ok {
				if output, ok := data["output"]; ok {
					fmt.Printf("结果: %v\n", output)
				}
			}
		}
	}
}
