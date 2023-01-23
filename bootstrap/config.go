package bootstrap

import (
	"log"
	"os"
	"strconv"
)

func getDebugMode() bool {
	isDebugMode := false
	debugModeString := os.Getenv("DEBUG_MODE")

	if debugModeString != "" {
		debugMode, err := strconv.ParseBool(debugModeString)
		if err != nil {
			log.Fatalln("environment DEBUG_MODE must be bool")
		}
		isDebugMode = debugMode
	}

	return isDebugMode
}
