package pg

import (
	"analytic-service/internal/logger"
	"analytic-service/internal/report"
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type ReportPG struct {
	Pool   *pgxpool.Pool
	Logger *logger.Logger
}

func (r ReportPG) Dispose() {
	r.Pool.Close()
}

func (r ReportPG) TaskListGet(ctx context.Context) ([]*report.TaskReport, error) {

	query := `select atl.task_id , atl.livetime  from public.aggr_task_livetime atl;`
	data, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	result := make([]*report.TaskReport, 0, 5)
	for data.Next() {
		var task report.TaskReport
		err = data.Scan(
			&task.TaskId,
			&task.LiveTime,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &task)
	}
	return result, nil
}
func (r ReportPG) ApprovedCntGet(ctx context.Context) (int, error) {

	query := `select cnt from aggr_state as2 where as2.state = $1;`
	data, err := r.Pool.Query(ctx, query, report.APPROVED)
	if err != nil {
		return -1, err
	}
	var cnt int = -1
	for data.Next() {
		if err := data.Scan(&cnt); err != nil {
			return -1, err
		}
	}
	return cnt, nil
}

func (r ReportPG) DeclinedCntGet(ctx context.Context) (int, error) {
	query := `select cnt from aggr_state as2 where as2.state = $1;`
	data, err := r.Pool.Query(ctx, query, report.DECLINED)
	if err != nil {
		return -1, err
	}
	var cnt int = -1
	for data.Next() {
		if err := data.Scan(&cnt); err != nil {
			return -1, err
		}
	}
	return cnt, nil
}

func (r ReportPG) TaskStateSet(ctx context.Context, task report.TaskState) error {
	//Если в бд есть более свежие записи - предполагаем, что пришел повтор и выходим
	var cnt int
	err := r.Pool.QueryRow(
		ctx,
		"select count(t.id)  from task_state_events t  where t.task_id = $1 and t.event_time >= $2", task.TaskId, task.EventTime).Scan(&cnt)
	if err != nil {
		return err
	}
	if cnt > 0 {
		r.Logger.Info("В базе присутствуют более свежии записи")
		return nil
	}

	if task.State == report.DELETED {

		tx, err := r.Pool.Begin(ctx)
		if err != nil {
			return err
		}

		defer func() {
			err := tx.Rollback(ctx)
			if err != nil {
				r.Logger.Error(err.Error())
			}
		}()

		batch := new(pgx.Batch)
		batch.Queue(`delete from aggr_task_livetime where task_id = $1;`, task.TaskId)
		res := tx.SendBatch(ctx, batch)

		err = res.Close()
		if err != nil {
			return err
		}

		if err = tx.Commit(ctx); err != nil {
			return err
		}
		return nil
	}

	_, err = r.Pool.Exec(ctx, "CALL sp_set_task_state($1, $2, $3);", task.EventTime, task.State, task.TaskId)
	return err
}

//Фиксируем отправку нотификации по задаче
func (r ReportPG) NtfEventSet(ctx context.Context, notification report.TaskEvent) error {
	var cnt int
	err := r.Pool.QueryRow(ctx, "select count(t.id)  from task_progress_events t  where t.task_id = $1 and t.event_time >= $2", notification.TaskId, notification.EventTime).Scan(&cnt)
	if err != nil {
		return err
	}
	if cnt > 0 {
		r.Logger.Info(fmt.Sprintf("В базе присутствуют более свежие записи task_id: %s, user_id: %s", notification.TaskId, notification.UserId))
		return nil
	}
	_, err = r.Pool.Exec(ctx, "CALL sp_set_task_progress($1,$2,$3,$4);", notification.EventTime, notification.State, notification.TaskId, notification.UserId)
	return err
}

//Фиксируем прогресс по задаче
func (r ReportPG) TaskApproveSet(ctx context.Context, event report.TaskEvent) error {
	//Если в бд есть более свежие записи - предполагаем, что пришел повтор и выходим
	var cnt int
	err := r.Pool.QueryRow(ctx, "select count(t.id)  from task_progress_events t  where t.task_id = $1 and t.event_time >= $2", event.TaskId, event.EventTime).Scan(&cnt)
	if err != nil {
		return err
	}
	if cnt > 0 {
		r.Logger.Info("В базе присутствуют более свежие записи", zap.String("task_id", event.TaskId), zap.String("user_id", event.UserId))
		return nil
	}
	_, err = r.Pool.Exec(
		ctx, "CALL sp_set_task_progress($1,$2,$3,$4);", event.EventTime, event.State, event.TaskId, event.UserId)
	return err
}
