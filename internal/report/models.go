package report

import (
	"context"
	"time"
)

type (
	TaskStateEnum string
	TaskEventEnum string
)

const (
	APPROVED TaskStateEnum = "APPROVED"
	DECLINED TaskStateEnum = "DECLINED"
	CREATED  TaskStateEnum = "CREATED"
	DELETED  TaskStateEnum = "DELETED"

	MAILSEND     TaskEventEnum = "MAIL_SEND"
	CLICKAPPROVE TaskEventEnum = "CLICK_APPROVE"
)

type TaskReport struct {
	TaskId   string
	LiveTime int
}

type Report interface {
	Dispose()
	TaskListGet(ctx context.Context) ([]*TaskReport, error)
	ApprovedCntGet(ctx context.Context) (int, error)
	DeclinedCntGet(ctx context.Context) (int, error)
	TaskStateSet(ctx context.Context, taskState TaskState) error
	NtfEventSet(ctx context.Context, notification TaskEvent) error
	TaskApproveSet(ctx context.Context, taskEvent TaskEvent) error
}

//состояния всей задачи
type TaskState struct {
	TaskId    string
	State     TaskStateEnum
	EventTime time.Time
}

//эвент - информация о том, что письмо по задаче ушло юзеру, информация о том, что по задаче пользователь кликнул согласование
type TaskEvent struct {
	TaskId    string
	EventTime time.Time
	UserId    string
	State     TaskEventEnum
}

type NotiUniqueKeyUniqueKey struct {
	ServiceId string        //тут всегда TaskService
	TaskId    string        //uuid
	UserId    string        //uuid
	State     TaskEventEnum //MAIL_SEND нужно отсюда будет обрабаывать
}

type NotificationContract struct {
	RecipientList    []string
	Subject          string
	Body             string
	NotificationTime time.Time
	UniqueKey        NotiUniqueKeyUniqueKey
}
