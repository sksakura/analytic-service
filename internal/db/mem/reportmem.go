package mem

import (
	"analytic-service/internal/report"
	"context"
)

type ReportMEM struct {
}

func (r ReportMEM) Dispose() {
}

func (r ReportMEM) TaskListGet(ctx context.Context) ([]*report.TaskReport, error) {
	content := []*report.TaskReport{}
	return content, nil
}
func (r ReportMEM) ApprovedCntGet(ctx context.Context) (int, error) {
	return 111, nil
}
func (r ReportMEM) DeclinedCntGet(ctx context.Context) (int, error) {
	return 2222, nil
}

func (r ReportMEM) TaskStateSet(ctx context.Context, task report.TaskState) error {
	return nil
}

func (r ReportMEM) NtfEventSet(ctx context.Context, notification report.TaskEvent) error {
	return nil
}

func (r ReportMEM) TaskApproveSet(ctx context.Context, event report.TaskEvent) error {
	return nil
}
