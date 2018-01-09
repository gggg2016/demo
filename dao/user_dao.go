package dao

import (
	_ "github.com/lib/pq"
	"strconv"
	"fmt"
	"errors"
)

const (
	LIKE    = 1
	DISLIKE = -1
	NONE    = 0
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
	UpdateRelationship(first_user_id, second_user_id int64, state int8) (bool, error)
	ListLikedUser(userId int64) ([]int64, error)
	ListDislikedUser(userId int64) ([]int64, error)
	ListMatchedUser(userId int64) ([]int64, error)
	GetRelationship(first_user_id, second_user_id int64) (int8, int8, error)
}

type UserDaoerImpl struct {
}

var G_usd *UserDaoerImpl

func init() {
	G_usd = new(UserDaoerImpl)
}

func (usd *UserDaoerImpl) CreateTableRelationships() error {
	_, err := G_db.Exec(`CREATE TABLE public.relationships
								(
									id bigint NOT NULL DEFAULT nextval('relationship_id_seq'::regclass),
									first_user_id bigint NOT NULL,
									second_user_id bigint NOT NULL,
									state smallint NOT NULL,
									CONSTRAINT relationship_pkey PRIMARY KEY (id),
									CONSTRAINT user_pair UNIQUE (first_user_id, second_user_id)
								)`)
	return err

}

func (usd *UserDaoerImpl) CreateTableUsers() error {
	_, err := G_db.Exec(`CREATE TABLE public.users
								(
									id bigint NOT NULL DEFAULT nextval('user_id_seq'::regclass),
									name character varying(10) COLLATE pg_catalog."default" NOT NULL,
									create_time timestamp with time zone,
									update_time timestamp with time zone,
									CONSTRAINT user_pkey PRIMARY KEY (id)
								)`)
	return err
}

func (usd *UserDaoerImpl) Register(userName string) (int64, error) {
	var uid int64
	err := G_db.QueryRow("INSERT INTO public.users(name,create_time,update_time) VALUES($1,now(),now()) RETURNING id",
		username).Scan(&uid)
	if err != nil {
		return -1, err
	}
	return uid, nil
}

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

	sql := fmt.Sprintf("SELECT id,name FROM public.users WHERE id IN (%s)", str)

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

func (usd *UserDaoerImpl) UpdateRelationship(first_user_id, second_user_id int64, state int8) (bool, error) {
	_, err := G_db.Exec(`INSERT INTO relationships(first_user_id,second_user_id,state) VALUES($1,$2,$3)
	ON CONFLICT(first_user_id,second_user_id) DO UPDATE SET state=$3`, first_user_id, second_user_id, state)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (usd *UserDaoerImpl) ListLikedUser(userId int64) ([]int64, error) {
	rows, err := G_db.Query(`SELECT second_user_id FROM relationships
		WHERE first_user_id=$1 AND state=1`, userId)

	if err != nil {
		return nil, err
	}

	userIds := make([]int64, 0)
	for rows.Next() {
		var userId int64
		rows.Scan(&userId)
		userIds = append(userIds, userId)
	}
	return userIds, nil
}
func (usd *UserDaoerImpl) ListDislikedUser(userId int64) ([]int64, error) {
	rows, err := G_db.Query(`SELECT second_user_id FROM relationships
		WHERE first_user_id=$1 AND state=-1`, userId)

	if err != nil {
		return nil, err
	}

	userIds := make([]int64, 0)
	for rows.Next() {
		var userId int64
		rows.Scan(&userId)
		userIds = append(userIds, userId)
	}
	return userIds, nil
}

func (usd *UserDaoerImpl) ListMatchedUser(userId int64) ([]int64, error) {
	rows, err := G_db.Query(`SELECT r1.second_user_id FROM relationships AS r1,relationships AS r2
		WHERE r1.first_user_id=$1 AND r1.state=1 AND r2.state=1 AND
		r1.first_user_id=r2.second_user_id AND r1.second_user_id=r2.first_user_id`, userId)
	if err != nil {
		return nil, err
	}

	userIds := make([]int64, 0)
	for rows.Next() {
		var userId int64
		rows.Scan(&userId)
		userIds = append(userIds, userId)
	}
	return userIds, nil
}

func (usd *UserDaoerImpl) GetRelationship(first_user_id, second_user_id int64) (int8, int8, error) {
	rows, err := G_db.Query(`SELECT first_user_id,second_user_id,state FROM relationships
		WHERE (first_user_id=$1 AND second_user_id=$2) OR (first_user_id=$2 AND second_user_id=$1)`, first_user_id, second_user_id)
	if err != nil {
		return 0, 0, err
	}

	var r1, r2 int8
	for rows.Next() {
		var u1, u2 int64
		var state int8
		rows.Scan(&u1, &u2, &state)
		switch {
		case u1 == first_user_id:
			r1 = state
		case u1 == second_user_id:
			r2 = state
		}
	}
	return r1, r2, nil
}
