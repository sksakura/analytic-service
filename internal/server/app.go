package server

import (
	"analytic-service/config"
	"analytic-service/internal/db"
	handlergrpc "analytic-service/internal/handler/grpc"
	handlerhttp "analytic-service/internal/handler/http"
	"analytic-service/internal/handler/kafka"
	"analytic-service/internal/logger"
	"analytic-service/internal/profiler"
	"context"
	"net"
	"net/http"

	analytic "gitlab.com/g6834/team21/event-proto.git/generated/grpc"

	"github.com/go-chi/chi"
	l "github.com/treastech/logger"
	"gitlab.com/g6834/team21/authproto.git/pkg/tokensservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Application struct {
	handlerHttp     *handlerhttp.HandlerHttp
	Logger          *logger.Logger
	authServiceConn *grpc.ClientConn
	grpcserver      *grpc.Server
	handlerKafka    *kafka.KafkaClient
	cfg             *config.Config
}

func (a *Application) Dispose() {
	a.handlerHttp.Dispose()
	a.authServiceConn.Close()
	a.handlerKafka.Dispose()
}

func New(cfg *config.Config, logger *logger.Logger, grpcDialOption grpc.DialOption) *Application {
	ctx := context.Background()

	if cfg.DbConnectionString == "" {
		logger.Fatal("не указан адрес подключения к БД")
	}
	report := db.NewReport(ctx, cfg.DbType, cfg.DbConnectionString, logger)
	server := grpc.NewServer()
	reflection.Register(server)
	analytic.RegisterAnalyticServiceServer(server, &handlergrpc.AnalyticService{Report: report})

	logger.Info("приготовили grpc для аналитики")

	authadress := cfg.AuthAddress
	if authadress == "" {
		logger.Fatal("не указан адрес подключения к Аутентификации")
	}

	authConnection, err := grpc.DialContext(ctx, cfg.AuthAddress, grpcDialOption, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(err.Error())
	}
	ts := tokensservice.NewTokensServiceClient(authConnection)

	logger.Info("подключились к сервису авторизации")
	kafka, err := kafka.NewKafka(
		cfg.Brokers,
		cfg.GroupId,
		logger,
		&kafka.HandlerKafka{Report: report},
		kafka.TopicSet{State: cfg.StateKafkaTopic, Event: cfg.EventKafkaTopic, Notify: cfg.NotifyKafkaTopic})
	if err != nil {
		logger.Fatal("не указан адрес подключения к kafka")
	}

	obj := Application{
		Logger:       logger,
		grpcserver:   server,
		cfg:          cfg,
		handlerKafka: kafka,
		handlerHttp: handlerhttp.New(
			report,
			profiler.NewProfiler(cfg.ProfilerStatus == "true"),
			logger,
			&handlerhttp.Auth{
				Tokensservice: ts,
			}),
	}
	return &obj
}

func (a *Application) Run() chan error {
	r := chi.NewRouter()
	r.Use(a.handlerHttp.AuthMiddleware)
	r.Use(l.Logger(a.Logger.Logger))

	r.Mount("/", Common(a.handlerHttp))
	r.With(a.handlerHttp.CheckProfiler).Mount("/debug", Profiler())

	a.Logger.Info("подготовили роутер")

	errCh := make(chan error)

	a.handlerKafka.StartReaders(context.Background())

	listener, err := net.Listen("tcp", a.cfg.GrpcAddress)
	if err != nil {
		a.Logger.Fatal(err.Error())
	}
	a.Logger.Info("создали grpc listener")

	//Залогиниться

	go func() {
		if err := a.grpcserver.Serve(listener); err != nil {
			errCh <- err
		}
	}()

	go func() {
		err := http.ListenAndServe(":3000", r)
		if err != nil {
			errCh <- err
		}
	}()
	return errCh
}
