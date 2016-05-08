package main

import (
	"github.com/eslng/oinker-go/controller"
	"github.com/eslng/oinker-go/service"

	"github.com/gocql/gocql"
	"github.com/karlkfi/inject"

	"net/http"
	"os"
)

type Oinker struct {
	CQLHosts             []string
	CQLReplicationFactor int
}

func (o *Oinker) NewGraph() inject.Graph {
	graph := inject.NewGraph()

	var instanceName string
	graph.Define(&instanceName, inject.NewProvider(func() string {
		name := os.Getenv("OINKER_INSTANCE_NAME")
		if name == "" {
			return "instance-unknown"
		}
		return name
	}))

	var server *http.ServeMux
	graph.Define(&server, inject.NewProvider(http.NewServeMux))

	var cqlCluster *gocql.ClusterConfig
	graph.Define(&cqlCluster, inject.NewProvider(func() *gocql.ClusterConfig {
		return gocql.NewCluster(o.CQLHosts...)
	}))

	var cqlSession *service.CQLSession
	graph.Define(&cqlSession, inject.NewProvider(service.NewCQLSession, &cqlCluster))

	var oinkRepo service.OinkRepo
	if len(o.CQLHosts) > 0 {
		graph.Define(&oinkRepo, inject.NewProvider(service.NewCQLOinkRepo, &cqlSession, &o.CQLReplicationFactor))
	} else {
		graph.Define(&oinkRepo, inject.NewProvider(service.NewMockOinkRepo))
	}

	var readyController *controller.ReadyController
	graph.Define(&readyController, inject.NewProvider(controller.NewReadyController, &oinkRepo, &instanceName))

	var assetsController *controller.AssetsController
	graph.Define(&assetsController, inject.NewProvider(controller.NewAssetsController))

	var indexController *controller.IndexController
	graph.Define(&indexController, inject.NewProvider(controller.NewIndexController, &oinkRepo, &instanceName))

	var oinkController *controller.OinkController
	graph.Define(&oinkController, inject.NewProvider(controller.NewOinkController, &oinkRepo))

	return graph
}
