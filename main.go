package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mitchellh/go-ps"
)

var (
	showVersion = flag.Bool("v", false, "show version")
)

var version = "1.0.0"

type process struct {
	Pid int
	Cmd string
}

func processes() ([]process, error) {
	var processes []process
	procs, err := ps.Processes()
	if err != nil {
		return nil, err
	}

	for _, p := range procs {
		// skip pid 0
		if p.Pid() == 0 {
			continue
		}
		processes = append(processes, process{
			Pid: p.Pid(),
			Cmd: p.Executable(),
		})
	}

	return processes, nil
}

func kill(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return p.Kill()
}

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Println("fk - " + version)
		return
	}

	procs, err := processes()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	idx, err := fuzzyfinder.FindMulti(
		procs,
		func(i int) string {
			return fmt.Sprintf("%d: %s", procs[i].Pid, procs[i].Cmd)
		},
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, i := range idx {
		pid, cmd := procs[i].Pid, procs[i].Cmd
		fmt.Println(pid, cmd)

		if err := kill(pid); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
