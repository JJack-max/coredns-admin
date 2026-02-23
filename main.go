package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/router"
	"github.com/kexirong/coredns-admin/service"
	"golang.org/x/term"
)

var adduser bool
var confPath string

func main() {
	flag.BoolVar(&adduser, "adduser", false, "add user")
	flag.StringVar(&confPath, "C", "config.yaml", "config file path")
	flag.Parse()

	config.Set(confPath)
	conf := config.Get()
	err := service.EtcdInitClient(conf)
	if err != nil {
		panic(err)
	}
	if adduser {
		addUser(conf.UserEtcdPath)
	}
	// 若设置了 ADMIN_USERNAME 和 ADMIN_PASSWORD 环境变量，自动初始化用户到 etcd
	addUserFromEnv(conf.UserEtcdPath)
	err = router.Router.Run(fmt.Sprintf("%s:%s", conf.Host, conf.Port))
	panic(err)

}

func addUser(prefixPath string) {
	var username, password, confirmPassword string
username:
	fmt.Print("Enter Username：")
	fmt.Scanln(&username)
	if username == "" {
		fmt.Println("Username can not be empty")
		goto username
	}
password:
	fmt.Print("Enter Password: ")
	password = getpassword()
	if len(password) < 6 {
		fmt.Println("\nPassword must be at least 6 characters")
		goto password
	}
	fmt.Print("\nConfirm Password: ")
	confirmPassword = getpassword()
	if confirmPassword != password {
		fmt.Println("\nPassword and confirm password doesn't match")
		goto password
	}

	path := prefixPath + username
	secret := service.MakeSecret(password)
	err := service.EtcdPutKv(path, secret)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nuser %s created\n", username)
	os.Exit(0)
}

func getpassword() string {
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(bytePassword))
}

// addUserFromEnv 从环境变量 ADMIN_USERNAME、ADMIN_PASSWORD 读取并创建用户（用于 Docker 等非交互场景）
func addUserFromEnv(prefixPath string) {
	username := strings.TrimSpace(os.Getenv("ADMIN_USERNAME"))
	password := strings.TrimSpace(os.Getenv("ADMIN_PASSWORD"))
	if username == "" || password == "" {
		fmt.Println("[init] ADMIN_USERNAME or ADMIN_PASSWORD not set, skip init admin user")
		return
	}
	if len(password) < 6 {
		fmt.Printf("[init] ADMIN_PASSWORD must be at least 6 characters, skip init user %s\n", username)
		return
	}
	path := prefixPath + username
	secret := service.MakeSecret(password)
	if err := service.EtcdPutKv(path, secret); err != nil {
		fmt.Printf("[init] failed to init admin user %s: %v\n", username, err)
		return
	}
	fmt.Printf("[init] admin user %s initialized from env\n", username)
}
