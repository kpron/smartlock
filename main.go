package main

import "flag"
import "fmt"
import "net/http"
import "os/exec"
import "log"

func lock() {
	cmd := exec.Command("/bin/bash", "-c", "'/home/kpron/.config/i3/lock2'")
	out, err := cmd.Output()
	if err != nil {
		log.Println(err, "from code")
		log.Printf("%s", out)
	}
}

func unlock() {
	cmd := exec.Command("killall", "i3lock")
	err := cmd.Run()
	if err != nil {
		log.Println(err, "from code")
	}
}

func toggle(w http.ResponseWriter, r *http.Request, t string) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	parserr := r.ParseForm()
	if parserr != nil {
		log.Println("Cannot parse form", parserr)
	}
	key := r.PostFormValue("key")
	if key != t {
		log.Println("Wrong token from", r.RemoteAddr)
		return
	}
	log.Println("Correct token from", r.RemoteAddr)
	err := exec.Command("pgrep", "i3lock").Run()
	if err != nil {
		log.Println("not locked, locking now")
		go lock()
		w.Write([]byte("lolked"))
	} else {
		log.Println("locked, unlocking now")
		go unlock()
		w.Write([]byte("unlolked"))
	}
}
func main() {
	token := flag.String("token", "test", "secret token")
	port := flag.String("port", "10666", "port for listen on")
	flag.Parse()
	listen := fmt.Sprintf(":%v", *port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		toggle(w, r, *token)
	})
	http.ListenAndServe(listen, nil)
}
