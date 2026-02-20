// package ui

// import (
// 	"embed"
// 	"io/fs"
// 	"net/http"
// )

// //go:embed assets/index.html
// var content embed.FS

// func Register(mux *http.ServeMux) {
// 	sub, err := fs.Sub(content, "assets")
// 	if err != nil {
// 		panic(err)
// 	}

// 	mux.Handle("/", http.FileServer(http.FS(sub)))

// 	mux.HandleFunc("/analytics/", func(w http.ResponseWriter, r *http.Request) {
// 		http.ServeFileFS(w, r, sub, "index.html")
// 	})
// }
package ui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed assets/index.html
var content embed.FS

func Register(mux *http.ServeMux) {
	sub, _ := fs.Sub(content, "assets")
	mux.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.FS(sub))))

	mux.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/", http.StatusFound)
	})
}