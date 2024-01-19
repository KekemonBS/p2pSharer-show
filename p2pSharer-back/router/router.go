package router

import (
	"net/http"
)

type trackerHandlersImpl interface {
	SaveUser(w http.ResponseWriter, r *http.Request)
	ReadUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	Save(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

// New returns router
func New(imp trackerHandlersImpl) *http.ServeMux {
	//API
	mux := http.NewServeMux()

	//Redis tbale #1 (active users)
	mux.HandleFunc("/tracker/init", imp.SaveUser)
	mux.HandleFunc("/tracker/check", imp.ReadUser)
	mux.HandleFunc("/tracker/leave", imp.DeleteUser)

	//Redis table #2 (folders)

	mux.HandleFunc("/tracker/host", imp.Save)
	mux.HandleFunc("/tracker/read", imp.Read)
	mux.HandleFunc("/tracker/remove", imp.Delete)

	prefixMux := http.NewServeMux()
	prefixMux.Handle("/api/v1/", http.StripPrefix("/api/v1", mux))

	return prefixMux
}
