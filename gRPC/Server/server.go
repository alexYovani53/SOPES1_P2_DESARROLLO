// Paquete principal, acá iniciará la ejecución
package main

// Importar dependencias, notar que estamos en un módulo llamado grpctuiter
import (
	"context"
	"fmt"

	"os"

	"strconv"

	"google.golang.org/grpc"
	"log"
	"net"
	"tuiterserver/greet.pb"

	"math/rand"

	"github.com/joho/godotenv"
)

// Iniciar una estructura que posteriormente gRPC utilizará para realizar un server
type server struct{
}


func GameMaximo(jugadores int) int {
	return jugadores
}

func GameRandom(jugadores int) int {

	fmt.Println(">> SERVER",jugadores)

	val:= rand.Intn(jugadores)
	return val
}


func GameMedio(jugadores int) int {
	if jugadores < 2 {
		return jugadores
	}

	medio := 0
	medio = int(jugadores/2)

	return medio;
}

func GameModulo(jugadores int) int {
	if jugadores < 2 {
		return jugadores
	}

	medio := 0
	medio = int(jugadores/2)

	nuevoValor := medio * rand.Intn(100)

	modulo := nuevoValor % jugadores

	return modulo
}

func GamePenultimo(jugadores int) int {
	if jugadores < 2 {
		return jugadores
	}

	return jugadores - 1
}



// Función que será llamada desde el cliente
// Debemos pasarle un contexto donde se ejecutara la funcion
// Y utilizar las clases que fueron generadas por nuestro proto file
// Retornara una respuesta como la definimos en nuestro protofile o un error
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf(">> SERVER: Función Greet llamada con éxito. Datos: %v\n", req)

	// Todos los datos podemos obtenerlos desde req
	// Tendra la misma estructura que definimos en el protofile
	// Para ello utilizamos en este caso el GetGreeting
	juego := req.GetGreeting().GetJuego()
	nombreJuego := req.GetGreeting().GetNombreJuego()
	jugadores := req.GetGreeting().GetJugadores()

	ganador := 0
	jugadores2:= int(jugadores)

	if juego == 1 {
		ganador = GameMaximo(jugadores2)
	}else if juego == 2{
		ganador = GameRandom(jugadores2)
	}else if juego == 3 {
		ganador = GameMedio(jugadores2)
	}else if juego == 4 {
		ganador = GameModulo(jugadores2)
	}else if juego == 5 {
		ganador = GamePenultimo(jugadores2)
	}

	result := " ||| No.Juego -> "+strconv.FormatInt(juego,10) + "| nombreJuego -> " + nombreJuego + "| Jugadores -> " + strconv.FormatInt(jugadores,10) + "| Ganador-> " +strconv.Itoa(ganador) + " |||"

	fmt.Printf(">> SERVER: %s\n", result)
	// Creamos un nuevo objeto GreetResponse definido en el protofile
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

// Funcion principal
func main() {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("sin .env")
	}


	// Leer el host de las variables del ambiente
	host := os.Getenv("HOST")
	fmt.Println(">> SERVER: Iniciando en ", host)

	// Primero abrir un puerto para poder escuchar
	// Lo abrimos en este puerto arbitrario
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf(">> SERVER: Error inicializando el servidor: %v", err)
	}

	fmt.Println(">> SERVER: Empezando server gRPC")

	
	// Ahora si podemos iniciar un server de gRPC
	s := grpc.NewServer()

	// Registrar el servicio utilizando el codigo que nos genero el protofile
	greetpb.RegisterGreetServiceServer(s, &server{})

	fmt.Println(">> SERVER: Escuchando servicio...")
	// Iniciar a servir el servidor, si hay un error salirse
	if err := s.Serve(lis); err != nil {
		log.Fatalf(">> SERVER: Error inicializando el listener: %v", err)
	}
}
