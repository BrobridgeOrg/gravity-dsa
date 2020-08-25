package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type App struct {
	id string
}

func CreateApp() *App {

	// Using hostname (pod name) by default
	host, err := os.Hostname()
	if err != nil {
		log.Error(err)
		return nil
	}

	host = strings.ReplaceAll(host, ".", "_")

	id := fmt.Sprintf("gravity-%s", host)

	return &App{
		id: id,
	}
}

func (a *App) Init() error {

	log.WithFields(log.Fields{
		"a_id": a.id,
	}).Info("Starting application")

	return nil
}

func (a *App) Uninit() {
}

func (a *App) Run() error {

	port := strconv.Itoa(viper.GetInt("service.port"))
	err := a.InitGRPCServer(":" + port)
	if err != nil {
		return err
	}

	return nil
}
