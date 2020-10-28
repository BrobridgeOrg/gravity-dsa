package main

import (
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "go.uber.org/automaxprocs"

	app "gravity-dsa/app"
)

func init() {

	// From the environment
	viper.SetEnvPrefix("GRAVITY_DSA")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// From config file
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Warn("No configuration file was loaded")
	}

	runtime.GOMAXPROCS(16)

	/*
		go func() {

			f, err := os.Create("cpu-profile.prof")
			if err != nil {
				log.Fatal(err)
			}

			pprof.StartCPUProfile(f)

			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt, os.Kill)
			<-sig
			pprof.StopCPUProfile()

			os.Exit(0)
		}()
	*/
}

func main() {

	// Initializing application
	a := app.CreateApp()

	err := a.Init()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Starting application
	err = a.Run()
	if err != nil {
		log.Fatal(err)
		return
	}
}
