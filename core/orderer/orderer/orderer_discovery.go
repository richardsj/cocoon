package orderer

import (
	"fmt"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/ellcrys/util"
	"github.com/hashicorp/consul/api"
	"github.com/ncodes/cocoon/core/config"
	"github.com/ncodes/cocoon/core/scheduler"
)

var discoveryLog = config.MakeLogger("orderer.discovery", "orderer")

// Discovery defines a structure for fetching a list of addresses of orderers
// accessible in the cluster.
type Discovery struct {
	sync.Mutex
	sd           scheduler.ServiceDiscovery
	orderersAddr []string
	ticker       *time.Ticker
	OnUpdateFunc func(addrs []string)
}

// NewDiscovery creates a new discovery object.
// Returns error if unable to connector to the service discovery.
// Setting the env variable `CONSUL_ADDR` will override the default config address.
func NewDiscovery() (*Discovery, error) {
	cfg := api.DefaultConfig()
	cfg.Address = util.Env("CONSUL_ADDR", cfg.Address)
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %s", err)
	}
	return &Discovery{
		sd: scheduler.NewNomadServiceDiscovery(client),
	}, nil
}

// Discover fetches a list of orderer service addresses
// via consul service discovery API.
// For development purpose, If DEV_ORDERER_ADDR is set,
// it will fetch the orderer address from the env variable.
func (od *Discovery) discover() error {

	var err error

	if len(os.Getenv("DEV_ORDERER_ADDR")) > 0 {
		od.orderersAddr = []string{os.Getenv("DEV_ORDERER_ADDR")}
		return nil
	}

	_orderers, err := od.sd.GetByID("orderers", nil)
	if err != nil {
		return err
	}

	var orderers []string
	for _, orderer := range _orderers {
		orderers = append(orderers, fmt.Sprintf("%s:%d", orderer.IP, int(orderer.Port)))
	}

	od.Lock()
	od.orderersAddr = orderers
	od.Unlock()
	return nil
}

// Discover starts a ticker that discovers and updates the list
// of orderer addresses. It will perform the discovery immediately
// and will return error if it fails, otherwise nil is returned and
// subsequent discovery will be performed periodically
func (od *Discovery) Discover() error {

	// run immediately
	if err := od.discover(); err != nil {
		return err
	}

	// run on interval
	od.ticker = time.NewTicker(15 * time.Second)
	for _ = range od.ticker.C {
		err := od.discover()
		if err != nil {
			discoveryLog.Error(err.Error())
			if od.OnUpdateFunc != nil {
				od.OnUpdateFunc(od.GetAddrs())
			}
		}
	}
	return nil
}

// GetAddrs returns the list of discovered addresses
func (od *Discovery) GetAddrs() []string {
	return od.orderersAddr
}

// GetGRPConn dials a random orderer address and returns a
// grpc connection. If no orderer address has been discovered, nil and are error are returned.
func (od *Discovery) GetGRPConn() (*grpc.ClientConn, error) {

	od.Lock()

	var selected string
	if len(od.orderersAddr) == 0 {
		od.Unlock()
		return nil, fmt.Errorf("no known orderer address")
	}

	if len(od.orderersAddr) == 1 {
		selected = od.orderersAddr[0]
	} else {
		selected = od.orderersAddr[util.RandNum(0, len(od.orderersAddr))]
	}

	od.Unlock()

	client, err := grpc.Dial(selected, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Stop stops the discovery ticker
func (od *Discovery) Stop() {
	od.ticker.Stop()
}
