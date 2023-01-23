package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var dit Dict

type Dict struct {
	elems [][]byte
}

type Router struct {
	routes []Route
}

type Route struct {
	pattern string
	method  string
	action  http.HandlerFunc
}

func (d *Dict) set(full []byte) string {
	d.elems = append(d.elems, full)

	return "https://" + "localhost:8080" + "/" + strconv.Itoa(len(d.elems)-1)
}

func (d *Dict) get(id int) ([]byte, error) {
	if len(d.elems) == 0 {
		return nil, errors.New("")
	}

	el := d.elems[id]

	fmt.Println("", el)

	if len(el) != 0 {
		return el, nil
	}

	return nil, errors.New("")
}

func RootGetHandler(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
	id, err := getID(r.URL.Path)

	if err != nil {
		http.Error(w, "", http.StatusNotFound)

		return
	}

	el, err := dit.get(id)

	if err != nil {
		http.Error(w, "", http.StatusNotFound)

		return
	}

	http.Redirect(w, r, string(el), http.StatusTemporaryRedirect)
}

func RootPostHandler(w http.ResponseWriter, r *http.Request) {
	url, err := io.ReadAll(r.Body)

	if err != nil || len(url) == 0 {
		http.Error(w, "", http.StatusBadRequest)

		return
	}

	short := dit.set(url)

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(short))

	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
	}
}

func main() {
	routes := []Route{
		{
			pattern: "/",
			method:  http.MethodGet,
			action:  RootGetHandler,
		},
		{
			pattern: "/",
			method:  http.MethodPost,
			action:  RootPostHandler,
		},
	}
	router := new(Router)
	router.routes = routes

	err := http.ListenAndServe("localhost:8080", router)

	if err != nil {
		log.Panic(err)
	}
}

func getID(s string) (int, error) {
	id := strings.Split(s, "/")[1]

	fmt.Println("id", id)

	return strconv.Atoi(id)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route, err := r.find(req)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(""))
		if err != nil {
			return
		}
	} else {
		route.action(w, req)
	}
}

func (r *Router) find(req *http.Request) (Route, error) {
	for _, route := range r.routes {
		if route.method == req.Method {
			return route, nil
		}
	}

	return Route{}, errors.New("")
}
