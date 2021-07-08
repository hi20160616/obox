package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hi20160616/obox/internal/server"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Valid
	if err := server.ValidPasswd(); err != nil {
		fmt.Println(err)
	}
	// randNum := func() string {
	//         rand.Seed(time.Now().UnixNano())
	//         return strconv.Itoa(rand.Intn(99999))
	// }
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	// Web server
	address := ":80"
	s, err := server.NewServer(address)
	if err != nil {
		log.Printf("%v", err)
	}
	g.Go(func() error {
		log.Printf("Server start on %s", address)
		return s.Start(ctx)
	})
	g.Go(func() error {
		defer log.Printf("Server stop done.")
		<-ctx.Done() // wait for stop signal
		log.Printf("Server stop now...")
		return s.Stop(ctx)
	})

	// Elegant stop
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	g.Go(func() error {
		select {
		case sig := <-sigs:
			fmt.Println()
			log.Printf("signal caught: %s, ready to quit...", sig.String())
			cancel()
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	})

	// Err treat
	if err := g.Wait(); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("not canceled by context: %s", err)
		} else {
			log.Println(err)
		}
	}
}
