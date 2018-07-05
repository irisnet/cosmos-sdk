package context

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestLoadBalancing(t *testing.T) {
	nodeURIs := "10.10.10.10:26657,20.20.20.20:26657,30.30.30.30:26657"
	clientMgr,err := NewClientManager(nodeURIs)
	assert.Empty(t,err)
	endpoint := clientMgr.getClient()
	clientMgr.getClient()
	clientMgr.getClient()
	assert.Equal(t,endpoint,clientMgr.getClient())
}


func TestClientManagerLB(t *testing.T) {
	nodeURIs := "10.10.10.10:26657,20.20.20.20:26657,30.30.30.30:26657"
	clientMgr, err := NewClientManagerLB(nodeURIs)
	assert.Empty(t, err)

	note := map[string]int{}
	for i := 0; i < 100; i++ {
		_,s := clientMgr.getClientByNameDebug("roundrobin")
		fmt.Println(s)
		if note[s] != 0 {
			note[s]++
		} else {
			note[s] = 1
		}
	}
	for k, v := range note {
		fmt.Println(k, " ", v)
	}


	note = map[string]int{}
	for i := 0; i < 100; i++ {
		_,s := clientMgr.getClientByNameDebug("random")
		fmt.Println(s)
		if note[s] != 0 {
			note[s]++
		} else {
			note[s] = 1
		}
	}
	for k, v := range note {
		fmt.Println(k, " ", v)
	}
}

func TestClientManagerLBwithWeight(t *testing.T) {
	nodeURIs := "10.10.10.10:26657:1,20.20.20.20:26657:2,30.30.30.30:26657:3"
	clientMgr,err := NewClientManagerLBwithWeight(nodeURIs)
	assert.Empty(t,err)

	note := map[string]int{}
	for i := 0; i < 100; i++ {
		_,s := clientMgr.getClientByNameDebug("roundrobinweight")
		fmt.Println(s)
		if note[s] != 0 {
			note[s]++
		} else {
			note[s] = 1
		}
	}
	for k, v := range note {
		fmt.Println(k, " ", v)
	}

	note = map[string]int{}
	for i := 0; i < 100; i++ {
		_,s := clientMgr.getClientByNameDebug("randomweight")
		fmt.Println(s)
		if note[s] != 0 {
			note[s]++
		} else {
			note[s] = 1
		}
	}
	for k, v := range note {
		fmt.Println(k, " ", v)
	}
}