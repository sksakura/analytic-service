package handlergrpc

import (
	"analytic-service/internal/report"
	"context"

	analytic "gitlab.com/g6834/team21/event-proto.git/generated/grpc"
)

type AnalyticService struct {
	Report report.Report
}

func (a AnalyticService) TaskState(
	ctx context.Context,
	request *analytic.TaskStateContract,
) (*analytic.Response, error) {
	err := a.Report.TaskStateSet(
		ctx,
		report.TaskState{
			TaskId:    request.TaskId,
			State:     report.TaskStateEnum(request.State.String()),
			EventTime: request.GetEventTime().AsTime()})
	return &analytic.Response{ResponseStatus: "OK"}, err
}

func (a AnalyticService) TaskEvent(
	ctx context.Context,
	request *analytic.TaskEventContract,
) (*analytic.Response, error) {
	err := a.Report.TaskApproveSet(ctx, report.TaskEvent{
		TaskId:    request.TaskId,
		UserId:    request.UserId,
		EventTime: request.GetEventTime().AsTime(),
		State:     report.TaskEventEnum(request.State.String()),
	})
	return &analytic.Response{ResponseStatus: "OK"}, err
}

func (a AnalyticService) NotificationEvent(
	ctx context.Context,
	request *analytic.NotificationContract,
) (*analytic.Response, error) {
	err := a.Report.NtfEventSet(ctx, report.TaskEvent{
		TaskId:    request.UniqueKey.TaskId,
		UserId:    request.UniqueKey.UserId,
		EventTime: request.NotificationTime.AsTime(),
		State:     report.TaskEventEnum(request.UniqueKey.State.String()),
	})
	return &analytic.Response{ResponseStatus: "OK"}, err
}
