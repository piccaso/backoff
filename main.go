package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/piccaso/backoff/command"
	"os"
	"time"
)

func main() {
	flagOut := os.Stderr
	flag.CommandLine.SetOutput(flagOut)
	duration := flag.Int("d", 5, "initial duration [sec]")
	increment := flag.Int("i", 5, "increment [sec]")
	reset := flag.Float64("s", 5, "reset timer after [sec]")
	max := flag.Int("m", 100, "max count")

	flag.Usage = func() {
		fmt.Fprintf(flagOut, "Usage:\n\tbackoff [options] -- <command> [args...]\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()
	//fmt.Fprintf(flagOut, "eol:%vsec\n", *duration+((*increment)*(*max)))

	if len(args) < 1 {
		flag.Usage()
		os.Exit(2)
	}

	sleepTime := *duration
	failCnt := 0
	for {
		cmd, ctx, err := command.NewCommandWithContext(context.Background(), args)
		if err != nil {
			panic(err)
		}

		func() {
			cancel, _ := context.WithCancel(ctx)
			defer cancel.Done()

			start := time.Now()
			err := cmd.Run()
			elapsed := time.Since(start).Seconds()
			if err == nil {
				fmt.Fprintf(flagOut, "no error, runtime %v sec\n", elapsed)
				os.Exit(0)
			}

			if elapsed > *reset {
				sleepTime = *duration
				failCnt = 0

				fmt.Fprint(flagOut, "reset\n")
			} else {

				if failCnt > 0 {
					sleepTime += *increment
				}
				failCnt++

				fmt.Fprintf(flagOut, "failed, incremented sleepTime: %v, failCnt: %v\n", sleepTime, failCnt)
			}
			if failCnt > *max {
				fmt.Fprintf(flagOut, "max retries reached: %v\n", failCnt)
				os.Exit(2)
			}
		}()

		fmt.Fprintf(flagOut, "backing off for %v seconds\n", sleepTime)
		time.Sleep(time.Second * time.Duration(sleepTime))
	}

}
