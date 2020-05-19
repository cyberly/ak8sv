package ak8sv

import (
	"fmt"
	"os"
)

func initEnvData(e string) string {
	v := os.Getenv(e)
	if len(v) == 0 {
		fmt.Printf("ERROR: Unable to read %v from environment!\n", e)
		os.Exit(1)
	}
	return v
}
