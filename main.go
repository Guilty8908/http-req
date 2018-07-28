package main

import (
	"net/http"
	"fmt"
	"runtime"
	"reflect"
	"time"
	"log"
)

func rootHandler(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(w, "hello world...")
}

func helloHandler(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(w, "hello ...")
}

func worldHandler(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(w, "world ...")
}



func logHandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request){
		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		fmt.Println("handler func called - ", name)
		h(w, req)
	}
}

func logHandler(h http.Handler) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Handler called - %T\n", h)
		h.ServeHTTP(w, r)
	})
}

func protectHandler(h http.Handler) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		// token...
		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		fmt.Printf("Handler called - %s\n", name)
		h.ServeHTTP(w, r)
	})
}

type myhandle struct{}

func (h *myhandle) ServeHTTP(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(w, "hello world...")
}

type hello struct{}
func (h *hello) ServeHTTP(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(w, "hello ...")
}

type world struct{}
func (h *world) ServeHTTP(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(w, "world...")
}



func main() {
	handle := myhandle{}
	handleHello := hello{}
	handleWorld:= world{}
	s := &http.Server{
		Addr:           ":7878",
		//Handler:        &handle,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	http.Handle("/", &handle)
	http.Handle("/hello/", protectHandler(logHandler(&handleHello)))
	http.Handle("/world/", &handleWorld)
	log.Fatal(s.ListenAndServe())
	//
	//http.HandleFunc("/", rootHandler)
	//http.HandleFunc("/hello/", logHandlerFunc(helloHandler)) // hello or hello/
	//http.HandleFunc("/world/", worldHandler)		//   /
	//
	//log.Fatal(http.ListenAndServe(":7878",nil))

	////err := http.ListenAndServeTLS(":7878", "cert.pem", "key.pem", nil)
	//log.Fatal(err)
}

