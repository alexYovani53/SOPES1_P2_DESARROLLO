package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"io"
	"strconv"
	"time"


	"net/http"
	"bytes"
	"encoding/json"

	"math/rand"
	"io/ioutil"
)

type comandoCLI struct{
	Timeout		int 
	Concurrence	int
	Rungames	int
	Players 	int
	IndiceJuegos		[]	int
	NombreJuegos		[]  string
}

type Juego struct{
	Juego int `json:"juego"`
	NombreJuego string `json:"nombreJuego"`
	Jugadores int `json:"jugadores"`
}

func main(){

	// path : ruta del archivo a leer
	var comando string
	var url 	string

	fmt.Printf("CLI $> Ingrese comando de juego \n")
	fmt.Printf("CLI $> ")


	var stdin *bufio.Reader
	var line []rune
	var c rune
	var err error

	stdin = bufio.NewReader(os.Stdin)

	fmt.Printf("Type something: ")

	for {
			c, _, err = stdin.ReadRune()
			if err == io.EOF || c == '\n' { break }
			if err != nil {
					fmt.Fprintf(os.Stderr,"Error reading standard input\n")
					os.Exit(1)
			}
			line = append(line,c)
	}


	comando = string(line[:len(line)])

	fmt.Printf("CLI $> Ingrese la url a enviar los juegos ")
	fmt.Scanf("%s\n",&url)

	estructura := splitComando(comando)

	ValidarEnvio(estructura,url)
}

func splitComando(comando string) comandoCLI{

	var indices [] int
	var nombres [] string
	var players string 
	var rungames string
	var concurrence string
	var timeout string

	comandoSeparado := strings.Split(comando,"--");
	
	for i := 0; i <=len(comandoSeparado) - 1; i++ {
		fmt.Println(comandoSeparado[i])
	}


	/*OBTENIENDO JUEGOS*/
	juegosLista := strings.TrimLeft(comandoSeparado[1],"gamename");
	juegosLista = strings.TrimSpace(juegosLista)
	juegosLista = strings.TrimLeft(juegosLista," ")
	juegosLista = strings.TrimLeft(juegosLista,"\"")
	juegosLista = strings.TrimRight(juegosLista," ")
	juegosLista = strings.TrimRight(juegosLista,"\"")

	fragmentosJuegos :=strings.Split(juegosLista,"|");



	for i := 0; i<= len(fragmentosJuegos) - 1; i++{
		
		if i%2 == 0{
			idJuego := strings.TrimSpace(fragmentosJuegos[i])
			if res,err := strconv.Atoi(idJuego)
			err == nil{				
				indices = append(indices,res);
			}
		}else{
			nombres = append(nombres,strings.TrimSpace(fragmentosJuegos[i]))
			fmt.Println(strings.TrimSpace(fragmentosJuegos[i]))
		}

	}


	/*OBTENIENDO players*/
	
	players = strings.TrimLeft(comandoSeparado[2],"players")
	players = strings.TrimSpace(players)

	/*OBTENIENDO rungames*/

	rungames = strings.TrimLeft(comandoSeparado[3],"rungames")
	rungames = strings.TrimSpace(rungames)

	/*OBTENIENDO concurrence*/

	concurrence = strings.TrimLeft(comandoSeparado[4],"concurrence")
	concurrence = strings.TrimSpace(concurrence)

	/*OBTENIENDO timeout*/

	timeout = strings.TrimLeft(comandoSeparado[5],"timeout")
	timeout = strings.TrimSpace(timeout)
	timeout = strings.TrimRight(timeout,"m")


	var comandoStruct comandoCLI

	if i,err := strconv.Atoi(timeout)
	err == nil{
		comandoStruct.Timeout = i
	}

	if i,err := strconv.Atoi(concurrence)
	err == nil{
		comandoStruct.Concurrence = i
	}

	if i,err := strconv.Atoi(rungames)
	err == nil{
		comandoStruct.Rungames = i
	}

	if i,err := strconv.Atoi(players)
	err == nil{
		comandoStruct.Players = i
	}

	comandoStruct.IndiceJuegos = indices	
	comandoStruct.NombreJuegos = nombres

	
	printSlice(comandoStruct.IndiceJuegos)
	fmt.Println(comandoStruct.Players)
	fmt.Println(comandoStruct.Rungames)
	fmt.Println(comandoStruct.Concurrence)
	fmt.Println(comandoStruct.Timeout)

	return comandoStruct

}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func ValidarEnvio(comando comandoCLI, url_access string) {
	var contador int
	var finish bool
	//start time
	start_time := time.Now()
	/// se define el tiempo maximo en minutos
	timeout := start_time.Add(time.Minute * time.Duration(comando.Timeout))
	for {
		for i := 0; i < len(comando.IndiceJuegos); i++ {
			actual_time := time.Now()
			//timeout break
			if actual_time.After(timeout) {
				fmt.Println("Tiempo alcanzado, se enviaron:", contador, " juegos a procesar")
				finish = true
				break
			}
			//concurrence: create n gorutines
			for j := 0; j < comando.Concurrence; j++ {
				//increment game_counter
				contador++


				///cambiar por funcion de envio a grpc
				go func() {
					fmt.Println("url: ",url_access);

					//resp, err:=http.Get("http://34.69.79.12.nip.io/myapp");

					indice := rand.Intn(len(comando.IndiceJuegos))
					jugadores := rand.Intn(comando.Players)

					data := Juego{comando.IndiceJuegos[indice],comando.NombreJuegos[indice],jugadores}
					data2,_:= json.Marshal(data)

					resp, err:= http.Post("http://localhost:5001","application/json",
					bytes.NewBuffer(data2))

					if err !=nil{						
						fmt.Printf("No se ha podido enviar la información: %s\n", err)
					}else{    

						defer resp.Body.Close()

						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							fmt.Println(err)
						}
						fmt.Println(string(body))

						//fmt.Printf("informacion enviada: %s\n",resp.Body)
					}


				}()
				//max games validation
				if contador >= int(comando.Rungames) {
					fmt.Println("Se enviaron todos los juegos a procesar")
					finish = true
					break
				}
			}
			time.Sleep(time.Second)
		}
		//finish loop
		if finish {
			break
		}
	}
}