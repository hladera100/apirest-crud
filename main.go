package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "apirest_crud"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {

	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Println("Servidor corriendo...")
	http.ListenAndServe(":2001", nil)

}

type Empleado struct {
	ID     int
	Nombre string
	Correo string
}

func Inicio(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Gola Crud")
	log.Println("Funcion Inicio ")
	conexionEstablecida := conexionDB()
	registros, err := conexionEstablecida.Query("SELECT * FROM empleado ")
	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for registros.Next() {
		var id int
		var nombre, correo string

		err = registros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.ID = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)

	}

	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Gola Crud")
	log.Println("Funcion crear")
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	log.Println("Funcion Insertar ")
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")
		conexionEstablecida := conexionDB()
		insertRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleado (nombre, correo) VALUES (?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insertRegistros.Exec(nombre, correo)
		http.Redirect(w, r, "/", 301)
	}
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	log.Println("Funcion Borrar ")
	id := r.URL.Query().Get("id")
	conexionEstablecida := conexionDB()
	borrarReg, err := conexionEstablecida.Prepare("DELETE FROM empleado WHERE ID =?")
	if err != nil {
		panic(err.Error())
	}
	borrarReg.Exec(id)
	http.Redirect(w, r, "/", 301)
}

func Editar(w http.ResponseWriter, r *http.Request) {
	log.Println("Funcion Editar ")
	id := r.URL.Query().Get("id")
	conexionEstablecida := conexionDB()
	registro, err := conexionEstablecida.Query("SELECT * FROM empleado WHERE id= ?", id)
	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}

	for registro.Next() {
		var id int
		var nombre, correo string

		err = registro.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.ID = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}
	fmt.Println(empleado)
	plantillas.ExecuteTemplate(w, "editar", empleado)
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	log.Println("Funcion Actualizar ")
	id := r.FormValue("id")
	nombre := r.FormValue("nombre")
	correo := r.FormValue("correo")
	log.Println("id: ", id)
	log.Println("Nombre: ", nombre)
	log.Println("correo: ", correo)

	conexionEstablecida := conexionDB()
	registro, err := conexionEstablecida.Prepare("UPDATE empleado SET nombre=?, correo= ? WHERE id= ?")
	if err != nil {
		panic(err.Error())
	}
	registro.Exec(nombre, correo, id)
	http.Redirect(w, r, "/", 301)
}
