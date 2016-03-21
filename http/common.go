package http

import (
	/*"log"*/
	"net/http"

	"github.com/dockermec/g"
)

func configCommonRoutes() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(g.VERSION))
	})
	//deal with the imformation of replication controll
	http.HandleFunc("/gethostinfo", func(w http.ResponseWriter, r *http.Request) {

	})

}
