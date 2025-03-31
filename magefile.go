//go:build mage
// +build mage

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
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
	cmd := exec.Command("go", "build", "-o", "cast", ".")
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
	os.RemoveAll("cast")
}

func TestCLI() error {
	// go test -v ./...
	err := Build()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		fmt.Println("Starting server on :1738")
		testHTTPServer()
	}()

	go func() {

		testCLIget()
		os.Exit(3)
	}()

	wg.Wait()
	return nil
}

func TestFileInput() {

}

func testCLIget() {
	var cmdOut []byte
	var err error
	fmt.Println("running test get")

	// Run commands using the newest "cast-via-mage" binary.
	cmd := exec.Command("./cast", "get", "http://localhost:1738/testGet")
	cmdOut, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error encountered: %s", err)
	}
	fmt.Println(string(cmdOut))
	if strings.Contains(string(cmdOut), "test value 1") {
		fmt.Println("Test CLI GET Success")
	} else {
		//fmt.Println("Test GET Failed")
		panic("Test CLI GET Failed")
	}
}

func testHTTPServer() {
	fmt.Println("starting server")

	// Start Go HTTP server on a specific port.
	http.HandleFunc("/testGet", func(w http.ResponseWriter, req *http.Request) {
		respBody := `
hello\n

<html>
<h1>test value 1</h1>
<p>this is another test area to grab a specific value</p>
</html>
		`

		fmt.Fprint(w, respBody)
	})

	http.HandleFunc("/testJSON", func(w http.ResponseWriter, req *http.Request) {
		respBody := `
{
"testVal1": "testVal2",
"name2": "name3",
"name4": true
}
		`

		fmt.Fprint(w, respBody)
	})

	if err := http.ListenAndServe(":1738", nil); err != nil {
		log.Fatal(err)
	}

	return
}
