package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var (
	infoLogger     *log.Logger
	errorLogger    *log.Logger
	azSubId        = flag.String("AZURE_SUBSCRIPTION_ID", "", "Azure Subscription ID.")
	azTenantId     = flag.String("AZURE_TENANT_ID", "", "AZURE_TENANT_ID. You save this after az login")
	azClientID     = flag.String("AZURE_CLIENT_ID", "", "AZURE_CLIENT_ID. You save this after az login")
	azClientSecret = flag.String("AZURE_CLIENT_SECRET", "", "AZURE_CLIENT_SECRET. You save this after az login")
)

func init() {
	infoLogger = log.New(os.Stdout, "INFO ", log.LstdFlags) // 5 -> 2023/05/24 14:52:51.507584, 4 -> 14:52:57.482177
	errorLogger = log.New(os.Stdout, "ERROR ", log.LstdFlags)
	flag.Parse()
}

func main() {
	os.Exit(run())
}

func run() int {

	if *azSubId == "" {
		value, isDefined := os.LookupEnv("AZURE_SUBSCRIPTION_ID")
		if isDefined {
			*azSubId = value
		}

		if !isDefined || value == "" {
			errorLogger.Println("neither azSubId nor $AZURE_SUBSCRIPTION_ID was defined. Exiting")
			return 1
		}
	}

	if *azTenantId == "" {
		value, isDefined := os.LookupEnv("AZURE_TENANT_ID")
		if isDefined {
			*azTenantId = value
		}

		if !isDefined || value == "" {
			errorLogger.Println("neither azTenantId nor $AZURE_TENANT_ID was defined. Exiting")
			return 1
		}
	}

	if *azClientID == "" {
		value, isDefined := os.LookupEnv("AZURE_CLIENT_ID")
		if isDefined {
			*azClientID = value
		}

		if !isDefined || value == "" {
			errorLogger.Println("neither azClientID nor $AZURE_CLIENT_ID was defined. Exiting")
			return 1
		}
	}

	if *azClientSecret == "" {
		value, isDefined := os.LookupEnv("AZURE_CLIENT_SECRET")
		if isDefined {
			*azClientSecret = value
		}

		if !isDefined || value == "" {
			errorLogger.Println("neither azClientSecret nor $AZURE_CLIENT_SECRET was defined. Exiting")
			return 1
		}
	}

	return 1
}
