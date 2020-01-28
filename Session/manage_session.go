package session

import (
	"fmt"
	"log"
	"net/http"
	"time"

	SessionRepo "github.com/Projects/Inovide/Session/Repository"
	entity "github.com/Projects/Inovide/models"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("This is my Secret Key For Encryption ")

type Cookiehandler struct {
	Sessionrepo *SessionRepo.SessionRepository
}

func NewCookieHandler(sessionrepo *SessionRepo.SessionRepository) *Cookiehandler {
	return &Cookiehandler{Sessionrepo: sessionrepo}
}
func (coockiehandler *Cookiehandler) SaveSession(writer http.ResponseWriter, session *entity.Session) bool {
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	Succesfull := false
	expirationTime := time.Now().Add(24 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time

	fmt.Print(session.Userid, session.Username, "\n\n\n\n\n")

	fmt.Println(session.Userid, "In  the session package ")
	usercookie := &entity.Claim{
		Username: session.Username,
		Id:       session.Userid,
		IsAdmin:  session.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
			// HttpOnly:  true,
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, usercookie)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return Succesfull
	}
	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	cookie := http.Cookie{
		Name:    "inovidetoken",
		Value:   tokenString,
		Expires: expirationTime,
		Domain:  "localhost:8080",
		Path:    "/",
	}
	http.SetCookie(writer, &cookie)
	log.Print(cookie)

	session.IsAdmin = false
	erro := coockiehandler.Sessionrepo.CreateSession(session)
	if erro != nil {
		Succesfull = false
	}
	Succesfull = true
	return Succesfull
}

func (sessionhandler *Cookiehandler) DeleteSession(writer http.ResponseWriter, request *http.Request) bool {
	sucess := false
	id, username, _ := sessionhandler.Valid(request)

	if id >= 0 {
		session := &entity.Session{
			Userid:   id,
			Username: username,
		}
		code := sessionhandler.Sessionrepo.DeleteSession(session)
		if code == 1 {
			sucess = false
			return sucess
		}
		cookie := http.Cookie{
			Name:    "inovidetoken",
			Domain:  "localhost:8080",
			Path:    "/",
			Expires: time.Now().Add(-10 * time.Second),
		}
		sucess = true
		http.SetCookie(writer, &cookie)
		return sucess
	} else {
		fmt.Println("Nor OOOOkk ")
		return false
	}

}
func (sessionhandler *Cookiehandler) Valid(request *http.Request) (int, string, bool) {
	sucess := false

	c, err := request.Cookie("inovidetoken")

	defer recover()
	// fmt.Println(c.Value)
	if err != nil {
		fmt.Println(err.Error())
	}
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
			fmt.Println("two")

			return -1, "", sucess
		}
		fmt.Println("three")

		return -1, "", sucess

	}
	if !tkn.Valid {
		fmt.Println("four")

		return -1, "", sucess
	}
	return claims.Id, claims.Username, sucess
}

func (sessionhandler *Cookiehandler) Authorize(request *http.Request) (IsAdmin bool, IsUser bool) {
	c, err := request.Cookie("inovidetoken")

	defer recover()
	// fmt.Println(c.Value)
	if err != nil {
		fmt.Println(err.Error())
	}
	if err != nil {
		IsAdmin = false
		IsUser = false
		return IsAdmin, IsUser
	}
	tknStr := c.Value
	claims := &entity.Claim{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Println("two")

			IsAdmin = false
			IsUser = false
			return IsAdmin, IsUser
		}
		fmt.Println("three")

		IsAdmin = false
		IsUser = false
		return IsAdmin, IsUser

	}
	if !tkn.Valid {
		fmt.Println("four")

		IsAdmin = false
		IsUser = false
		return IsAdmin, IsUser
	}

	IsUser = true
	if claims.IsAdmin {
		IsAdmin = true
	} else {
		IsAdmin = false
	}
	return IsAdmin, IsUser
}

func (sessionhandle *Cookiehandler) RandomToken() string {

	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString

}

func (sessionhandle *Cookiehandler) ValidateForm(tokenstring string) bool {

	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return false
	}
	return true
}
