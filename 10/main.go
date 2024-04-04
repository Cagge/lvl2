package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	defaultTimeout = 10 * time.Second
)

type Options struct {
	Host    string
	Port    string
	Timeout time.Duration
}

func parseArgs(args []string) (Options, error) {
	options := Options{}

	flagSet := flag.NewFlagSet("", flag.ExitOnError)
	flagSet.DurationVar(&options.Timeout, "timeout", defaultTimeout, "timeout")

	err := flagSet.Parse(args)
	if err != nil {
		return options, fmt.Errorf("failed to parse flags: %s", err)
	}

	options.Host = flagSet.Arg(0)
	options.Port = flagSet.Arg(1)

	return options, nil
}

func Run(args []string) error {
	options, err := parseArgs(args)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(options.Host, options.Port), options.Timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return HandleInput(ctx, conn)
	})

	group.Go(func() error {
		return HandleOutput(ctx, conn)
	})

	go func() {
		err = group.Wait()
		if err != nil && err != io.EOF {
			log.Println(err)
			return
		}
	}()

	<-ctx.Done()

	return nil
}

func HandleInput(ctx context.Context, conn net.Conn) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line, _, err := reader.ReadLine()
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(conn, string(line))
			if err != nil {
				return err
			}
		}
	}
}

func HandleOutput(ctx context.Context, conn net.Conn) error {
	reader := bufio.NewReader(conn)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line, _, err := reader.ReadLine()
			if err != nil {
				return err
			}

			fmt.Println(string(line))
		}
	}
}

func main() {
	err := Run(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
