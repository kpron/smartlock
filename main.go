package main

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

func toggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		toggle(w, r)
	})
	http.ListenAndServe(":10666", nil)
}
