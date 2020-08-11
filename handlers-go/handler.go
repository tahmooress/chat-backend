package handlers

import (
	dbf "app/db-go"
	"app/util-go"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

//LoginHandler is handler for login route
func LoginHandler(env *dbf.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enbaleCors(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		var user util.NewUser
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &user)
		fmt.Println(user)
		var res util.ResponseResult
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		var temp util.NewUser
		err = env.DB.QueryRow(dbf.CheckExist, user.Email).Scan(&temp.Email, &temp.Password)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(temp.Password), []byte(user.Password))
		if err != nil {
			res.Error = "password is wrong!"
			json.NewEncoder(w).Encode(res)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": user.Email,
			"exp":  time.Now().Add(time.Hour * time.Duration(1)).Unix(),
			"iat":  time.Now().Unix(),
		})
		tokenString, err := token.SignedString([]byte(jwtKey))
		if err != nil {
			res.Error = "error occurd while generating token"
			json.NewEncoder(w).Encode(res)
			return
		}
		res.Result = tokenString
		user.Token = tokenString
		json.NewEncoder(w).Encode(res)
		fmt.Println("login", user.Token, res.Result)
	})
}

const jwtKey = "my_secret_key"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//RegisterHandler is a handler for register route
func RegisterHandler(env *dbf.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enbaleCors(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		var user util.NewUser
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &user)
		var res util.ResponseResult
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		cost := 5
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), cost)
		if err != nil {
			res.Error = "Error while Creating user, Try again"
			json.NewEncoder(w).Encode(res)
			return
		}
		_, err = env.DB.Exec(dbf.RegisterNewUserStatement, user.Username, user.Email, string(hash))
		if err != nil {
			res.Error = err.Error()
			res.Result = "user already exist"
			json.NewEncoder(w).Encode(res)
			return
		}
		res.Result = "Registered succesfully"
		json.NewEncoder(w).Encode(res)
	})
}

//WsHandler is handler for chat and websocket
func WsHandler(env *dbf.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enbaleCors(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("there  was an error")
				}
				return []byte(jwtKey), nil
			})
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
			if token.Valid {
				con, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					log.Println(err)
				}
				user := util.New()
				user.Con = con
				user.Clients = util.Client
				if err = websocket.ReadJSON(con, user); err != nil {
					log.Println(err)
					return
				}
				go user.Reader(env)
				go user.Writer()
			}
		} else {
			fmt.Fprintf(w, "NOT AUTHORIZED")
		}
	})
}

func enbaleCors(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
