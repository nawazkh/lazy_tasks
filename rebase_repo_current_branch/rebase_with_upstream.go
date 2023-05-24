package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger

	reposRoot       = flag.String("reposRoot", "", "The root of all the repos you want to rebase with UPSTREAM.")
	defaultUpstream = flag.String("defaultUpstream", "main", "The default source of truth at UPSTREAM. It could be `main` or `master`.")
)

func init() {
	infoLogger = log.New(os.Stdout, "INFO", log.LstdFlags)
	errorLogger = log.New(os.Stdout, "ERROR", log.LstdFlags)
	debugLogger = log.New(os.Stdout, "DEBUG", log.LstdFlags)

	flag.Parse()
}

func main() {
	os.Exit(run())
}

func run() int {
	// Setup

	// get all the dirs in root

	// launch parallel go routines to fetch UPSTREAM and attempt rebase

	// return 0 on all successful rebases
	// return any error code if any one of them fails

	// find the right root
	if reposRoot != nil {
		if *reposRoot == "" {
			value, isDefined := os.LookupEnv("UPSTREAM")
			if isDefined {
				*reposRoot = value
			}

			if !isDefined || value == "" {
				// Select current path as the root
				os.getP
			}
		}
	}
	log.Println()
	fmt.Println("Root of all the repositories", reposRoot)
	return 0
}

// Filename is the __filename equivalent
func Filename() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filename, nil
}

// Dirname is the __dirname equivalent
func Dirname() (string, error) {
	filename, err := Filename()
	if err != nil {
		return "", err
	}
	return filepath.Dir(filename), nil
}
