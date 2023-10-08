package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

type server struct {
	fs http.Handler
	c  *counter
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s%s", r.Proto, r.Method, r.Host, r.URL)
	switch r.URL.Path {
	case "/count":
		fmt.Fprintln(w, s.c.i)
	case "/count/inc":
		if r.Method != http.MethodPost {
			http.Error(
				w,
				"method not allowed",
				http.StatusMethodNotAllowed,
			)
			return
		}
		i := s.c.increment()
		fmt.Fprintln(w, i)
	case "/count/dec":
		if r.Method != http.MethodPost {
			http.Error(
				w,
				"method not allowed",
				http.StatusMethodNotAllowed,
			)
			return
		}
		i := s.c.decrement()
		fmt.Fprintln(w, i)
	case "/":
		tmpl, err := template.ParseFiles("./html/index.html")
		if err != nil {
			http.Error(
				w,
				fmt.Sprintln(err),
				http.StatusInternalServerError,
			)
		}
		tmpl.Execute(w, s.c)
	default:
		s.fs.ServeHTTP(w, r)
	}
}

type counter struct {
	i int
	m sync.Mutex
}

func (c *counter) String() string {
	return fmt.Sprintf("%d", c.i)
}

func (c *counter) increment() int {
	c.m.Lock()
	defer c.m.Unlock()
	time.Sleep(400 * time.Millisecond)
	c.i++
	return c.i
}

func (c *counter) decrement() int {
	c.m.Lock()
	defer c.m.Unlock()
	time.Sleep(400 * time.Millisecond)
	c.i--
	return c.i
}

func main() {
	s := server{
		fs: http.FileServer(http.Dir("./html")),
		c:  new(counter),
	}
	log.Print("Server starting on port 8000")
	log.Fatal(http.ListenAndServe(":8000", &s))
}
