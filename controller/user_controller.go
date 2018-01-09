package controller

import (
	"github.com/gorilla/mux"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"demo/service"
	"net/http"
	"errors"
)

type UserController struct {
}

var G_usc *UserController

func init() {
	G_usc = new(UserController)
}

type NameReq struct {
	Name string `json:"name"`
}
type UpdateRelationshipReq struct {
	State string `json:"state"`
}

type ErrorResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func httpWriteResultAsJson(w http.ResponseWriter, obj interface{}) {
	resp, err := json.Marshal(obj)
	if err != nil {
		httpWriteResultAsErr(w, err)
	} else {
		w.Write(resp)
	}
}

func httpWriteResultAsErr(w http.ResponseWriter, err error) {
	result := ErrorResult{Code: 500, Msg: err.Error()}
	resp, _ := json.Marshal(result)
	w.Write(resp)
}

func (usc *UserController) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	ums, err := service.G_uss.ListAllUsers()

	if err != nil {
		httpWriteResultAsErr(w, err)
		return
	}

	httpWriteResultAsJson(w, ums)
}

func (usc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var nameReq NameReq
	json.Unmarshal(body, &nameReq)

	if len(nameReq.Name) <= 0 {
		httpWriteResultAsErr(w, errors.New("arguments, name, is invalid."))
		return
	}

	user, err := service.G_uss.CreateUser(nameReq.Name)
	if err != nil {
		httpWriteResultAsErr(w, err)
		return
	}

	httpWriteResultAsJson(w, user)
}

func (usc *UserController) ListAllRelationshipOfUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	suid := vars["user_id"]
	uid, err := strconv.ParseInt(suid, 10, 64)
	if err != nil {
		httpWriteResultAsErr(w, err)
		return
	}
	relationships, err := service.G_uss.ListAllRelationshipOfUser(uid)
	if err != nil {
		httpWriteResultAsErr(w, err)
		return
	}
	httpWriteResultAsJson(w, relationships)
}

func (usc *UserController) UpdateRelationship(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	suserId := vars["user_id"]
	sotherUserId := vars["other_user_id"]
	userId, err := strconv.ParseInt(suserId, 10, 64)
	if err != nil {
		return
	}
	otherUserId, err := strconv.ParseInt(sotherUserId, 10, 64)
	if err != nil {
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var bodyReq UpdateRelationshipReq
	json.Unmarshal(body, &bodyReq)

	if len(bodyReq.State) <= 0 {
		httpWriteResultAsErr(w, errors.New("invalid argument : state"))
		return
	}

	relationship, err := service.G_uss.UpdateRelationship(userId, otherUserId, bodyReq.State)
	if err != nil {
		httpWriteResultAsErr(w, err)
	} else {
		httpWriteResultAsJson(w, relationship)
	}
}
