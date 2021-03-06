package main

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Id      int
	Name    string
	Profile *Profile       `orm:"rel(one)"`      // OneToOne relation
	Posts   []*ArticlePost `orm:"reverse(many)"` // 设置一对多的反向关系
}

type Profile struct {
	Id   int
	Age  int16
	User *User `orm:"reverse(one)"` // 设置一对一反向关系(可选)
}

type ArticlePost struct {
	Id    int
	Title string
	User  *User         `orm:"rel(fk)"` //设置一对多关系
	Tags  []*ArticleTag `orm:"rel(m2m)"`
}

type ArticleTag struct {
	Id    int
	Name  string
	Posts []*ArticlePost `orm:"reverse(many)"` //设置多对多反向关系
}

func init() {
	// 需要在init中注册定义的model
	// 多对多反向关系 beego@v2.0.1 有 bug，需要升级到最新版
	orm.RegisterModel(new(User), new(ArticlePost), new(Profile), new(ArticleTag))
	// 自动创建数据表
	err := orm.RunSyncdb("default", true, true)
	if err != nil {
		fmt.Println(err)
		return
	}
}
