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
	Msg string `json:"msg"`
}

/**
 *正常返回json数据
 */
func httpWriteJsonResult(w http.ResponseWriter, obj interface{}) {
	resp, err := json.Marshal(obj)
	if err != nil {
		httpWriteErrResult(w)
	} else {
		w.Header().Set("Content-type", "application/json;charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

/**
 *服务内部错误
 */
func httpWriteErrResult(w http.ResponseWriter) {
	result := ErrorResult{Msg: "internal error"}
	resp, _ := json.Marshal(result)
	w.Header().Set("Content-type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(resp)
}

/**
 *参数错误
 */
func httpWriteIllegalArgResul(w http.ResponseWriter) {
	result := ErrorResult{Msg: "illegal argument error"}
	resp, _ := json.Marshal(result)
	w.Header().Set("Content-type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(resp)
}

/**
 *列出所有用户
 */
func (usc *UserController) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	ums, err := service.G_uss.ListAllUsers()
	if err != nil {
		httpWriteErrResult(w)
		return
	}
	httpWriteJsonResult(w, ums)
}

/**
 *创建新用户
 */
func (usc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var nameReq NameReq
	json.Unmarshal(body, &nameReq)

	if len(nameReq.Name) <= 0 {
		httpWriteIllegalArgResul(w)
		return
	}

	user, err := service.G_uss.CreateUser(nameReq.Name)
	if err != nil {
		httpWriteErrResult(w)
		return
	}

	httpWriteJsonResult(w, user)
}

/**
 *获取指定用户的所有关系
 */
func (usc *UserController) ListAllRelationshipOfUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	suid := vars["user_id"]
	uid, err := strconv.ParseInt(suid, 10, 64)
	if err != nil {
		httpWriteIllegalArgResul(w)
		return
	}
	relationships, err := service.G_uss.ListAllRelationshipOfUser(uid)
	if err != nil {
		httpWriteErrResult(w)
		return
	}
	httpWriteJsonResult(w, relationships)
}

/**
 *更新指定用户之间的关系(左滑/右滑)
 */
func (usc *UserController) UpdateRelationship(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	suserId := vars["user_id"]
	sotherUserId := vars["other_user_id"]
	userId, err := strconv.ParseInt(suserId, 10, 64)
	if err != nil {
		httpWriteIllegalArgResul(w)
		return
	}
	otherUserId, err := strconv.ParseInt(sotherUserId, 10, 64)
	if err != nil {
		httpWriteIllegalArgResul(w)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var bodyReq UpdateRelationshipReq
	json.Unmarshal(body, &bodyReq)

	if bodyReq.State != "liked" && bodyReq.State != "disliked" {
		httpWriteIllegalArgResul(w)
		return
	}

	relationship, err := service.G_uss.UpdateRelationship(userId, otherUserId, bodyReq.State)
	if err != nil {
		httpWriteErrResult(w)
	} else {
		httpWriteJsonResult(w, relationship)
	}
}
