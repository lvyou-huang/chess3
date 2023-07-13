package service

import (
	"chess/dao"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func Register(name string, password string) error {
	Db, err := dao.OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = Db.Exec("insert into users (id,password,name) value (?,?,?)", nil, password, name)
	if err != nil {
		log.Println(err)
		return errors.New("创建用户失败")
	}
	return nil
}

func Login(name string, password string, c *gin.Context) error {
	var (
		user struct {
			Id       int
			Name     string
			Password string
		}
		Db, err = dao.OpenDb()
	)

	if err != nil {
		log.Println(err)
		return err
	}
	_, err = Db.Query("select * from users where name =? ", name)
	if err != nil {
		log.Println(err)
		return err
	}

	row := Db.QueryRow("select * from users where name = ?", name)

	if err != nil {
		log.Println(err)
		return err
	}
	err = row.Scan(&user.Id, &user.Password, &user.Name)
	if err != nil {
		log.Println(err)
		return err
	}
	if user.Password == password {
		log.Println(strconv.Itoa(user.Id))
		c.SetCookie("user_id", strconv.Itoa(user.Id), 60*60*24, "/", "localhost", false, true)
		c.SetCookie("name", user.Name, 60*60*24, "/", "localhost", false, true)

		return nil
	} else {
		return errors.New("密码错误")
	}
}
