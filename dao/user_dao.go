package dao

import (
	_ "github.com/lib/pq"
	"strconv"
	"fmt"
)

const(
	LIKE = 1
	DISLIKE = -1
	NONE = 0
)
/**
 * define user model
 */
type User struct{
	Id int64
	Name string
}
/**
 *define user service interface
 */
type IuserDao interface{
	Register(userName string)(int64,error)
	GetUser(userIds []int64)map[int64]User
	ListAllUser()[]User
	UpdateRelationship(first_user_id,second_user_id int64,state int8)bool
	ListLikedUser(userId int64)[]int64
	ListDislikedUser(userId int64)[]int64
	ListMatchedUser(userId int64)[]int64
	GetRelationship(first_user_id,second_user_id int64)int8
}

type UserDao struct{
}

var G_usd *UserDao
func init(){
	G_usd = new (UserDao)
}

func (usd *UserDao) Register(userName string)(int64,error){
	var uid int64
	err := G_db.QueryRow("INSERT INTO public.user(name,create_time,update_time) VALUES($1,now(),now()) RETURNING id",
	      				username).Scan(&uid)
	if err != nil {
		return -1,err
	}	
	return uid, nil
}

func (usd *UserDao) GetUser(userIds []int64)map[int64]User{
	str := ``
	for _,val := range userIds{
		str = str + strconv.FormatInt(val,10) + ","
	}
	if(len(str)>1){
		str = str[0:len(str)-1]
	}

	sql := fmt.Sprintf("SELECT id,name FROM public.user WHERE id IN (%s)",str)

	rows,err := G_db.Query(sql)
	if err != nil{
		return nil
	}
	usermap := make(map[int64]User)
	for rows.Next(){
		row := new (User)
		rows.Scan(&row.Id,&row.Name)
		usermap[row.Id]=*row
	}
	return usermap
}

func (usd *UserDao) ListAllUser()[]User{
	stmt,err := G_db.Prepare("SELECT id,name FROM public.user ORDER BY id")
	if err!=nil{
		return nil
	}
	rows,err := stmt.Query()
	if err!=nil{
		return nil
	}
	users := make([]User,0)
	for rows.Next(){
		row := new (User)
		rows.Scan(&row.Id,&row.Name)
		users = append(users,*row)
	}
	return users
}

func (usd *UserDao)UpdateRelationship(first_user_id,second_user_id int64,state int8)bool{
	_,err := G_db.Exec(`INSERT INTO public.relationship(first_user_id,second_user_id,state) VALUES($1,$2,$3)
	ON CONFLICT(first_user_id,second_user_id) DO UPDATE SET state=$3`,first_user_id,second_user_id,state)
    return err == nil
} 

func (usd *UserDao)ListLikedUser(userId int64)[]int64{
	rows,err := G_db.Query(`SELECT second_user_id FROM relationship 
		WHERE first_user_id=$1 AND state=1`,userId)

	if err != nil{
		return nil
	}

	userIds := make([]int64, 0)
	for rows.Next(){
		var userId int64
		rows.Scan(&userId)
		userIds = append(userIds,userId)
	}
	return userIds
}
func (usd *UserDao)ListDislikedUser(userId int64)[]int64{
	rows,err := G_db.Query(`SELECT second_user_id FROM relationship 
		WHERE first_user_id=$1 AND state=-1`,userId)

	if err != nil{
		return nil
	}

	userIds := make([]int64, 0)
	for rows.Next(){
		var userId int64
		rows.Scan(&userId)
		userIds = append(userIds,userId)
	}
	return userIds
}

func (usd *UserDao)ListMatchedUser(userId int64)[]int64{
	rows,err := G_db.Query(`SELECT r1.second_user_id FROM relationship AS r1,relationship AS r2
		WHERE r1.first_user_id=$1 AND r1.state=1 AND r2.state=1 AND
		r1.first_user_id=r2.second_user_id AND r1.second_user_id=r2.first_user_id`,userId)
	if err != nil {
		return nil
	}

	userIds := make([]int64,0)
	for rows.Next() {
		var userId int64
		rows.Scan(&userId)
		userIds = append(userIds,userId)
	}
	return userIds
}

func (usd *UserDao)GetRelationship(first_user_id,second_user_id int64)(int8,int8){
	rows,err := G_db.Query(`SELECT first_user_id,second_user_id,state FROM relationship
		WHERE (first_user_id=$1 AND second_user_id=$2) OR (first_user_id=$2 AND second_user_id=$1)`,first_user_id,second_user_id)
	if err != nil {
		return 0,0
	}

	var r1,r2 int8
	for rows.Next() {
		var u1,u2 int64
		var state int8
		rows.Scan(&u1,&u2,&state)
		switch{
		case u1 == first_user_id:
			r1 = state
		case u1 ==second_user_id:
			r2 = state
		}
	}
	return r1,r2
}




