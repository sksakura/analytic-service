package server

import (
	handlerhttp "analytic-service/internal/handler/http"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "analytic-service/docs"
)

func Profiler() http.Handler {
	rProfiler := chi.NewRouter()

	rProfiler.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/pprof/", http.StatusMovedPermanently)
	})
	rProfiler.HandleFunc("/pprof", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
	})

	// Получение списка всех профилей
	rProfiler.HandleFunc("/pprof/*", pprof.Index)
	// Отображение строки запуска (например: /go-observability-course/examples/caching/redis/__debug_bin)
	rProfiler.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	// профиль ЦПУ, в query-параметрах можно указать seconds со значением времени в секундах для снимка (по-умолчанию 30с)
	rProfiler.HandleFunc("/pprof/profile", pprof.Profile)
	rProfiler.HandleFunc("/pprof/symbol", pprof.Symbol)
	// профиль для получения трассировки (последовательности инструкций) выполнения приложения за время seconds из query-параметров ( по-умолчанию 1с)
	rProfiler.HandleFunc("/pprof/trace", pprof.Trace)
	return rProfiler
}

func Common(a *handlerhttp.HandlerHttp) http.Handler {
	r := chi.NewRouter()
	r.Get("/", a.SayHello)
	r.Get("/tasklist", a.ReportTaskList)
	r.Get("/approved", a.ReportApproved)
	r.Get("/declined", a.ReportDeclined)

	r.Route("/pprof_state", func(r chi.Router) {
		r.Get("/", a.GetPprof)
		r.With(a.PprofOn).Get("/on", a.GetPprof)
		r.With(a.PprofOff).Get("/off", a.GetPprof)
	})
	// установка маршрута для документации
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	return r
}
