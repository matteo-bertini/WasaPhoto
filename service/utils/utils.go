package utils

import (
	"errors"
	"strings"
)

// Errori

// CheckAuthorization //
var ErrorAuthorizationNotSpecified error = errors.New("Authorization non specificata nell'header.")
var ErrorBearerTokenNotSpecifiedWell error = errors.New("Bearer Token non specificato correttamente nel campo Authorization dell'header.")
var ErrorUnauthorized error = errors.New("L'id passato nel campo Authorization non è autorizzato ad effettuare l'operazione.")

// AddUser //
var ErrorUserAlreadyExists error = errors.New("L'utente è già presente nel sistema quindi non è stato creato un nuovo profilo.")

// FollowUser //
var ErrorFollowerAlreadyAdded error = errors.New("L'user ")

// BanUser //
var ErrorUserAlreadyBanned error = errors.New("L'user è già bannato.")

// CheckUserExistence //
var ErrUserDoesNotExist error = errors.New("L'utente cercato non esiste: non ha ancora creato un profilo o non è ancora registrato.")

// IsAllowed //
var ErrUserNotAllowed error = errors.New("L'utente cercato non è autorizzato ad ottenere le informazioni.")

// CheckPhotoExistence
var ErrPhotoDoesNotExist error = errors.New("La foto non esiste.")

// CheckUsername controlla che l'Username passato sia una stringa conforme alle specifiche dichiarate
// La funzione ritorna true quando l'Username passato è valido,false altrimenti.
func CheckUsername(username string) bool {

	// L'Username passato è composto solo da spazi bianchi quindi non è valido
	if strings.TrimSpace(username) == "" {
		return false
	} else {
		len := len(username)
		if len > 30 || len < 3 {
			return false
		} else {
			return true
		}
	}
}
