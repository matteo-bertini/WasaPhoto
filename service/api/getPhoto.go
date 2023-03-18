package api

import (
	"WasaPhoto/service/api/reqcontext"
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	photoid := strings.Split(r.URL.Path, "/")[4]
	urlusername := strings.Split(r.URL.Path, "/")[2]
	id2, err := rt.db.IdFromUsername(urlusername)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Si è verificato un errore nell'estrazione dell'id dell'username nell'URL.")
		return
	}
	path := "/tmp/WasaPhoto/" + *id2 + "/" + photoid + ".png"
	photofile, err := os.Open(path)
	if err != nil {
		ctx.Logger.WithError(err).Error("La foto cercata non esiste.")
		w.WriteHeader(http.StatusNotFound)
		return
	} else {

		// Setting dell'header della risposta
		w.Header().Set("Content-Type", "image/png")

		buf := bytes.NewBuffer(nil)
		_, err := io.Copy(buf, photofile)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Si è verificato un errore nella copia sul buffer.")
			return
		} else {
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(buf.Bytes())
			if err != nil {
				ctx.Logger.WithError(err).Error("Si è verificato un errore nella scrittura della riposta.")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
	}
}
