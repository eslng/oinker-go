package main

import (
	"github.com/eslng/oinker-go/controller"
	"github.com/eslng/oinker-go/service"

	log "github.com/Sirupsen/logrus"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/karlkfi/inject"

	"flag"
	"net/http"
)

func main() {
	flagSet := flag.CommandLine
	flags := parseFlags(flagSet)
	log.Infof("Flags: %+v", flags)

	oinker := &Oinker{}

	if *flags.cassandraAddr != "" {
		oinker.CQLHosts = []string{*flags.cassandraAddr}
		oinker.CQLReplicationFactor = *flags.cassandraRepl
	}

	graph := oinker.NewGraph()
	defer graph.Finalize()

	// initialize cassandra (connection, keyspace, tables)
	var oinkRepo service.OinkRepo
	inject.ExtractAssignable(graph, &oinkRepo)
	svc, ok := oinkRepo.(inject.Initializable)
	if ok {
		svc.Initialize()
	}

	var mux controller.MuxServer
	inject.ExtractAssignable(graph, &mux)

	var controllers []controller.Controller
	inject.FindAssignable(graph, &controllers)
	for _, c := range controllers {
		log.Infof("Registering controller:", c.Name())
		c.RegisterHandlers(mux)
	}

	// serve and listen for shutdown signals
	gracehttp.Serve(
		&http.Server{Addr: *flags.address, Handler: mux},
	)
}
