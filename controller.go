package earlyreturn

import (
	"context"
	"time"

	"go.temporal.io/sdk/client"
)

type WorkflowCaller struct {
	client client.Client
}

func (w *WorkflowCaller) Call(tx Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	update := client.NewUpdateWithStartWorkflowOperation(client.UpdateWorkflowOptions{
		UpdateID:     "early-return",
		WaitForStage: client.WorkflowUpdateStageCompleted,
	})
	_, err := w.client.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		WithStartOperation: update,
		ID:                 "early-return",
		TaskQueue:          "default-task-queue",
	}, Workflow, tx)
	if err != nil {
		return err
	}

	fut, err := update.Get(ctx)
	if err != nil {
		return err
	}
	if err := fut.Get(ctx, nil); err != nil {
		return err
	}

	return nil
}
