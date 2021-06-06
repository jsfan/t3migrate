package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/jsfan/t3migrate/internal/config"
	"github.com/jsfan/t3migrate/internal/storage"
	"os"
)

func main() {
	cfgFileOpt := flag.String("config", "databases.yaml", "Configuration YAML file")
	dryOpt := flag.Bool("dry", false, "DRY RUN")
	flag.Parse()

	if *dryOpt {
		fmt.Println("*** DRY RUN ***")
	}

	cfg, err := config.ReadConfig(*cfgFileOpt)
	if err != nil {
		glog.Fatalf("Could not load configuration: %+v", err)
	}

	srcStore := storage.Store{}
	if err = srcStore.Connect(&cfg.Source, *dryOpt); err != nil {
		glog.Fatalf("Could not connect to source database: %+v", err)
	}
	if err = srcStore.DeleteUnusedElements(); err != nil {
		glog.Errorf("Failed to delete unused elements: %+v", err)
		os.Exit(1)
	}
}
