package u

import (
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
		log.Printf("[DEV] %v\n", m)
	}
}

func LogDev(m string) {
	if MODE_DEV || MODE_DEBUG {
		log.Printf("[DEV] %v\n", m)
	}
}

func LogDebug(m string) {
	if MODE_DEBUG {
		log.Printf("[DEBUG] %v\n", m)
	}
}

func Log(m string) {
	log.Printf("[INFO] %v\n", m)
}

func LogError(m string) {
	log.Printf("[ERROR] %v\n", m)
}

func LogWarning(m string) {
	log.Printf("[WARN] %v\n", m)
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
