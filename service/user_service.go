package service

import (
	"strconv"
	"demo/dao"
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
	if err != nil {
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

func stateString(state int8) string {
	switch state {
	case dao.LIKED:
		return "liked"
	case dao.MATCHED:
		return "matched"
	default:
		return "disliked"
	}
}
func (uss *UserService) CreateUser(name string) (UserModel, error) {
	id, err := dao.G_usd.Register(name)
	if err != nil {
		return UserModel{}, err
	}

	return UserModel{strconv.FormatInt(id, 10), name, "user"}, nil
}

func (uss *UserService) ListAllRelationshipOfUser(userId int64) ([]RelationShip, error) {
	m, err := dao.G_usd.GetRelationshipsOfUser(userId)
	if err != nil {
		return nil, err
	}

	relationships := make([]RelationShip, 0)
	for userId, state := range m {
		relationships = append(relationships, RelationShip{User_id: strconv.FormatInt(userId, 10), State: stateString(state), Type: "relationship"})
	}
	return relationships, nil
}

func (uss *UserService) UpdateRelationship(userId, otherUserId int64, state string) (RelationShip, error) {
	//update relationship state
	dao.G_usd.UpdateRelationship(userId, otherUserId, state)

	//query result state
	r, err := dao.G_usd.GetRelationship(userId, otherUserId)

	if err != nil {
		return RelationShip{}, err
	}

	return RelationShip{strconv.FormatInt(otherUserId, 10), stateString(r), "relationship"}, nil
}
