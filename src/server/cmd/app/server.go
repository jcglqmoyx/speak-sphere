package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"speak-sphere/pkg/server/conf"
	"speak-sphere/pkg/server/router"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
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
	// 直接使用相对路径
	configFilePath := "cmd/app/config.yaml"
	if err := readConfig(&conf.Cfg, configFilePath); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// SQLite路径改为相对于当前目录
	if !filepath.IsAbs(conf.Cfg.Sqlite.Path) {
		conf.Cfg.Sqlite.Path = filepath.Join(".", conf.Cfg.Sqlite.Path)
	}

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
