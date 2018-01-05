package service

import (
	"strconv"
	"demo/dao"
	"demo/util"
)

type UserModel struct{
	Id string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type RelationShip struct{
	User_id string `json:"user_id"`
	State string `json:"state"`
	Type string `json:"type"`
}

type UserService struct{

}

var G_uss *UserService
func init(){
	G_uss = new (UserService)
}

func (uss *UserService)ListAllUsers()[]UserModel{
	users := dao.G_usd.ListAllUser()
	if users == nil {
		return nil
	}

	ums := make([]UserModel,0)
	for _,user := range users{
		ums = append(ums,transformUserToUserModel(user))
	}
	return ums
}

func transformUserToUserModel(u dao.User)UserModel{
	return UserModel{strconv.FormatInt(u.Id,10),u.Name,"user"}
}

func (uss *UserService)CreateUser(name string)(UserModel,error){
	id,err := dao.G_usd.Register(name)
	if err != nil {
		return UserModel{},err
	}

	return UserModel{strconv.FormatInt(id,10),name,"user"},nil
}

func (uss *UserService)ListAllRelationshipOfUser(userId int64)[]RelationShip{
	//liked, matched included
	likedUserIds := dao.G_usd.ListLikedUser(userId)
	//matched
	matchedUserIds := dao.G_usd.ListMatchedUser(userId)
	//disliked
	dislikedUserIds := dao.G_usd.ListDislikedUser(userId)
	//singel-liked, matched not included
	singleLikedUserIds := util.Sub(likedUserIds,matchedUserIds)

	relationships := make([]RelationShip, 0)
	for _,u := range singleLikedUserIds {
		relationships = append(relationships,RelationShip{strconv.FormatInt(u,10),"liked","relationship"})
	}
	for _,u := range matchedUserIds {
		relationships = append(relationships,RelationShip{strconv.FormatInt(u,10),"matched","relationship"})
	}
	for _,u := range dislikedUserIds {
		relationships = append(relationships,RelationShip{strconv.FormatInt(u,10),"disliked","relationship"})
	}
	return relationships
}

func (uss *UserService)UpdateRelationship(userId, otherUserId int64, state string)RelationShip{
	var st int8
	switch state {
	case "liked":
		st = dao.LIKE
	case "disliked":
		st = dao.DISLIKE
	}
	
	//update relationship state
	if st == 1 || st == -1{
		dao.G_usd.UpdateRelationship(userId,otherUserId,st)
	}

	//query result state
	var resState string
	r1,r2 := dao.G_usd.GetRelationship(userId,otherUserId)
	switch {
	case r1 == -1:
		resState = "disliked"
	case r1 == 1 && r2 != 1:
		resState = "liked"
	case r1 == 1 && r2 == 1:
		resState = "matched"
	}

	return RelationShip{strconv.FormatInt(otherUserId,10),resState,"relationship"}
}
