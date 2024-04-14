package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	var cmd *exec.Cmd

	log.Println("MAAAAAAAAAAAAIIIIIIIIIIIIIINNNNNNNNNNNN")

	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "./script/swagger-codegen.ps1")
	} else {
		cmd = exec.Command("sh", "./script/swagger-codegen.sh")
	}
	// Attach STDOUT stream
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}

	// Attach STDIN stream
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println(err)
	}

	// Attach STDERR stream
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
	}

	// Spawn go-routine to copy os's stdin to command's stdin
	go io.Copy(stdin, os.Stdin)

	// Spawn go-routine to copy command's stdout to os's stdout
	go io.Copy(os.Stdout, stdout)

	// Spawn go-routine to copy command's stderr to os's stderr
	go io.Copy(os.Stderr, stderr)

	// Run() under the hood calls Start() and Wait()
	cmd.Run()

	// Note: The PIPES above will be closed automatically after Wait sees the command exit.
	// A caller need only call Close to force the pipe to close sooner.
	log.Println("Command complete")

}
