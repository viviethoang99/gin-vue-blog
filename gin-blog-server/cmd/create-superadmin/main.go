package main

import (
	"errors"
	"flag"
	"fmt"
	ginblog "gin-blog/internal"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"log"
	"log/slog"
	"os"

	"gorm.io/gorm"
)

func main() {
	username := flag.String("username", "", "超级管理员账户")
	password := flag.String("password", "", "super administrator password")
	configPath := flag.String("c", "../../config.yml", "configuration file path")
	flag.Parse()

	// Read configuration file based on command-line parameters; other variable initialization depends on the config object
	conf := g.ReadConfig(*configPath)

	//! Handle sqlite3 database path
	conf.SQLite.Dsn = "../" + conf.SQLite.Dsn
	conf.Server.DbLogMode = "silent"

	db := ginblog.InitDatabase(conf)

	if *username == "" || *password == "" {
		log.Fatal("Please specify super administrator account and password")
	}

	createSuperAdmin(db, *username, *password)
}

// Create super administrator
func createSuperAdmin(db *gorm.DB, username, password string) {
	err := db.Transaction(func(tx *gorm.DB) error {
		var userAuth model.UserAuth
		err := db.Where("username = ?", username).First(&userAuth).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if userAuth.ID != 0 {
			return errors.New(userAuth.Username + " account already exists")
		}

		slog.Info("Start creating super administrator")

		// 默认生成一个 super admin 用户
		hashPassword, err := utils.BcryptHash(password)
		if err != nil {
			return errors.New("Password hash generation failed: " + err.Error())
		}

		userAuth = model.UserAuth{
			Username: username,
			Password: hashPassword,
			IsSuper:  true,
			UserInfo: &model.UserInfo{
				Nickname: username,
				Avatar:   "https://raw.githubusercontent.com/szluyu99/gin-vue-blog/main/images/config/superadmin_avatar.jpg",
				Intro:    "This person is lazy and left nothing.",
				Website:  "https://www.hahacode.cn",
			},
		}
		if err := db.Create(&userAuth).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		slog.Error("Failed to create super administrator: " + err.Error())
		os.Exit(0)
	}

	slog.Info(fmt.Sprintf("Successfully created super administrator: %s, password: %s\n", username, password))
}
