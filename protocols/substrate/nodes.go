package substrate

import (
	"fmt"

	"github.com/taubyte/odo/config"
	counters "github.com/taubyte/odo/protocols/substrate/components/counters"
	database "github.com/taubyte/odo/protocols/substrate/components/database"
	http "github.com/taubyte/odo/protocols/substrate/components/http"
	ipfs "github.com/taubyte/odo/protocols/substrate/components/ipfs"
	p2p "github.com/taubyte/odo/protocols/substrate/components/p2p"
	pubSub "github.com/taubyte/odo/protocols/substrate/components/pubsub"
	smartOps "github.com/taubyte/odo/protocols/substrate/components/smartops"
	storage "github.com/taubyte/odo/protocols/substrate/components/storage"
)

func attachNodesError(name string, err error) error {
	err = fmt.Errorf("creating node %s failed with %s", name, err.Error())
	logger.Error(err)

	return err
}

func (srv *Service) Verbose() bool {
	return srv.verbose
}

func (srv *Service) attachNodes(config *config.Protocol) (err error) {
	// Needs to happen first, as others depend on it
	if err = srv.attachNodeCounters(config); err != nil {
		return attachNodesError("counters", err)
	}

	// Needs to happen second, as others depend on it
	if err = srv.attachNodeSmartOps(config); err != nil {
		return attachNodesError("smartops", err)
	}

	if err = srv.attachNodePubSub(config); err != nil {
		return attachNodesError("pubsub", err)
	}

	if err = srv.attachNodeIpfs(config); err != nil {
		return attachNodesError("ipfs", err)
	}

	if err = srv.attachNodeDatabase(config); err != nil {
		return attachNodesError("database", err)
	}

	if err = srv.attachNodeStorage(config); err != nil {
		return attachNodesError("storage", err)
	}

	if err = srv.attachNodeP2P(config); err != nil {
		return attachNodesError("p2p", err)
	}

	if err = srv.attachNodeHttp(config); err != nil {
		return attachNodesError("http", err)
	}

	return nil
}

func (srv *Service) attachNodeHttp(config *config.Protocol) (err error) {
	ops := []http.Option{}

	if config.DevMode {
		ops = append(ops, http.DvKey(config.DomainValidation.PublicKey))
	}

	srv.nodeHttp, err = http.New(srv, ops...)
	return
}

func (srv *Service) attachNodePubSub(config *config.Protocol) (err error) {
	srv.nodePubSub, err = pubSub.New(srv)
	return
}

func (srv *Service) attachNodeIpfs(config *config.Protocol) (err error) {
	ipfsPort, ok := config.Ports["ipfs"]
	if !ok {
		err = fmt.Errorf("did not find ipfs port in config")
		return

	}

	srv.nodeIpfs, err = ipfs.New(srv.node.Context(), ipfs.Public(), ipfs.Listen([]string{fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", ipfsPort)}))
	return
}

func (srv *Service) attachNodeDatabase(config *config.Protocol) (err error) {
	srv.nodeDatabase, err = database.New(srv)
	return
}

func (srv *Service) attachNodeStorage(config *config.Protocol) (err error) {
	srv.nodeStorage, err = storage.New(srv)
	return
}

func (srv *Service) attachNodeP2P(config *config.Protocol) (err error) {
	srv.nodeP2P, err = p2p.New(srv)
	return
}

func (srv *Service) attachNodeCounters(config *config.Protocol) (err error) {
	srv.nodeCounters, err = counters.New(srv)
	return
}

func (srv *Service) attachNodeSmartOps(config *config.Protocol) (err error) {
	ops := []smartOps.Option{}

	if config.DevMode {
		ops = append(ops, smartOps.Dev())
	}

	if config.Verbose {
		ops = append(ops, smartOps.Verbose())
	}

	srv.nodeSmartOps, err = smartOps.New(srv, ops...)
	return
}