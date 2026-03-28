package utils

import (
	"errors"
	"strings"
)

// Costanti //
const Bearer_Authorization = "Bearer"

// Errori

// DoLogin //
var ErrUserNotRegistered error = errors.New("User not registered yet.")
var ErrUserAlreadyExists error = errors.New("User already exists.")
var ErrInvalidCredentials = errors.New("invalid username or password")

// CheckAuthorization //
var ErrAuthorizationNotSpecified error = errors.New("Authorization non specificata nell'header.")
var ErrBearerTokenNotSpecifiedWell error = errors.New("Bearer Token non specificato correttamente nel campo Authorization dell'header.")
var ErrUnauthorized error = errors.New("L'id passato nel campo Authorization non è autorizzato ad effettuare l'operazione.")

// AddUser //

// FollowUser //
var ErrFollowerAlreadyAdded error = errors.New("L'user ")

// BanUser //
var ErrUserAlreadyBanned error = errors.New("L'user è già bannato.")

// CheckUserExistence //
var ErrUserDoesNotExist error = errors.New("L'utente cercato non esiste: non ha ancora creato un profilo o non è ancora registrato.")

// LikePhoto //
var ErrLikeAlreadyExists error = errors.New("Il like con id specificato è già presente sulla foto.")

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

func CheckPassword(password string) bool {
	// La password passata è composta solo da spazi bianchi quindi non è valida
	if strings.TrimSpace(password) == "" {
		return false
	} else {
		len := len(password)
		if len < 8 {
			return false
		} else {
			return true
		}
	}

}
