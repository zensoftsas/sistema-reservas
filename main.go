
package main

import "fmt"

func saludar(nombre string) string {
    return fmt.Sprintf("Hola, %s", nombre)
}

func main() {
    mensaje := saludar("Clínica Internacional")
    fmt.Println("Sistema de Reservas -", mensaje)
}
