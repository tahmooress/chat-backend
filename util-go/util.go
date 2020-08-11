package util

import (
	dbf "app/db-go"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

//User is the model of any user who connected
type User struct {
	Email   Email
	Token   string
	Con     *websocket.Conn
	Ch      chan Pocket
	Clients map[Email]chan Pocket
}

//Pocket is the wraper of any user message which comminicate in the server
type Pocket struct {
	Message string `json:"message"`
	From    Email  `json:"from"`
}

//Email is the identifier of users
type Email string

//RecieveP is the type for modeling the messages from client over websocket which are recieving
type RecieveP struct {
	To      Email  `json:"to"`
	Message string `json:"message"`
}

//Writer is the method that users wirte to websocket mesasges which are reciving from other users
func (u User) Writer() {
	defer func() { fmt.Println("client disconnected") }()
	for {
		select {
		case p := <-u.Ch:
			if err := websocket.WriteJSON(u.Con, p); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

//Reader read from websocket and send the message to the targeted users chanel
func (u User) Reader(env *dbf.Env) {
	defer func() {
		delete(u.Clients, u.Email)
		close(u.Ch)
	}()
	var r RecieveP
	var p Pocket
	p.From = u.Email
	p.Message = r.Message
	for {
		err := websocket.ReadJSON(u.Con, &r)
		if err != nil {
			log.Println(err)
			return
		}
		if _, ok := u.Clients[r.To]; ok {
			u.Clients[r.To] <- p
		} else {
			go func(env *dbf.Env) {
				//inja aval bayad query bezanim 2ta ke ba email , id ro begirim az data base va bezarim too Exec e paieen

				_, err := env.DB.Exec(dbf.NewMessageStatement, p.From, r.To, p.Message)
				if err != nil {
					fmt.Println(err)
					return
				}
			}(env)
		}

	}
}

//New create a new user with specified channel and return a pointer to that new created user
func New() *User {
	var u User
	u.Ch = make(chan Pocket)
	return &u
}

//NewUser is type for modeling any users that register to app
type NewUser struct {
	Username string `json:"username"`
	Email    Email  `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

//ResponseResult is type for modeling result of response to client for registeration
type ResponseResult struct {
	Error  string `json:"er"`
	Result string `json:"result"`
}

//Client is the list of all the users who are now connected
var Client = make(map[Email]chan Pocket)
