package helper

import (
	"fmt"
	"log"
)

// PauseProcess pause program and wait user press Enter to exit
func PauseProcess(erors []error) {
    if len(erors) == 0 {
        fmt.Printf("All done!\n\n")
    } else {
        fmt.Printf("Error!\n\n")
        for _, err := range erors {
            log.Println(err)
        }
    }

    fmt.Printf("\nPress Enter Key to exit...")
    fmt.Scanln()
}