package main

import (
	"context"
	"log"
	"server-template/config"
	"server-template/internal/server"
	"server-template/pkg/logger"

	authGRPC "server-template/internal/auth/handler/delivery/grpc"
	"server-template/internal/auth/repository"
	UC "server-template/internal/auth/usecase"

	"os"
	"os/signal"
	"server-template/pkg/storage/redis"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	zerolog "github.com/philip-bui/grpc-zerolog"
	redisSource "github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"google.golang.org/grpc"
)

func main() {
	cfgFile, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig: %s", err.Error())
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %s", err.Error())
	}
	log.Println("Config loaded")

	config.C = cfg
	appLogger := logger.NewAPILogger(cfg)

	err = appLogger.InitLogger()
	if err != nil {
		log.Fatalf("Cannot init logger: %s", err.Error())
	}
	appLogger.Infof("Logger successfully started with level: %s, InFile: %t (filePath: %s), InTG: %t (chatID: %d)",
		cfg.Logger.Level,
		cfg.Logger.InFile,
		cfg.Logger.FilePath,
		cfg.Logger.InTG,
		cfg.Logger.ChatID,
	)

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.OpenTelemetry.Host)))
	if err != nil {
		log.Fatalf("Cannot create Jaeger exporter - %s", err.Error())
	}

	defer func() {
		err = exp.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}()

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.OpenTelemetry.ServiceName),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	defer func() {
		err = tp.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("OpenTelemetry closed properly")
		}
	}()

	ctx := context.Background()

	redisClient, err := redis.NewRedisClient(cfg)
	if err != nil {
		appLogger.Fatalf("Redis init error: %s", err.Error())
	} else {
		appLogger.Infof("Redis connected, status: %#v", redisClient.PoolStats())
	}
	defer func(redisClient *redisSource.Client) {
		err = redisClient.Close()
		if err != nil {
			appLogger.Info(err.Error())
		} else {
			appLogger.Info("Redis closed properly")
		}
	}(redisClient)

	redisRepo := repository.NewAuthRedisRepo(redisClient, appLogger, cfg)
	usecase := UC.NewAuthUC(
		redisRepo,
		appLogger,
		cfg,
	)
	authGRPCHandler := authGRPC.NewAuthHandlers(usecase)

	deps := server.Deps{
		AuthHandlers: authGRPCHandler,
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				otelgrpc.UnaryServerInterceptor(),
				zerolog.NewUnaryServerInterceptorWithLogger(appLogger.GetLogger()),
			),
		),
	)
	grpcServer := server.NewGRPCServer(srv, deps, cfg, appLogger)
	go func() {
		appLogger.Infof("GRPC server listening %s:%s", cfg.Server.Host, cfg.Server.GRPCPort)
		if err = grpcServer.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	httpServer := server.NewServer(cfg, redisClient, appLogger)
	go func() {
		appLogger.Infof("HTTP Server listening %s:%s", cfg.Server.Host, cfg.Server.HTTPPort)
		if err = httpServer.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	grpcServer.GracefulShutdown()
	if err = httpServer.Shutdown(); err != nil {
		appLogger.Error(err)
	}
}
