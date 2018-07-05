package context

import (
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"strings"
	"sync"
	"github.com/pkg/errors"
	"github.com/MrXJC/GoLoadBalance"
	"strconv"
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


type ClientManagerLB struct {
	clients map[string]rpcclient.Client
	mgr *balance.BalanceMgr
	current string
	mutex sync.Mutex
}

func NewClientManagerLB(nodeURIs string) (*ClientManagerLB,error) {
	if nodeURIs != "" {
		clientMgr := &ClientManagerLB{
			mutex: sync.Mutex{},
			clients: make(map[string]rpcclient.Client),
		}

		var addr *balance.NodeAddr
		var addrs balance.NodeAddrs

		nodeUrlArray := strings.Split(nodeURIs, ",")
		for _, url := range nodeUrlArray {
			client := rpcclient.NewHTTP(url, "/websocket")
			clientMgr.clients[nodeURIs] = client

			arr := strings.Split(url, ":")
			port,err := strconv.Atoi(arr[1])

			if err!=nil{
				return nil,errors.New("Port isn't a number string")
			}

			ip := arr[0]
			addr = balance.NewNodeAddr(ip,port,1)
			addrs = append(addrs, addr)
		}
		clientMgr.mgr = balance.NewBalanceMgr(addrs)

		clientMgr.mgr.RegisterBalancer("randomweight",&balance.RandomBalance{})
		clientMgr.mgr.RegisterBalancer("random",&balance.RandomWeightBalance{})
		clientMgr.mgr.RegisterBalancer("roundrobin",&balance.RoundRobinBalance{})
		clientMgr.mgr.RegisterBalancer("roundrobinweight",&balance.RoundRobinWeightBalance{})

		clientMgr.current,_,_ = clientMgr.mgr.GetAddrString("roundrobin")
		return clientMgr, nil

	} else {
		return nil, errors.New("missing node URIs")
	}
}

func NewClientManagerLBwithWeight(nodeURIs string) (*ClientManagerLB,error) {
	if nodeURIs != "" {
		clientMgr := &ClientManagerLB{
			mutex: sync.Mutex{},
			clients: make(map[string]rpcclient.Client),
		}

		var addr *balance.NodeAddr
		var addrs balance.NodeAddrs

		nodeUrlArray := strings.Split(nodeURIs, ",")
		for _, url := range nodeUrlArray {
			client := rpcclient.NewHTTP(url, "/websocket")
			clientMgr.clients[nodeURIs] = client

			arr := strings.Split(url, ":")
			port,err := strconv.Atoi(arr[1])

			if err!=nil{
				return nil,errors.New("Port isn't a number string")
			}

			weight,err := strconv.Atoi(arr[2])

			if err!=nil{
				return nil,errors.New("Weight isn't a number string")
			}

			ip := arr[0]
			addr = balance.NewNodeAddr(ip,port,weight)
			addrs = append(addrs, addr)
		}
		clientMgr.mgr = balance.NewBalanceMgr(addrs)

		clientMgr.mgr.RegisterBalancer("random",&balance.RandomBalance{})
		clientMgr.mgr.RegisterBalancer("randomweight",&balance.RandomWeightBalance{})
		clientMgr.mgr.RegisterBalancer("roundrobin",&balance.RoundRobinBalance{})
		clientMgr.mgr.RegisterBalancer("roundrobinweight",&balance.RoundRobinWeightBalance{})

		clientMgr.current,_,_ = clientMgr.mgr.GetAddrString("roundrobin")
		return clientMgr, nil

	} else {
		return nil, errors.New("missing node URIs")
	}
}

func (clientMgr *ClientManagerLB) getClient() rpcclient.Client {
	clientMgr.mutex.Lock()
	defer clientMgr.mutex.Unlock()

	clientMgr.current,_,_ = clientMgr.mgr.GetAddrString("roundrobin")
	client := clientMgr.clients[clientMgr.current]

	return client
}

func (clientMgr *ClientManagerLB) getClientByName(name string) rpcclient.Client {
	clientMgr.mutex.Lock()
	defer clientMgr.mutex.Unlock()

	clientMgr.current,_,_ = clientMgr.mgr.GetAddrString(name)
	client := clientMgr.clients[clientMgr.current]

	return client
}

func (clientMgr *ClientManagerLB) getClientByNameDebug(name string) (rpcclient.Client,string) {
	clientMgr.mutex.Lock()
	defer clientMgr.mutex.Unlock()

	clientMgr.current,_,_ = clientMgr.mgr.GetAddrString(name)
	client := clientMgr.clients[clientMgr.current]

	return client,clientMgr.current
}