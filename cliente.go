package main

import (
	"fmt"
	"net/rpc"
	"bufio"
	"os"
	"io/ioutil"
	//"strings"
)
type Persona struct{
	Nombre string
	Mensaje string
	Archivo []byte
	Archivos string
	Nombre_A string
}
type error interface {
    Error() string
}
type errorString struct {
    s string
}
func (e *errorString) Error() string {
    return e.s
}
var persona Persona
func cliente()  {
	c,err:=rpc.Dial("tcp","127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	var op, result int64
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	persona.Nombre = scanner.Text()
	err = c.Call("Server.Conexion", persona.Nombre,&result)
	if err != nil {
		fmt.Println(err)
	}
	go chat()
	for{
		fmt.Println("\t- Menu -")
		fmt.Println("1.- Enviar Mensaje")
		fmt.Println("2.- Enviar Archvio")
		fmt.Println("3.- Mensajes Servidor")
		fmt.Println("4.- Salir")
		fmt.Scan(&op)
		if op == 4{
			err = c.Call("Server.Desconectar", persona,&result)
			if err != nil {
				fmt.Println(err)
			}
			break
		}else if op ==1{
			fmt.Println("Mensaje es: ")
			scanner.Scan()
			scanner.Scan()
			persona.Mensaje = scanner.Text()
			err = c.Call("Server.Mensaje", persona,&result)
			if err != nil {
				err.Error()
			}
		}else if op == 2{
			var archivo string
			fmt.Println("Ruta del Archivo a enviar" )
			scanner.Scan()
			scanner.Scan()
			archivo = scanner.Text()
			f, err := ioutil.ReadFile(archivo)
			if err != nil{
				fmt.Println(err)
			}
			persona.Nombre_A = archivo
			persona.Archivo = f
			k := c.Call("Server.Archivos", persona,&result)
			if k != nil{
				fmt.Println(k)
			}
		}else if op == 3{
			err = c.Call("Server.Todo", "",&result)
			if err != nil {
				fmt.Println("Mensajes en Servidor:\n",err)
			}
		}
	}
}
func chat()  {
	var result int64
	c,err:=rpc.Dial("tcp","127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for{
		err = c.Call("Server.Ret","",&result)
		if err != nil{
			fmt.Println(err)
		}
		k := c.Call("Server.Compartir", "",&result)
		if k != nil{
			l:= k.Error()
			destino, err := os.OpenFile("destino.txt", os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				fmt.Println(err)
			}
			destino.WriteString(string(l))
			defer destino.Close()
		}
	}
}
func main()  {
	cliente()
}