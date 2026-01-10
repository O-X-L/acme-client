package u

import (
	"fmt"
	"log"
	"os"
)

var (
	MODE_DEV   = os.Getenv("MODE_DEV") == "1"
	MODE_TEST  = os.Getenv("MODE_TEST") == "1"
	MODE_DEBUG = os.Getenv("MODE_DEBUG") == "1"
)

func LogDevTest(m string) {
	if MODE_DEV || MODE_TEST || MODE_DEBUG {
		log.SetOutput(os.Stdout)
		log.Printf("[DEV] %v\n", m)
	}
}

func LogDev(m string) {
	if MODE_DEV || MODE_DEBUG {
		log.SetOutput(os.Stdout)
		log.Printf("[DEV] %v\n", m)
	}
}

func LogDebug(m string) {
	if MODE_DEBUG {
		log.SetOutput(os.Stdout)
		log.Printf("[DEBUG] %v\n", m)
	}
}

func LogDebugf(m string, params ...any) {
	LogDebug(fmt.Sprintf(m, params...))
}

func Log(m string) {
	log.SetOutput(os.Stdout)
	log.Printf("[INFO] %v\n", m)
}

func Logf(m string, params ...any) {
	Log(fmt.Sprintf(m, params...))
}

func LogError(m string) {
	log.SetOutput(os.Stderr)
	log.Printf("[ERROR] %v\n", m)
}

func LogErrorf(m string, params ...any) {
	LogError(fmt.Sprintf(m, params...))
}

func LogWarning(m string) {
	log.SetOutput(os.Stderr)
	log.Printf("[WARN] %v\n", m)
}

func LogWarningf(m string, params ...any) {
	LogWarning(fmt.Sprintf(m, params...))
}

func InitLogModes() {
	if MODE_DEV {
		Log("DEV-MODE enabled")
	}
	if MODE_TEST {
		Log("TEST-MODE enabled")
	}
	if MODE_DEBUG {
		Log("DEBUG-MODE enabled")
	}
}
