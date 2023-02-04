package repository

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"user/pkg/pb"
)

type User struct {
	Id           uint32 `json:"id"`
	Name         string `json:"name"`
	NickName     string `json:"nick_name"`
	Password     string
	EncryptedPwd []byte
	db           *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{db: db}
}

func (u *User) ShowUserInfo(req *pb.UserRequest) (*User, error) {
	//if err := u.db.Where("name=?", req.UserName).First(&u).Error; err != nil && err == gorm.ErrRecordNotFound {
	if err := u.db.Table("user").Where("name=?", req.UserName).First(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) CreateUser(req *pb.UserRequest) error {
	var count int64 = 0
	err := u.db.Table("user").Where("name=?", req.UserName).Count(&count).Error
	if count > 0 || err != nil {
		return fmt.Errorf("user %s exist: %s", req.UserName, err.Error())
	}
	user := User{
		Name:     req.UserName,
		NickName: req.NickName,
	}
	encryptedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return fmt.Errorf("encrypt user %s failed", req.UserName)
	}
	user.EncryptedPwd = encryptedPwd

	return u.db.Create(&user).Error
}

func (u *User) Translate() *pb.UserModel {
	return &pb.UserModel{
		UserId:   u.Id,
		UserName: u.Name,
		NickName: u.NickName,
	}
}
