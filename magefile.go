//go:build mage
// +build mage

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	// mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	//mg.Deps(InstallDeps)
	fmt.Println("Prepping...")
	err := Prep()
	if err != nil {
		return err
	}

	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", "cast-via-mage", ".")
	return cmd.Run()
}

func Prep() error {
	//gofmt -s -w .
	fmt.Println("Formatting, Vetting, Etc...")
	cmd := exec.Command("gofmt", "-s", "-w", ".")

	//go vet ./...
	cmd = exec.Command("go", "vet", "./...")

	//go mod tidy
	cmd = exec.Command("go", "mod", "tidy")

	//go mod download
	cmd = exec.Command("go", "mod", "download")

	//go mod verify
	cmd = exec.Command("go", "mod", "verify")

	return cmd.Run()
}

// A custom install step if you need your bin someplace other than go/bin
//func Install() error {
//mg.Deps(Build)
//fmt.Println("Installing...")
//return os.Rename("./MyApp", "/usr/bin/MyApp")
//}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "get", "github.com/stretchr/piglatin")
	return cmd.Run()
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("cast-via-mage")
}

func TestCLI() error {
	// go test -v ./...
	err := Build()
	if err != nil {
		return err
	}

	fmt.Println("starting server")

	// Start Go HTTP server on a specific port.
	http.HandleFunc("/testGet", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello\n")
	})

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		fmt.Println("Starting server on :1738")
		if err := http.ListenAndServe(":1738", nil); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		var cmdOut []byte
		fmt.Println("running test get")

		// Run commands using the newest "cast-via-mage" binary.
		cmd := exec.Command("./cast-via-mage", "get", "http://localhost:1738/testGet")
		cmdOut, err = cmd.CombinedOutput()
		fmt.Println(string(cmdOut))
		fmt.Println("Test Get Success")
		os.Exit(3)
	}()

	wg.Wait()
	return nil
}

func TestFileInput() {

}
