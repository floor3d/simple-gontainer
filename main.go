/*

*** Below documentation made by chatgpt :) 

Package main is a simple code to create a basic container in Go.

This code demonstrates how to create a container using Go by using Linux namespaces. It includes two functions, `parent` and `child`, which are used to run the container processes.

The `main` function is the entry point of the program and determines whether to run the parent or child process based on the command-line arguments.

The `parent` function is responsible for setting up the container environment and executing the child process inside the container. It uses the `os/exec` package to execute the child process with the necessary configuration.

The `child` function is the actual container process. It sets up the container's root filesystem, mounts it, changes the root directory to the new filesystem, and executes the specified command inside the container.

The `must` function is a helper function that panics if an error occurs. It is used to handle errors in a simple and straightforward manner.

Usage:
	go run main.go run [command] [arguments]

The command-line argument "run" is used to indicate that the container should be executed. The subsequent arguments specify the command to run inside the container.

Example:
	go run main.go run /bin/bash

This will run a container and execute the "/bin/bash" command inside it.

Note:
This code assumes a Linux-based operating system, as it relies on Linux namespaces and system calls specific to Linux.

References:
- Linux Namespaces: https://man7.org/linux/man-pages/man7/namespaces.7.html
- `os/exec` package: https://golang.org/pkg/os/exec/
- `syscall` package: https://golang.org/pkg/syscall/
*/
package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
        panic("Wrong argument supplied: need 'run' or 'child'")
	}
}

// parent sets up the container environment and executes the child process inside the container.
func parent() {
	// Create a command to execute the current program (self) with arguments for the child process.
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	// Configure the command to use Linux namespaces for isolation.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	// Set the standard input, output, and error streams for the child process.
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and handle any errors.
	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

// child is the actual container process.
func child() {
	// Mount the root filesystem to itself as a bind mount.
	must(syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))

	// Create a directory to serve as the old root for the pivot operation.
	must(os.MkdirAll("rootfs/oldrootfs", 0700))

	// Change the root directory to the new filesystem using the pivot_root system call.
	must(syscall.PivotRoot("rootfs", "rootfs/oldrootfs"))

	// Change the current working directory to the new root.
	must(os.Chdir("/"))

	// Create a command to execute the specified command inside the container.
	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	// Set the standard input, output, and error streams for the command.
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and handle any errors.
	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

// Panic if error occurs
func must(err error) {
    // TODO: issue where it cant find file; debug
	if err != nil {
		panic(err)
	}
}
