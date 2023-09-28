package helpers

import (
	"os"
)

//WriteFile se encarga de escribir los logs
func WriteFile(datosGuarda []byte, nombre string) {
	f, err := os.OpenFile(nombre, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	ValidError(err)
	if _, err := f.Write(datosGuarda); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}
