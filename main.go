package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/reportes/ok", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id": 101, "mensaje": "Reporte creado con exito", "estado": "activo"}`)
	})

	http.HandleFunc("/reportes/error", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, `{"error": "Faltan datos obligatorios: nombre_mascota y foto_url"}`)
	})

	fmt.Println("Servidor falso corriendo en el puerto 8080...")
	http.ListenAndServe(":8080", nil)
}
