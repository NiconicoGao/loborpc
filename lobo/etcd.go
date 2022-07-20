package lobo

import (
	"context"
	"log"
	sync "sync"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type ServiceRegister struct {
	cli           *clientv3.Client
	leaseID       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string //key
	val           string //value
}

func NewServiceRegister(endpoints []string, key, val string, lease int64) (*ServiceRegister, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	ser := &ServiceRegister{
		cli: cli,
		key: key,
		val: val,
	}

	// Register with lease
	if err := ser.putKeyWithLease(lease); err != nil {
		return nil, err
	}

	return ser, nil
}

//Set lease info
func (s *ServiceRegister) putKeyWithLease(lease int64) error {
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}

	//Register service
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	// Keep Lease Alive
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)

	if err != nil {
		return err
	}
	s.leaseID = resp.ID
	log.Println(s.leaseID)
	s.keepAliveChan = leaseRespChan
	log.Printf("Put key:%s  val:%s  success!", s.key, s.val)
	return nil
}

// Lease timeout
func (s *ServiceRegister) ListenLeaseRespChan() {
	for range s.keepAliveChan {

	}
	log.Println("Close leasing")
}

// Terminate Leasing
func (s *ServiceRegister) Close() error {
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}

	log.Println("Close leasing")
	return s.cli.Close()
}

//ServiceDiscovery
type ServiceDiscovery struct {
	cli        *clientv3.Client  //etcd client
	serverList map[string]string //
	lock       sync.Mutex
}

//NewServiceDiscovery
func NewServiceDiscovery(endpoints []string) *ServiceDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 1 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &ServiceDiscovery{
		cli:        cli,
		serverList: make(map[string]string),
	}
}

//WatchService
func (s *ServiceDiscovery) WatchService(prefix string) error {
	resp, err := s.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, ev := range resp.Kvs {
		s.SetServiceList(string(ev.Key), string(ev.Value))
	}

	go s.watcher(prefix)
	return nil
}

//watcher
func (s *ServiceDiscovery) watcher(prefix string) {
	rch := s.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	log.Printf("watching prefix:%s now...", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case clientv3.EventTypeDelete:
				s.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

// New Service Register
func (s *ServiceDiscovery) SetServiceList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serverList[key] = string(val)
	log.Println("put key :", key, "val:", val)
}

// Delete an Active Service
func (s *ServiceDiscovery) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serverList, key)
	log.Println("del key:", key)
}

//Get Service List
func (s *ServiceDiscovery) GetServices() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	addrs := make([]string, 0)

	for _, v := range s.serverList {
		addrs = append(addrs, v)
	}
	return addrs
}

func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}
