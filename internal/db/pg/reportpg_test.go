package pg_test

import (
	"analytic-service/config"
	"analytic-service/internal/db"
	"analytic-service/internal/db/pg"
	"analytic-service/internal/logger"
	"analytic-service/internal/report"
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v4"
)

var ctx context.Context

func flush(ctx context.Context, task_id string, approved int, declined int, t *testing.T) {

	tx, err := pgreport.Pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			t.Log(err)
			return
		}
	}()

	batch := new(pgx.Batch)
	batch.Queue(`delete from task_progress_events t where t.task_id = $1;`, task_id)
	batch.Queue(`delete from aggr_task_livetime where task_id = $1;`, task_id)
	batch.Queue(`update aggr_state set cnt = $1 where state = $2;`, declined, report.DECLINED)
	batch.Queue(`update aggr_state set cnt = $1 where state = $2;`, approved, report.APPROVED)
	batch.Queue(`update aggr_state set cnt = 0 where state = $1;`, report.CREATED)
	res := tx.SendBatch(ctx, batch)

	err = res.Close()
	if err != nil {
		t.Fatal(err)
	}

	if err = tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
}

var pgreport *pg.ReportPG

// Инициализация среды тестирования.
func TestMain(m *testing.M) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 50*time.Minute)
	defer cancel()

	os.Setenv("DB_CONNECTION_STRING", "postgres://team21:mNgdxITbhVGd@91.185.93.23:5432/postgres")
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err.Error())
	}
	connstr := cfg.DbConnectionString
	if connstr == "" {
		log.Fatal("не указан адрес подключения к БД")
	}

	var ok bool
	pgreport, ok = db.NewReport(ctx, "PG", connstr, logger.Mock(nil)).(*pg.ReportPG)
	if !ok {
		log.Fatal("отчет не создался")
	}

	defer pgreport.Dispose()
	// Запуск тестов.
	os.Exit(m.Run())
}

func TestReport_PG_declined(t *testing.T) {
	taskId := uuid.NewString()
	approveDft, _ := pgreport.ApprovedCntGet(ctx)
	declineDft, _ := pgreport.DeclinedCntGet(ctx)
	defer flush(ctx, taskId, approveDft, declineDft, t)

	if err := pgreport.TaskStateSet(ctx, report.TaskState{TaskId: taskId, EventTime: time.Now(), State: report.CREATED}); err != nil {
		t.Error(err.Error(), nil)
	}
	approved, _ := pgreport.ApprovedCntGet(ctx)
	declined, _ := pgreport.DeclinedCntGet(ctx)
	if approved != approveDft || declined != declineDft {
		t.Errorf("%s:  %v != %v , %v != %v  ", report.CREATED, approved, approveDft, declined, declineDft)
	}

	if err := pgreport.TaskStateSet(ctx, report.TaskState{TaskId: taskId, EventTime: time.Now(), State: report.DECLINED}); err != nil {
		t.Error(err.Error(), nil)
	}
	approved, _ = pgreport.ApprovedCntGet(ctx)
	declined, _ = pgreport.DeclinedCntGet(ctx)
	if !(approved == approveDft && declined == 1+declineDft) {
		t.Errorf("%s:  %v != %v , %v != %v  ", report.DECLINED, approved, approveDft, declined, declineDft)
	}
}

func TestReport_PG_approved(t *testing.T) {
	//Получаем текущие значения полностью согласованных, отклоненных и время жизни задачи
	taskId := uuid.NewString()
	approveDft, _ := pgreport.ApprovedCntGet(ctx)
	declineDft, _ := pgreport.DeclinedCntGet(ctx)
	defer flush(ctx, taskId, approveDft, declineDft, t)
	if err := pgreport.TaskStateSet(ctx, report.TaskState{TaskId: taskId, EventTime: time.Now(), State: report.CREATED}); err != nil {
		t.Error(err.Error(), nil)
	}
	approved, _ := pgreport.ApprovedCntGet(ctx)
	declined, _ := pgreport.DeclinedCntGet(ctx)
	if approved != approveDft || declined != declineDft {
		t.Errorf("%s:  %v != %v , %v != %v  ", report.CREATED, approved, approveDft, declined, declineDft)
	}

	if err := pgreport.TaskStateSet(ctx, report.TaskState{TaskId: taskId, EventTime: time.Now(), State: report.APPROVED}); err != nil {
		t.Error(err.Error(), nil)
	}
	approved, _ = pgreport.ApprovedCntGet(ctx)
	declined, _ = pgreport.DeclinedCntGet(ctx)
	if !(approved == 1+approveDft && declined == declineDft) {
		t.Errorf("%s:  %v != %v , %v != %v  ", report.APPROVED, approved, approveDft, declined, declineDft)
	}
}

func TestReport_Livetime(t *testing.T) {
	taskId := uuid.NewString()
	_user_id := uuid.NewString()

	if err := pgreport.TaskStateSet(ctx, report.TaskState{TaskId: taskId, EventTime: time.Now(), State: report.CREATED}); err != nil {
		t.Error(err.Error(), nil)
	}

	if err := pgreport.NtfEventSet(ctx, report.TaskEvent{TaskId: taskId, EventTime: time.Now(), UserId: _user_id, State: report.MAILSEND}); err != nil {
		t.Error(err.Error(), nil)
	}
	//Проверим, что задача видна и живет 0 секунд
	tasklist, _ := pgreport.TaskListGet(ctx)
	if len(tasklist) == 0 {
		t.Errorf("tasklist is empty")
	}
	livetime := -1
	for _, element := range tasklist {
		if element.TaskId == taskId {
			livetime = element.LiveTime
			break
		}
	}
	if livetime != 0 {
		t.Errorf("bad task livetime")
	}

	time.Sleep(time.Second * 3)

	if err := pgreport.TaskApproveSet(ctx, report.TaskEvent{TaskId: taskId, EventTime: time.Now(), UserId: _user_id, State: report.CLICKAPPROVE}); err != nil {
		t.Error(err.Error(), nil)
	}
	livetime = 0
	list, _ := pgreport.TaskListGet(ctx)
	for _, element := range list {
		if element.TaskId == taskId {
			livetime += element.LiveTime
		}
	}
	if livetime < 3 {
		t.Errorf("bad task livetime %d<%d", livetime, 3)
	}
}
