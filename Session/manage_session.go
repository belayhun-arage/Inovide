package session

import (
	SessionRepo "github.com/Projects/Inovide/Session/Repository"
	entity "github.com/Projects/Inovide/models"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var jwtKey = []byte("This is my Secret Key For Encryption ")

type Cookiehandler struct {
	Sessionrepo *SessionRepo.SessionRepository
}

func NewCookieHandler() *Cookiehandler {
	return &Cookiehandler{}
}
func (coockiehandler *Cookiehandler) SaveSession(writer http.ResponseWriter, session *entity.Session) bool {
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	Succesfull := false
	expirationTime := time.Now().Add(12 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	usercookie := &entity.Claim{
		Username: session.Username,
		Id:       session.Userid,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, usercookie)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		writer.WriteHeader(http.StatusInternalServerError)
		return Succesfull
	}
	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(writer, &http.Cookie{
		Name:     "InovideToken",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	})
	erro := coockiehandler.Sessionrepo.CreateSession(session)
	if erro != nil {
		Succesfull = false
	}
	Succesfull = true
	return Succesfull
}

func (sessionhandler *Cookiehandler) DeleteSession(writer http.ResponseWriter, request *http.Request) bool {
	sucess := false
	c, err := request.Cookie("InovideToken")
	if err != nil {
		return sucess
	}
	tknStr := c.Value
	claims := &entity.Claim{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return sucess
		}
		return sucess
	}
	if !tkn.Valid {
		return sucess
	}
	session := &entity.Session{
		Userid:   claims.Id,
		Username: claims.Username,
	}
	errors := sessionhandler.Sessionrepo.DeleteSession(session)

	if errors != nil {
		sucess = false
		return sucess
	}
	cookie := http.Cookie{
		Name:    "InovideToken",
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
		Value:   "",
	}
	sucess = true
	http.SetCookie(writer, &cookie)
	return sucess

}
func (sessionhandler *Cookiehandler) Valid(request *http.Request) (int, string, bool) {
	sucess := false
	c, err := request.Cookie("InovideToken")
	if err != nil {
		return -1, "", sucess
	}
	tknStr := c.Value
	claims := &entity.Claim{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return -1, "", sucess
		}
		return -1, "", sucess

	}
	if !tkn.Valid {
		return -1, "", sucess
	}
	return claims.Id, claims.Username, sucess
}
