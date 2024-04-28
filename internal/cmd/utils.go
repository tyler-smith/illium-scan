package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-cz/devslog"
)

func SetLogger() {
	slog.SetDefault(newDevLogger())
}

func newDevLogger() *slog.Logger {
	return slog.New(devslog.NewHandler(os.Stdout, &devslog.Options{
		HandlerOptions: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}}))
}

func WaitForExit() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop()
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	//for sig := range c {
	//	if sig == syscall.SIGINT || sig == syscall.SIGTERM {
	//		break
	//	}
	//}
}
