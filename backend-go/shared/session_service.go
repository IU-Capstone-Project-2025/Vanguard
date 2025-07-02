package shared

import (
	"fmt"
	"os"
)

const wsEndpoint = "/ws"

func GetWsEndpoint() string {
	return fmt.Sprintf("ws://%s:%s/ws",
		os.Getenv("REALTIME_SERVICE_HOST"), os.Getenv("REALTIME_SERVICE_PORT"))
}

const CodeLength int = 6
