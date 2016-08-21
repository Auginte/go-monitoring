package common

import "log"

// LogError - expecting happy path, otherwise log and os.Exit(1)
func LogError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
