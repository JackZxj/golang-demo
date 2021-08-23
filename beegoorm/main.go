package main

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/orm_test?charset=utf8")
}

func main() {

	// 默认使用 default，你可以指定为其他数据库
	o := orm.NewOrm()

	profile1 := new(Profile)
	profile1.Age = 10

	profile2 := new(Profile)
	profile2.Age = 20

	user1 := new(User)
	user1.Profile = profile1
	user1.Name = "user1"

	user2 := new(User)
	user2.Profile = profile2
	user2.Name = "user2"

	tag1 := new(Tag)
	tag1.Name = "tag1"

	tag2 := new(Tag)
	tag2.Name = "tag2"

	post1 := new(Post)
	post1.Title = "hello1"
	post1.User = user1
	// post1.Tags = []*Tag{tag1, tag2}
	// post1.Tags = append(post1.Tags, tag1, tag2)

	post2 := new(Post)
	post2.Title = "hello2"
	post2.User = user2
	// post2.Tags = []*Tag{tag1}
	// post1.Tags = append(post1.Tags, tag1)

	o.Insert(profile1)
	o.Insert(profile2)
	o.Insert(user1)
	o.Insert(user2)
	o.Insert(tag1)
	o.Insert(tag2)
	o.Insert(post1)
	o.Insert(post2)

	// 保存多对多关系
	m2m1 := o.QueryM2M(post1, "Tags")
	m2m1.Add(tag1, tag2)

	m2m2 := o.QueryM2M(post2, "Tags")
	m2m2.Add(tag1)
	fmt.Println()

	{
		// 外键（一对多）是加载的
		// 多对多是不加载的
		var posts []*Post
		num, err := o.QueryTable("post").Filter("Tags__Tag__Name", "tag1").All(&posts)
		if err == nil {
			fmt.Printf("%d posts read\n", num)
			for _, post := range posts {
				v1, _ := json.Marshal(*post)
				fmt.Println("默认值：", string(v1))
				// 默认不载入关系字段，需要手动载入
				// https://beego.me/docs/mvc/model/query.md#%E8%BD%BD%E5%85%A5%E5%85%B3%E7%B3%BB%E5%AD%97%E6%AE%B5
				o.LoadRelated(post, "Tags")
				v2, _ := json.Marshal(*post)
				fmt.Println("载入值：", string(v2))
			}
		} else {
			fmt.Println(err)
		}
	}

	fmt.Println()

	{
		// 正向一对一是加载的
		// 反向一对多是不加载的
		user := User{Id: 1}
		o.Read(&user)
		v1, _ := json.Marshal(user)
		fmt.Println("默认值：", string(v1))
		_, _ = o.LoadRelated(&user, "Posts")
		v2, _ := json.Marshal(user)
		fmt.Println("载入值：", string(v2))
	}

	fmt.Println()

	{
		// 反向一对一关系是不加载的
		profile := Profile{Id: 1}
		o.Read(&profile)
		v1, _ := json.Marshal(profile)
		fmt.Println("默认值：", string(v1))
		_, _ = o.LoadRelated(&profile, "User")
		v2, _ := json.Marshal(profile)
		fmt.Println("载入值：", string(v2))
	}
}
