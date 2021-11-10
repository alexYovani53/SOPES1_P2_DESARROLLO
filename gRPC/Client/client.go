// Paquete principal, acá iniciará la ejecución
package main

// Importar dependencias, notar que estamos en un módulo llamado tuiterclient
import (
	"context"

	"os"

	"fmt"
	"log"

	"net/http"

	"tuiterclient/greet.pb"

	"google.golang.org/grpc"
	"encoding/json"

	"github.com/joho/godotenv"
)


type Juego struct{
	Juego int `json:"juego"`
	NombreJuego string `json:"nombreJuego"`
	Jugadores int `json:"jugadores"`
}

type server struct{}

// Funcion que realiza una llamada unaria
func sendMessage(Juego int, NombreJuego string, Jugadores int) {
	server_host := os.Getenv("SERVER_HOST")

	fmt.Println(">> CLIENT: Iniciando cliente")
	fmt.Println(">> CLIENT: Iniciando conexion con el servidor gRPC ", server_host)

	// Crear una conexion con el servidor (que esta corriendo en el puerto 50051)
	// grpc.WithInsecure nos permite realizar una conexion sin tener que utilizar SSL
	cc, err := grpc.Dial(server_host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf(">> CLIENT: Error inicializando la conexion con el server %v", err)
	}

	// Defer realiza una axion al final de la ejecucion (en este caso, desconectar la conexion)
	defer cc.Close()

	// Iniciar un servicio NewGreetServiceClient obtenido del codigo que genero el protofile
	// Esto crea un cliente con el cual podemos escuchar
	// Le enviamos como parametro el Dial de gRPC
	c := greetpb.NewGreetServiceClient(cc)

	fmt.Println(">> CLIENT: Iniciando llamada a Unary RPC")

	// Crear una llamada de GreetRequest
	// Este codigo lo obtenemos desde el archivo que generamos con protofile
	req := &greetpb.GreetRequest{
		// Enviar un Greeting
		// Esta estructura la obtenemos desde el archivo que generamos con protofile
		Greeting: &greetpb.Greeting{
			Juego: int64(Juego),
			NombreJuego: NombreJuego,
			Jugadores: int64(Jugadores),
		},
	}

	fmt.Println(">> CLIENT: Enviando datos al server")
	// Iniciar un greet, en background con la peticion que estamos realizando
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf(">> CLIENT: Error realizando la peticion %v", err)
	}

	fmt.Println(">> CLIENT: El servidor nos respondio con el siguiente mensaje: ", res.Result)
}

// Creamos un server sencillo que unicamente acepte peticiones GET y POST a '/'
func http_server(w http.ResponseWriter, r *http.Request) {
	instance_name := os.Getenv("NAME")
	fmt.Println(">> CLIENT: Manejando peticion HTTP CLIENTE: ", instance_name)
	// Comprobamos que el path sea exactamente '/' sin parámetros
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Comprobamos el tipo de peticion HTTP
	switch r.Method {
	// Devolver una página sencilla con una forma html para enviar un mensaje
	case "GET":
		fmt.Println(">> CLIENT: Devolviendo form.html")
		// Leer y devolver el archivo form.html contenido en la carpeta del proyecto
		http.ServeFile(w, r, "form.html")

	// Publicar un mensaje a Google PubSub
	case "POST":
		fmt.Println(">> CLIENT: Iniciando envio de mensajes")
		// Si existe un error con la forma enviada entonces no seguir

		var p Juego

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil{
			http.Error(w,err.Error(),http.StatusBadRequest)
			return
		}

	
		fmt.Println(p)
		// Publicar el mensaje, convertimos el objeto JSON a String
		sendMessage(p.Juego, p.NombreJuego,p.Jugadores)



		// Enviamos informacion de vuelta, indicando que fue generada la peticion
		fmt.Fprintf(w, "¡Juego completado!\n")
		fmt.Fprintf(w, "NombreJuego = %s\n", p.NombreJuego)
		fmt.Fprintf(w, "Juego = %d\n", p.Juego)
		fmt.Fprintf(w, "Jugadores = %d\n", p.Jugadores)

	// Cualquier otro metodo no sera soportado
	default:
		fmt.Fprintf(w, "Metodo %s no soportado \n", r.Method)
		return
	}
}

// Funcion principal
func main() {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("sin .env")
	}

	instance_name := os.Getenv("NAME")
	client_host := os.Getenv("CLIENT_HOST")

	fmt.Println(">> -------- CLIENTE ", instance_name, " --------")

	fmt.Println(">> CLIENT: Iniciando servidor http en ", client_host)

	// Asignar la funcion que controlara las llamadas http
	http.HandleFunc("/", http_server)

	// Levantar el server, si existe un error levantandolo hay que apagarlo
	if err := http.ListenAndServe(client_host, nil); err != nil {
		log.Fatal(err)
	}


	
}
