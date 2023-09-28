package helpers

import (
	"fmt"
	"time"
)

//ValidError valida los errores Fatales que pueden surgir durante el codigo
func ValidError(err error) {
	if err != nil {
		WriteFile([]byte(fmt.Sprintf("%s:::%s\n", err.Error(), time.Now().Format("2006-01-02 15:04:05"))), "log.txt")
	}
}
