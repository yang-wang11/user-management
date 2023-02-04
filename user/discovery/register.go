package discovery

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"log"
	"strings"
	"time"
)

type BindInfo struct {
	Name   string
	Detail string
	Server *Server
}

type Register struct {
	context    context.Context
	closeCh    chan struct{}
	leaseId    clientv3.LeaseID
	aliveCh    <-chan *clientv3.LeaseKeepAliveResponse
	etcdClient *clientv3.Client
	bind       *BindInfo
	logger     *logrus.Logger
	epManager  endpoints.Manager
}

func NewRegister(etcdAddress []string, logger *logrus.Logger) *Register {

	ctx := context.Background()

	dialTimeout := 3

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdAddress,
		DialTimeout: time.Duration(dialTimeout) * time.Second,
	})
	if err != nil {
		return nil
	}

	return &Register{
		context:    ctx,
		etcdClient: client,
		logger:     logger,
		closeCh:    make(chan struct{}),
	}
}

func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

func (r *Register) Register(srv *Server, ttl int64) error {

	if path := strings.Split(srv.Addr, ":")[0]; path == "" {
		return errors.New("invalid registered path")
	}

	// create lease resource
	leaseResp, err := r.etcdClient.Grant(r.context, ttl)
	if err != nil {
		return err
	}
	r.leaseId = leaseResp.ID

	registeredPath := BuildRegisteredPath(srv)

	r.bind = &BindInfo{
		Name:   registeredPath,
		Detail: fmt.Sprintf("%s/detail", registeredPath),
		Server: srv,
	}

	// create the endpoint manager
	if r.epManager, err = endpoints.NewManager(r.etcdClient, r.bind.Server.Name); err != nil {
		return err
	}

	if err = r.register(); err != nil {
		return err
	}

	go r.keepAlive()

	return nil
}

func (r *Register) keepAlive() {

	aliveCh, err := r.etcdClient.KeepAlive(r.context, r.leaseId)
	if err != nil {
		log.Fatalf("failed to active the %s", r.bind.Name)
	}

	for {
		select {
		case res := <-aliveCh:
			if res == nil {
				log.Fatalf("failed to keep the %s alive", r.bind.Name)
			}
		case <-r.closeCh:
			r.logger.Infof("start to clean the resources")
			if _, err := r.etcdClient.Delete(r.context, r.bind.Detail); err != nil {
				r.logger.Fatalf("failed to unregister %s", err.Error())
			}
			if err = r.epManager.DeleteEndpoint(r.context, r.bind.Name); err != nil {
				r.logger.Fatalf("failed to remove the endpoint of %s: %s", r.bind.Name, err.Error())
			}
			if _, err = r.etcdClient.Revoke(r.context, r.leaseId); err != nil {
				r.logger.Fatalf("failed to revoke lease id %v: %s", r.leaseId, err.Error())
			}
		}
	}
}

func (r *Register) register() error {

	data, err := MarshalRegisteredServer(r.bind.Server)
	if err != nil {
		return err
	}

	// update detailed information
	if _, err = r.etcdClient.Put(r.context, r.bind.Detail, string(data), clientv3.WithLease(r.leaseId)); err != nil {
		return err
	}

	// bind endpoints
	if err = r.epManager.AddEndpoint(r.context, r.bind.Name, endpoints.Endpoint{
		Addr: r.bind.Server.Addr,
	}, clientv3.WithLease(r.leaseId)); err != nil {
		return err
	}

	return nil
}
