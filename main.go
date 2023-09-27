package main

import (
	"github.com/ghodss/yaml"
	"github.com/harshabangi/url-shortener/internal/service"
	"io/ioutil"
	"log"
	"os"
)

// @title URL Shortener API
// @version 1.0
// @description URL Shortener Server
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	data, err := ioutil.ReadFile(os.Getenv("CONFIG_FILE"))
	if err != nil {
		log.Fatal(err)
	}
	conf := service.NewConfig()
	if err = yaml.Unmarshal(data, conf); err != nil {
		log.Fatal(err)
	}
	svc, err := service.NewService(conf)
	if err != nil {
		log.Fatal(err)
	}
	svc.Run()
}
