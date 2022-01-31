package injection

import (
	"net/http"
	handlers "new/internal/article/delivery/http"

	"github.com/gorilla/mux"
)

func StartApp() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/", handlers.Index).Methods("GET")
	rtr.HandleFunc("/create", handlers.Create).Methods("GET")
	rtr.HandleFunc("/save_article", handlers.Save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", handlers.Show_post).Methods("GET")
	rtr.HandleFunc("/delete/{id:[0-9]+}", handlers.Delete_post).Methods("POST") // https://medium.com/@carlospineda/why-no-methods-for-put-delete-in-html-f483b66d8874

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)
}