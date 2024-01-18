package main

import (
	"flag"
	"fmt"
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

type ServerOptions struct {
	EndpointAddr string
}

// config{
// 	host1 string
// 	host2 string
// 	}

// func getConfig(s *string) {
// 	flag.StringVar(s, "a", "localhost:8080", "endpoint address")
// 	flag.Parse()
// }

func getConfig() string {
	defaultHost := "localhost:8080"
	host := flag.String("a", defaultHost, "Адрес HTTP-сервера. По умолчанию localhost:8080")
	flag.Parse()
	return *host
}

func main() {

	// opt := ServerOptions{}

	// getConfig(&opt.EndpointAddr)
	host := getConfig()
	fmt.Println("host: ", host)
	appConfig := app.NewAppConfig(host)
	app := app.NewApp(appConfig)

	app.Run()

}
