package main

import (
	"fmt"
	"os"
	"log"
	"os/exec"
	"syscall"

	"github.com/docker/docker/pkg/reexec"
)


func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("run func init()")
	reexec.Register("nsInitialisation", nsInitialisation)
	log.Println("finish reexec.Register()")
	if reexec.Init() {
		log.Println("reexec.init() have been init()")
		os.Exit(0)
	}
	log.Println("run func init() finish")
}

func nsInitialisation() {
	log.Println(">> namespace setup code goes here <<")
	nsRun()
}

func nsRun() {
	cmd := exec.Command("/bin/sh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"PS1=-[ns-process]- # "}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the /bin/sh command - %s\n", err)
		os.Exit(1)
	}
}

func main() {
	log.Println("main() begin in first line")
	cmd := reexec.Command("nsInitialisation")
	log.Println("main() construct  reexec.Command()")
	log.Println(cmd.Path)
	log.Println(cmd.Args[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the reexec.Command - %s\n", err)
		os.Exit(1)
	}
}
