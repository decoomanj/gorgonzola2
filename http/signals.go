package gorgonzola

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"
)

var ShutdownChannel chan (int)
var TerminateChannel chan (int)

// Start signal handling
func init() {
	log.Println("Registering signal handler")
	ShutdownChannel = make(chan (int))
	TerminateChannel = make(chan (int))
	go handleSignals()
}

func handleSignals() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for sig := range c {
		switch sig {

		default:
			log.Panicf("unexpected signal: %v", sig)

		case syscall.SIGINT, syscall.SIGTERM:
			log.Printf("Shutdown started (signal %#v)\n", sig)
			close(ShutdownChannel)
			time.Sleep(8 * time.Second) // wait for a grace time TODO mark as "down"
			shutdown()

		case syscall.SIGQUIT:
			fmt.Printf("received signal %#v: printing stacktraces...\n", sig)
			stacktrace()

		}
	}
}

func shutdown() {
	r := recover()
	if r != nil {
		fmt.Printf("panic: %v\n", r)
		stacktrace()
	}
	flushLogs()
	fmt.Println("Goodbye!")
	close(TerminateChannel)
}

func stacktrace() {
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	pprof.Lookup("heap").WriteTo(os.Stdout, 1)
	pprof.Lookup("threadcreate").WriteTo(os.Stdout, 1)
	pprof.Lookup("block").WriteTo(os.Stdout, 1)
}
