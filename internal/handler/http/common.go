package handlerhttp

import (
	"fmt"
	"net/http"
)

// SayHello godoc
// @Summary SayHello
// @Description Say hello to auth and to no-auth user
// @Tags         common
// @Success      200
// @Router       / [get]
// @Produce text/plain
func (h *HandlerHttp) SayHello(w http.ResponseWriter, r *http.Request) {
	helloString := "welcome"
	prof, ok := r.Context().Value(ctxKey{}).(Profile)
	if ok {
		helloString += ", " + prof.UserId
	}
	_, err := w.Write([]byte(helloString))
	if err != nil {
		h.logger.Error(err.Error())
	}
}

// PprofOn godoc
// @Summary Включение профилировщика
// @Description апи включает доступ к роутам профилирования
// @Tags         pprof
// @Success 200
// @Router /pprof_state/on [get]
// @Produce text/plain
func (h *HandlerHttp) PprofOn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.profiler.On()
		next.ServeHTTP(w, r)
	})
}

// PprofOff godoc
// @Summary Выключение профилировщика
// @Description апи выключает доступ к роутам профилирования
// @Tags         pprof
// @Success 200
// @Router /pprof_state/off [get]
// @Produce text/plain
func (h *HandlerHttp) PprofOff(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.profiler.Off()
		next.ServeHTTP(w, r)
	})
}

// GetPprof godoc
// @Summary Статус профилировщика
// @Description Получение текущего состояния профилировщика (включен - выключен)
// @Tags         pprof
// @Success 200
// @Router /pprof_state/ [get]
// @Produce text/plain
func (h *HandlerHttp) GetPprof(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(fmt.Sprintf("pprof_on state:%v", h.profiler.Status())))
	if err != nil {
		h.logger.Error(err.Error())
	}
}
