package earlyreturn

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/mocks"
	"go.temporal.io/sdk/testsuite"
)

func TestWorkflow(t *testing.T) {
	ts := &testsuite.WorkflowTestSuite{}
	env := ts.NewTestWorkflowEnvironment()

	tx := Transaction{
		ID:            "123",
		SourceAccount: "1234",
		TargetAccount: "5678",
		Amount:        1000,
	}

	env.RegisterWorkflow(Workflow)
	env.RegisterActivity(tx.InitTransaction)
	env.RegisterActivity(tx.CompleteTransaction)

	c := mocks.NewClient(t)
	caller := &WorkflowCaller{client: c}
	c.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil).
		Run(func(args mock.Arguments) {
			wf := args.Get(2)
			tx := args.Get(3)
			env.ExecuteWorkflow(wf, tx)
		})

	err := caller.Call(tx)
	require.NoError(t, err)
}
