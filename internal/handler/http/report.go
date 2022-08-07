package handlerhttp

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Report struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Error struct {
	Message string `json:"content"`
}

func Model2View(tp string, content string) interface{} {
	return Report{Type: tp, Content: content}
}

// ReportTaskList godoc
// @Summary Время реакций по задачам
// @Description Суммарное время ожидания реакций по каждой задаче - сумма всех времен от создания задачи до нажатия на ссылку каждого согласующего
// @Tags         report
// @Success 200 {object} Report
// @Failure 401 {object} Error "не найден access token; авторизация не удалась"
// @Router /tasklist [get]
// @Produce json
func (h *HandlerHttp) ReportTaskList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, ok := r.Context().Value(ctxKey{}).(Profile)
	if !ok {
		err := respondwithJSON(w, http.StatusUnauthorized, &Error{Message: "не найден access token; авторизация не удалась"})
		if err != nil {
			h.logger.Error(err.Error())
		}
		return
	}
	result, err := h.report.TaskListGet(ctx)
	if err != nil {
		h.logger.Error(err.Error())
	}

	jsonString, err := json.Marshal(result)
	if err != nil {
		h.logger.Error(err.Error())
		respondwithJSON(w, http.StatusInternalServerError, err)
		return
	}

	if err := respondwithJSON(w, http.StatusOK, Model2View("livetime", string(jsonString))); err != nil {
		h.logger.Error(err.Error())
	}
}

// ReportApproved godoc
// @Summary Согласованные задачи
// @Description Количество полностью согласованных задач
// @Tags         report
// @Success 200 {object} Report
// @Failure 401 {object} Error "не найден access token; авторизация не удалась"
// @Router /approved [get]
// @Produce json
func (h *HandlerHttp) ReportApproved(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, ok := ctx.Value(ctxKey{}).(Profile)
	if !ok {
		err := respondwithJSON(w, http.StatusUnauthorized, &Error{Message: "не найден access token; авторизация не удалась"})
		if err != nil {
			h.logger.Error(err.Error())
		}
		return
	}
	result, err := h.report.ApprovedCntGet(ctx)
	if err != nil {
		h.logger.Error(err.Error())
		respondwithJSON(w, http.StatusInternalServerError, err)
		return
	}

	if err := respondwithJSON(w, http.StatusOK, Model2View("approved", strconv.Itoa(result))); err != nil {
		h.logger.Error(err.Error())
	}
}

// ReportDeclined godoc
// @Summary Отклоненные задачи
// @Description Количество несогласованных задач
// @Tags         report
// @Success 200 {object} Report
// @Failure 401 {object} Error "не найден access token; авторизация не удалась"
// @Router /declined [get]
// @Produce json
func (h *HandlerHttp) ReportDeclined(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, ok := ctx.Value(ctxKey{}).(Profile)
	if !ok {
		err := respondwithJSON(w, http.StatusUnauthorized, &Error{Message: "не найден access token; авторизация не удалась"})
		if err != nil {
			h.logger.Error(err.Error())
		}
		return
	}

	result, err := h.report.DeclinedCntGet(ctx)
	if err != nil {
		h.logger.Error(err.Error())
		respondwithJSON(w, http.StatusInternalServerError, err)
		return
	}

	if err := respondwithJSON(w, http.StatusOK, Model2View("declined", strconv.Itoa(result))); err != nil {
		h.logger.Error(err.Error())
	}
}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	return err
}
