package main

import (
	"flag"
	gohttp "net/http"
	"os"

	"github.com/ian-kent/go-log/log"
	comcfg "github.com/sskorol/mailhog/config"
	"github.com/sskorol/mailhog/http"
	"github.com/sskorol/mailhog/server/api"
	"github.com/sskorol/mailhog/server/config"
	"github.com/sskorol/mailhog/server/smtp"
	"github.com/sskorol/mailhog/ui/assets"
)

var (
	conf    *config.Config
	comconf *comcfg.Config
	exitCh  chan int
)

func configure() {
	comcfg.RegisterFlags()
	config.RegisterFlags()
	flag.Parse()
	conf = config.Configure()
	comconf = comcfg.Configure()
}

func main() {
	configure()

	if comconf.AuthFile != "" {
		http.AuthFile(comconf.AuthFile)
	}

	exitCh = make(chan int)
	cb := func(r gohttp.Handler) {
		api.CreateAPI(conf, r)
	}
	go http.Listen(conf.APIBindAddr, assets.Asset, exitCh, cb)
	go smtp.Listen(conf, exitCh)

	for {
		select {
		case <-exitCh:
			log.Printf("Received exit signal")
			os.Exit(0)
		}
	}
}
