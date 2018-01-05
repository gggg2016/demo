package dao

import(
	"testing"
)

func setUpTest() {

}
func Test_Db(t *testing.T){
	usd := new (UserDao)

	r,err := usd.Register("fan12")
	t.Error(r)
	t.Log(r,"\n",err)
}

func Test_GetUser(t *testing.T){
	usd := new (UserDao)
	r := usd.GetUser([]int64{1,5,7})
	t.Error("result: ",r)
}

func Test_ListAllUser(t *testing.T){
	usd := new (UserDao)
	r:= usd.ListAllUser()
	t.Error("result size: ",len(r))
	t.Log("result: ",r)
}

func Test_UpdateRelationship(t *testing.T){
	usd := new (UserDao)
	r:=usd.UpdateRelationship(2,8,1)
	t.Error(r)
}