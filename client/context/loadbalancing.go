package context

import (
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"strings"
	"sync"
	"github.com/pkg/errors"
)

type ClientManager struct {
	clients []rpcclient.Client
	currentIndex int
	mutex sync.Mutex
}

func NewClientManager(nodeURIs string) (*ClientManager,error) {
	if nodeURIs != "" {
		mgr := &ClientManager{
			mutex: sync.Mutex{},
		}
		nodeUrlArray := strings.Split(nodeURIs, ",")
		for _, url := range nodeUrlArray {
			client := rpcclient.NewHTTP(url, "/websocket")
			mgr.clients = append(mgr.clients, client)
		}
		mgr.currentIndex = 0
		return mgr, nil
	} else {
		return nil, errors.New("missing node URIs")
	}
}

func (mgr *ClientManager) getClient() rpcclient.Client {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	client := mgr.clients[mgr.currentIndex]
	mgr.currentIndex++
	if mgr.currentIndex >= len(mgr.clients){
		mgr.currentIndex = 0
	}
	return client
}