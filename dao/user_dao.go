package dao

import (
	_ "github.com/lib/pq"
	"strconv"
	"fmt"
	"errors"
)

const (
	LIKED    = iota
	DISLIKED
	MATCHED
)

/**
 * define user model
 */
type User struct {
	Id   int64
	Name string
}

/**
 *define user service interface
 */
type UserDaoer interface {
	Register(userName string) (int64, error)
	GetUser(userIds []int64) (map[int64]User, error)
	ListAllUser() ([]User, error)
	UpdateRelationship(user_id, other_user_id int64, state string) (bool, error)
	GetRelationshipsOfUser(user_id int64) (map[int64]int8, error)
	GetRelationship(user_id, other_user_id int64) (int8, error)
}

type UserDaoerImpl struct {
}

var G_usd *UserDaoerImpl

func init() {
	G_usd = new(UserDaoerImpl)
}

/**
 *注册新用户
 */
func (usd *UserDaoerImpl) Register(userName string) (int64, error) {
	var uid int64
	err := G_db.QueryRow("INSERT INTO users(name) VALUES($1) RETURNING id",
		username).Scan(&uid)
	if err != nil {
		return -1, err
	}
	return uid, nil
}

/**
 *批量获取用户信息
 */
func (usd *UserDaoerImpl) GetUser(userIds []int64) (map[int64]User, error) {
	str := ``
	for _, val := range userIds {
		str = str + strconv.FormatInt(val, 10) + ","
	}
	if (len(str) > 1) {
		str = str[0:len(str)-1]
	} else {
		return nil, errors.New("invalid args : userIds")
	}

	sql := fmt.Sprintf("SELECT id,name FROM users WHERE id IN (%s)", str)

	rows, err := G_db.Query(sql)
	if err != nil {
		return nil, err
	}
	usermap := make(map[int64]User)
	for rows.Next() {
		row := new(User)
		rows.Scan(&row.Id, &row.Name)
		usermap[row.Id] = *row
	}
	return usermap, nil
}

/**
 *获取DB中的所有用户
 */
func (usd *UserDaoerImpl) ListAllUser() ([]User, error) {
	stmt, err := G_db.Prepare("SELECT id,name FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	users := make([]User, 0)
	for rows.Next() {
		row := new(User)
		rows.Scan(&row.Id, &row.Name)
		users = append(users, *row)
	}
	return users, nil
}

/**
 *更新两个用户之间的关系
 *state取值仅限liked或disliked
 */
func (usd *UserDaoerImpl) UpdateRelationship(user_id, other_user_id int64, state string) (bool, error) {
	_, err := G_db.Exec(`INSERT INTO relationships(user_id,other_user_id,state) VALUES($1,$2,$3)
	ON CONFLICT(user_id,other_user_id) DO UPDATE SET state=$3`, user_id, other_user_id, state)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

/**
 *获取与指定用户有关系的所有相关用户及关系状态
 *注：状态包包括（liked,matched,disliked）
 */
func (usd *UserDaoerImpl) GetRelationshipsOfUser(userId int64) (map[int64]int8, error) {
	rows, err := G_db.Query(`SELECT r1.other_user_id,r1.state,r2.state
									FROM relationships AS r1 LEFT JOIN relationships AS r2
									ON r1.user_id=r2.other_user_id AND r1.other_user_id=r2.user_id
									WHERE r1.user_id=$1`, userId)
	if err != nil {
		return nil, err
	}

	m := make(map[int64]int8)
	for rows.Next() {
		var u int64
		var state1, state2 string
		rows.Scan(&u, &state1, &state2)
		m[u] = getIntState(state1, state2)
	}
	return m, nil
}

/**
 *获取一个用户对另一个用户的关系
 */
func (usd *UserDaoerImpl) GetRelationship(userId, otherUserID int64) (int8, error) {
	var state1, state2 string
	err := G_db.QueryRow(`SELECT r1.state,r2.state
								FROM relationships AS r1 LEFT JOIN relationships AS r2
								ON r1.user_id=r2.other_user_id AND r1.other_user_id=r2.user_id
								WHERE r1.user_id=$1 AND r1.other_user_id=$2`, userId, otherUserID).Scan(&state1, &state2)
	if err != nil {
		return 0, err
	}
	return getIntState(state1, state2), nil
}

/**
  *根据user和other_user的双边关系，判定user对other_user的关系
  *注：只适用于state1为liked或disliked的场景
 */
func getIntState(state1, state2 string) int8 {
	var state int8
	if state1 == "disliked" {
		state = DISLIKED
	} else if state1 == "liked" && state2 == "liked" {
		state = MATCHED
	} else {
		state = LIKED
	}
	return state
}
