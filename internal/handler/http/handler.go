package handlerhttp

import (
	"analytic-service/internal/logger"
	"analytic-service/internal/profiler"
	"analytic-service/internal/report"

	"gitlab.com/g6834/team21/authproto.git/pkg/tokensservice"
)

type Auth struct {
	Tokensservice tokensservice.TokensServiceClient
}

type HandlerHttp struct {
	report   report.Report
	profiler *profiler.Profiler
	logger   *logger.Logger
	auth     *Auth
}

func (h *HandlerHttp) Dispose() {
	h.report.Dispose()
}

func New(report report.Report, profiler *profiler.Profiler, logger *logger.Logger, auth *Auth) *HandlerHttp {
	obj := HandlerHttp{
		report:   report,
		logger:   logger,
		auth:     auth,
		profiler: profiler,
	}
	return &obj
}
