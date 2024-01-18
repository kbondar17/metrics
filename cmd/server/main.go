package main

import (
	"flag"
	"metrics/internal/app"
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
// if ennvHost := os.Getenv("ADDRESS"); ennvHost != "" {
// 	Host = ennvHost
// }

func getConfig(s *string) {
	flag.StringVar(s, "a", "localhost:8080", "endpoint address")
	flag.Parse()
}

func main() {
	var host string
	getConfig(&host)

	appConfig := app.NewAppConfig(&host)
	app := app.NewApp(appConfig)

	app.Run()

}
