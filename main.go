package main

import (
	"fmt"
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
	cfgFile, found := os.LookupEnv("CONFIG_FILE")
	if !found {
		log.Fatal(fmt.Errorf("configuration file must be specified in CONFIG_FILE enviornment variable"))
	}
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	conf := service.NewConfig()
	if err = yaml.Unmarshal(data, conf); err != nil {
		log.Fatal(err)
	}
	svc, err := service.New(conf)
	if err != nil {
		log.Fatal(err)
	}
	svc.Run()
}
