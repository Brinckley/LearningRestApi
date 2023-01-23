package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

const keyServerAddr = "serverAddr"

// handler funcs
func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s, got from / request", ctx.Value(keyServerAddr))
	_, err := io.WriteString(w, "Hello, this is my house ;o")
	if err != nil {
		fmt.Println("Error writing with ResponseWriter")
		return
	}
}

// ResponseWriter writes the info back to client
// Request gets the info about the request got from user, like body of post
func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s, got from /hello request", ctx.Value(keyServerAddr))

	myName := r.PostFormValue("myName")
	if myName == "" {
		w.Header().Set("x-missing-field", "myName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := io.WriteString(w, fmt.Sprintf("Glad to see you, %s!", myName))
	if err != nil {
		fmt.Println("Error writing /hello")
		return
	}
}

func getParams(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	fmt.Printf("%s. Param1 found - %v, param1 val - %s, Param2 found - %v, param2 val - %s",
		ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second)

	_, err := io.WriteString(w, "this is my params getter, btw")
	if err != nil {
		fmt.Println("Error writing /params")
		return
	}
}

func getBody(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	fmt.Printf("%s, body value is : %s\n", ctx.Value(keyServerAddr), string(body))

	_, err = io.WriteString(w, "Got your body")
	if err != nil {
		fmt.Println("Error writing /body")
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)
	mux.HandleFunc("/params", getParams)
	mux.HandleFunc("/body", getBody)

	ctx, cancelCtx := context.WithCancel(context.Background())
	serverOne := &http.Server{
		Addr:    ":8081",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	go func() {
		err := serverOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server one closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server one: %s\n", err)
		}
		cancelCtx()
	}()

	//serverTwo := &http.Server{
	//	Addr:    ":8082",
	//	Handler: mux,
	//	BaseContext: func(l net.Listener) context.Context {
	//		ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
	//		return ctx
	//	},
	//}

	//go func() {
	//	err := serverTwo.ListenAndServe()
	//	if errors.Is(err, http.ErrServerClosed) {
	//		fmt.Printf("server one closed\n")
	//	} else if err != nil {
	//		fmt.Printf("error listening for server one: %s\n", err)
	//	}
	//	cancelCtx()
	//}()

	<-ctx.Done()
}
