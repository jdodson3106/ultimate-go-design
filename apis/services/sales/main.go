package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ardanlabs/service/foundation/logger"
)

// build represents the version of the code was used to build the binary.
// this can be linked to a git tag on a release and set through CI/CD
// using the -ldflag "-X main.build=${<env-name>}" on the go build command.
// see dockerfile.sales for a working example of this.
var build = "develop"

func main() {
	var log *logger.Logger
	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT HERE *******")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return ""
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SALES", traceIDFn, events)

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "build", build)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	sig := <-shutdown

	log.Info(ctx, "shutdown", "status", "shutdown initiated", "signal", sig)
	defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

	return nil
}
