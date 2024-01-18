package main

import (
	"flag"
	"metrics/internal/app"
	"os"
)

// type ServerOptions struct {
// 	EndpointAddr string `env:"ADDRESS"`
// }

// func (o *ServerOptions) ParseArgs() {
// 	flag.StringVar(&o.EndpointAddr, "a", "localhost:8080", "endpoint address")
// 	flag.Parse()
// }

// func (o *ServerOptions) ParseEnv() {
// 	host := os.Getenv("ADDRESS")
// 	if host != "" {
// 		o.EndpointAddr = os.Getenv("ADDRESS")
// 	}
// }

// попробуйте передать строку через указатель.
func getConfig() *string {
	var Host string

	flag.StringVar(&Host, "a", "localhost:8080", "endpoint address")
	flag.Parse()

	if ennvHost := os.Getenv("ADDRESS"); ennvHost != "" {
		Host = ennvHost
	}

	return &Host
}

func main() {

	host := getConfig()
	appConfig := app.NewAppConfig(host)
	app := app.NewApp(appConfig)

	app.Run()

}
