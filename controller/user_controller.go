package controller

import (
	"github.com/gorilla/mux"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"demo/service"
	"net/http"
)

type UserController struct {

}

var G_usc *UserController
func init(){
	G_usc = new (UserController)
}

type NameReq struct {
	Name string `json:"name"`
}
type UpdateRelationshipReq struct {
	State string `json:"state"`
}

func (usc *UserController)ListAllUsers(w http.ResponseWriter, r *http.Request) {
	ums := service.G_uss.ListAllUsers()
	json.NewEncoder(w).Encode(ums)
}

func (usc *UserController)CreateUser(w http.ResponseWriter, r *http.Request) {
	body,_ := ioutil.ReadAll(r.Body)
	var nameReq NameReq
	json.Unmarshal(body,&nameReq)

	if len(nameReq.Name) <= 0 {
		return
	}

	user,err := service.G_uss.CreateUser(nameReq.Name)
	if err != nil {
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (usc *UserController)ListAllRelationshipOfUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	suid := vars["user_id"]
	uid,err := strconv.ParseInt(suid,10,64)
	if err != nil{
		return
	}
	relationships := service.G_uss.ListAllRelationshipOfUser(uid)
	json.NewEncoder(w).Encode(relationships)
}

func (usc *UserController)UpdateRelationship(w http.ResponseWriter,r *http.Request) {
	vars := mux.Vars(r)
	suserId := vars["user_id"]
	sotherUserId := vars["other_user_id"]
	userId,err := strconv.ParseInt(suserId,10,64)
	if err != nil{
		return
	}
	otherUserId,err := strconv.ParseInt(sotherUserId,10,64)
	if err != nil{
		return
	}

	body,_ := ioutil.ReadAll(r.Body)
	var bodyReq UpdateRelationshipReq
	json.Unmarshal(body,&bodyReq)

	if len(bodyReq.State) <= 0 {
		return
	}

	relationship := service.G_uss.UpdateRelationship(userId,otherUserId,bodyReq.State)
	json.NewEncoder(w).Encode(relationship)
}

