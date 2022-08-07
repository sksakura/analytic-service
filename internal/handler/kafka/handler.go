package kafka

import (
	"analytic-service/internal/report"
	"context"

	analytic "gitlab.com/g6834/team21/event-proto.git/generated/grpc"

	"google.golang.org/protobuf/proto"
)

type HandlerKafka struct {
	Report report.Report
}

func (a HandlerKafka) ProcessState(ctx context.Context, message Message) error {
	request := analytic.TaskStateContract{}
	if err := proto.Unmarshal(message.Value, &request); err != nil {
		return err
	}
	return a.Report.TaskStateSet(
		ctx,
		report.TaskState{
			TaskId:    request.TaskId,
			State:     report.TaskStateEnum(request.State.String()),
			EventTime: request.GetEventTime().AsTime()})
}

func (a HandlerKafka) ProcessEvent(ctx context.Context, message Message) error {
	request := analytic.TaskEventContract{}
	if err := proto.Unmarshal(message.Value, &request); err != nil {
		return err
	}
	return a.Report.TaskApproveSet(ctx, report.TaskEvent{
		TaskId:    request.TaskId,
		UserId:    request.UserId,
		EventTime: request.GetEventTime().AsTime(),
		State:     report.TaskEventEnum(request.State.String()),
	})
}

func (a HandlerKafka) ProcessNotify(ctx context.Context, message Message) error {
	request := analytic.NotificationContract{}
	if err := proto.Unmarshal(message.Value, &request); err != nil {
		return err
	}
	return a.Report.NtfEventSet(ctx, report.TaskEvent{
		TaskId:    request.UniqueKey.TaskId,
		UserId:    request.UniqueKey.UserId,
		EventTime: request.NotificationTime.AsTime(),
		State:     report.TaskEventEnum(request.UniqueKey.State.String()),
	})
}
