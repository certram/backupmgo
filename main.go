package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"
)

var (
	config   Config
	filePath string = "./config/atlas.yaml"
	// 具体备份时间（默认是5个域）；分 时 月的某天 某月 星期几
	concreteTime string = "15 9 * * *" // 每天的9:15分备份一次
)

func main() {
	// 只读取一次，如果更换了yaml文件，需重新go build
	config.InitConfig(filePath)
	fmt.Println("完成读取配置")
	// 周期性执行备份操作
	c := cron.New()
	command := "mongodump --uri mongodb+srv://" + config.UserName + ":" + config.Password + "@" + config.ClusterUrl + "/" + config.DBName
	c.AddFunc(concreteTime, func() {
		// 如果有dump文件夹，则先删除dump文件夹及其包含的所有子目录和所有文件
		err := os.RemoveAll("./dump/")
		if err != nil {
			panic("remove fileFolder failed")
		}
		fmt.Println("执行备份")
		time.Sleep(1 * time.Minute)
		exec_shell(command)
	})
	c.Start()

	select {} // 阻塞主程序

}

type Config struct {
	UserName   string `yaml:"userName"`
	Password   string `yaml:"password"`
	ClusterUrl string `yaml:"clusterUrl"`
	DBName     string `yaml:"dbName"`
}

// 读取yaml文件到Config struct
func (conf *Config) InitConfig(filePath string) *Config {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("get yamlFile err %v ", err)
	}

	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return conf

}

// mongodump --uri mongodb+srv://firstUser:<PASSWORD>@cluster0.deeze6y.mongodb.net/<DATABASE>
// mongodump --uri mongodb+srv://firstUser:qweasd123@cluster0.deeze6y.mongodb.net/market-cli-test
func exec_shell(command string) {
	// var once sync.Once
	// var command string
	// once.Do(func() {
	// 	command = "mongodump --uri mongodb+srv://" + config.UserName + ":" + config.Password + "@" + config.ClusterUrl + "/" + config.DBName
	// })

	cmd := exec.Command("/bin/bash", "-c", command)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("backup successfully\n")
}
