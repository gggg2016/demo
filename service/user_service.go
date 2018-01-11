package service

import (
	"strconv"
	"demo/dao"
	"demo/util"
)

type UserModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type RelationShip struct {
	User_id string `json:"user_id"`
	State   string `json:"state"`
	Type    string `json:"type"`
}

type UserService struct {
}

var G_uss *UserService

func init() {
	G_uss = new(UserService)
}

func (uss *UserService) ListAllUsers() ([]UserModel, error) {
	users, err := dao.G_usd.ListAllUser()
	// len == 0 不应该返回错误吧？
	if err != nil || len(users) <= 0 {
		return nil, err
	}

	ums := make([]UserModel, 0)
	for _, user := range users {
		ums = append(ums, transformUserToUserModel(user))
	}
	return ums, nil
}

func transformUserToUserModel(u dao.User) UserModel {
	return UserModel{strconv.FormatInt(u.Id, 10), u.Name, "user"}
}

func (uss *UserService) CreateUser(name string) (UserModel, error) {
	id, err := dao.G_usd.Register(name)
	if err != nil {
		return UserModel{}, err
	}

	return UserModel{strconv.FormatInt(id, 10), name, "user"}, nil
}

func (uss *UserService) ListAllRelationshipOfUser(userId int64) ([]RelationShip, error) {
	//liked, matched included
	likedUserIds, err := dao.G_usd.ListLikedUser(userId)
	if err != nil {
		return nil, err
	}
	//matched
	matchedUserIds, err := dao.G_usd.ListMatchedUser(userId)
	if err != nil {
		return nil, err
	}
	//disliked
	dislikedUserIds, err := dao.G_usd.ListDislikedUser(userId)
	if err != nil {
		return nil, err
	}
	//singel-liked, matched not included
	singleLikedUserIds := util.Sub(likedUserIds, matchedUserIds)

	relationships := make([]RelationShip, 0)
	for _, u := range singleLikedUserIds {
		relationships = append(relationships, RelationShip{strconv.FormatInt(u, 10), "liked", "relationship"})
	}
	for _, u := range matchedUserIds {
		relationships = append(relationships, RelationShip{strconv.FormatInt(u, 10), "matched", "relationship"})
	}
	for _, u := range dislikedUserIds {
		relationships = append(relationships, RelationShip{strconv.FormatInt(u, 10), "disliked", "relationship"})
	}
	return relationships, nil
}

func (uss *UserService) UpdateRelationship(userId, otherUserId int64, state string) (RelationShip, error) {
	var st int8
	switch state {
	case "liked":
		st = dao.LIKE
	case "disliked":
		st = dao.DISLIKE
	}

	// 为啥不使用定义好的dao.LIKE dao.DISLIKE
	//update relationship state
	if st == 1 || st == -1 {
		dao.G_usd.UpdateRelationship(userId, otherUserId, st)
	}

	//query result state
	var resState string
	r1, r2, err := dao.G_usd.GetRelationship(userId, otherUserId)

	if err != nil {
		return RelationShip{}, err
	}

	switch {
	case r1 == -1:
		resState = "disliked"
	case r1 == 1 && r2 != 1:
		resState = "liked"
	case r1 == 1 && r2 == 1:
		resState = "matched"
	}

	return RelationShip{strconv.FormatInt(otherUserId, 10), resState, "relationship"}, nil
}
