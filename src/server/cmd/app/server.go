package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"speak-sphere/pkg/server/conf"
	"speak-sphere/pkg/server/router"
)

func readConfig(cfg *conf.Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return err
	}
	return nil
}

func Run() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}
	exeDir := filepath.Dir(exePath)
	configFilePath := exeDir + "/cmd/app/config.yaml"
	if err := readConfig(&conf.Cfg, configFilePath); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	conf.Cfg.Sqlite.Path = exeDir + "/" + conf.Cfg.Sqlite.Path

	gin.SetMode(conf.Cfg.Mode)
	conf.InitLogger(conf.Cfg.Log)
	conf.InitDatabase(conf.Cfg.Sqlite)

	r := gin.Default()
	r.Use(router.CorsMiddleware())
	r.Use(router.AuthMiddleware())

	router.RegisterBookRouter(r)
	router.RegisterEntryRouter(r)
	router.RegisterDictionaryRouter(r)
	router.RegisterUserRouter(r)

	if err := r.Run(fmt.Sprintf(":%d", conf.Cfg.Server.Port)); err != nil {
		log.Fatalf("启动服务失败: %v", err)
		return
	}
}
