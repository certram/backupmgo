package main

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	config   Config
	filePath string = "./config/atlas_local.yaml"
	// 具体备份时间（默认是5个域）；分 时 月的某天 某月 星期几
	concreteTime string // 每天的9:15分备份一次
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	config.InitConfig(filePath)
	err := config.CheckConfig()
	if err != nil {
		logrus.Panic("验证配置文件错误")
	} else {
		logrus.Info("完成读取配置")
	}
	concreteTime = config.ConcreteTime

}

func main() {
	// 周期性执行备份操作
	c := cron.New()
	command := "mongodump --uri mongodb+srv://" + config.UserName + ":" + config.Password + "@" + config.ClusterUrl + "/" + config.DBName
	c.AddFunc(concreteTime, func() {
		// 如果有dump文件夹，则先删除dump文件夹及其包含的所有子目录和所有文件
		err := os.RemoveAll("./dump/")
		if err != nil {
			logrus.Panic("remove fileFolder failed")
		}

		time.Sleep(1 * time.Minute)
		exec_shell(command)
	})
	c.Start()

	select {} // 阻塞主程序

}

type Config struct {
	UserName     string `yaml:"userName"`
	Password     string `yaml:"password"`
	ClusterUrl   string `yaml:"clusterUrl"`
	DBName       string `yaml:"dbName"`
	ConcreteTime string `yaml:"concreteTime"`
}

// 读取yaml文件到Config struct
func (conf *Config) InitConfig(filePath string) *Config {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("get yamlFile err %v ", err)
	}

	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		logrus.Fatalf("Unmarshal: %v", err)
	}
	return conf

}
func (conf *Config) CheckConfig() error {
	if conf.UserName == "" || conf.Password == "" || conf.ClusterUrl == "" || conf.DBName == "" || conf.ConcreteTime == "" {
		logrus.Error("please configure the atlas_local.yaml correctly,every field can't be \"\"")
		return errors.New("config is not right")
	}
	return nil
}

// mongodump --uri mongodb+srv://firstUser:<PASSWORD>@cluster0.deeze6y.mongodb.net/<DATABASE>
func exec_shell(command string) {
	// var once sync.Once
	// var command string
	// once.Do(func() {
	// 	command = "mongodump --uri mongodb+srv://" + config.UserName + ":" + config.Password + "@" + config.ClusterUrl + "/" + config.DBName
	// })

	cmd := exec.Command("/bin/bash", "-c", command)
	err := cmd.Run()
	if err != nil {
		logrus.Fatal(err.Error())
	}
	logrus.Info("backup successfully\n")

}
