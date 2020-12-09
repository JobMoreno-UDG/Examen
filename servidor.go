package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strings"
)
type Server struct{}
type Persona struct{
	Nombre string
	Mensaje string
	Archivo []byte
	Archivos string
	Nombre_A string
}
type Clientes struct {
	P []Persona
	Mensaje,Temp string
	lista_Mensajes []string
	Archivo []byte
	Nombre_A string
}
var lista Clientes
type Error struct{
	msg string
}
func NewErrGOT(mensaje string) *Error {
    return &Error{
        msg: mensaje,
    }
}
func (e *Error) Error() string {
    return fmt.Sprintf("%s ", e.msg)
}
///// Archvio
// mensaje a servidor
func (this *Server)Mensaje(datos Persona, reply *int64)error{
	nombre := datos.Nombre
	mensaje:= datos.Mensaje	
	lista.Temp = nombre +" : "+mensaje
	fmt.Println(lista.Temp)
	lista.lista_Mensajes=append(lista.lista_Mensajes,nombre+" : "+mensaje)
	return nil
}
// archivos
func (this *Server)Archivos(datos Persona, reply *int64)error {
	lista.Nombre_A = datos.Nombre_A
	lista.Archivo = datos.Archivo
	l := strings.Split(datos.Nombre_A,"\\")	
	lista.Temp = datos.Nombre +" : "+l[len(l)-1]		
	lista.lista_Mensajes=append(lista.lista_Mensajes,lista.Nombre_A+" : "+l[len(l)-1])
	return nil
}
func (this *Server)Compartir(datos string, reply *int64)error {
	if  lista.Nombre_A != ""{
		lista.Nombre_A=""
		return NewErrGOT(string(lista.Archivo))		
	}else{
		return nil
	}
}
//conexion clientes
func (this *Server)Conexion(c string, reply *int64)error{
	var persona Persona
	persona.Nombre = c
	fmt.Println("Se Conecto: ", c)
	lista.P=append(lista.P,persona)
	lista.Mensaje = ""
	lista.Temp = ""
	return nil
}
func server() {
	rpc.Register(new(Server))
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}
}
func menu(){
	var op int
	go server()
	for{		
		fmt.Println("Menu")
		fmt.Println("1) Mostrar Mesajes/Nombre Archivos")
		fmt.Println("2) Respaldar Mensajes/Nombre Archivos")
		fmt.Println("3) Terminar")
		fmt.Scan(&op)
		if op ==3{
			return
		}else if op ==1{
			for _,i := range lista.lista_Mensajes{
				fmt.Println(i)
			}
		}else if op == 2{
			file,err:=os.Create("Respaldo.txt")
			if err != nil{
				fmt.Println(err)
				return
			}
			for _,i := range lista.lista_Mensajes{
				file.WriteString(i+"\n")
			}
			defer file.Close()
		}
	}
}
//conexion clientes
func (this *Server)Ret(c string, reply *int64)error{
	for{
		if lista.Mensaje != lista.Temp{
			lista.Mensaje = lista.Temp
			return NewErrGOT(lista.Mensaje)				
		}
	}
}
// Desconexion Clientes
func (this *Server)Desconectar(p Persona, reply *int64)error{
	cont:=0
	fmt.Println("Se Desconecto: ", p.Nombre)
	for _,i := range lista.P{
		if i.Nombre == p.Nombre{
			lista.P=append(lista.P[:cont],lista.P[cont+1:]...)
		}
		cont +=1
	}
	return nil
}
// Envia todo los mensajes
func (this *Server)Todo(dato string, reply *int64)error  {
	final:=""
	for _,i := range lista.lista_Mensajes{
		final += i +"\n"
	}
	return NewErrGOT(final)
}
func main() {
	menu()
}