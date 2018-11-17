package mainutil

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type CancelableFunc func(ctx context.Context) error

func TrapAndWait(ctx context.Context, cancelables ...CancelableFunc) int {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigch := make(chan os.Signal, 1)
	errch := make(chan error)

	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	for _, c := range cancelables {
		go func(c CancelableFunc) {
			errch <- c(ctx)
		}(c)
	}

	var rv, rets int
	for {
		select {
		case sig := <-sigch:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt:
				cancel()
			}
		case err := <-errch:
			cancel()

			if err != nil {
				log(ctx).Error("process ended with error", "error", err)
				rv = 1
			}

			rets++
			if rets == len(cancelables) {
				return rv
			}
		}
	}
}
