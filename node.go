package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	filestore "github.com/ipfs/go-ipfs/filestore"

	core "github.com/ipfs/go-ipfs/core"
	coreapi "github.com/ipfs/go-ipfs/core/coreapi"
	mock "github.com/ipfs/go-ipfs/core/mock"
	keystore "github.com/ipfs/go-ipfs/keystore"
	repo "github.com/ipfs/go-ipfs/repo"

	datastore "github.com/ipfs/go-datastore"
	syncds "github.com/ipfs/go-datastore/sync"
	config "github.com/ipfs/go-ipfs-config"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
	ci "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
)

const testPeerID = "QmTFauExutTsy4XP6JbMFcw2Wa9645HJt2bTqL6qYDCKfe"

// MakeAPISwarm returns the CoreInterfase
func MakeAPISwarm(ctx context.Context, fullIdentity bool) (coreiface.CoreAPI, error) {
	mn := mocknet.New(ctx)

	nodes := make([]*core.IpfsNode, 1)
	var api coreiface.CoreAPI
	var ident config.Identity

	if fullIdentity {
		sk, pk, err := ci.GenerateKeyPair(ci.RSA, 512)
		if err != nil {
			return nil, err
		}

		id, err := peer.IDFromPublicKey(pk)
		if err != nil {
			return nil, err
		}

		kbytes, err := sk.Bytes()
		if err != nil {
			return nil, err
		}

		ident = config.Identity{
			PeerID:  id.Pretty(),
			PrivKey: base64.StdEncoding.EncodeToString(kbytes),
		}
	} else {
		ident = config.Identity{
			PeerID: testPeerID,
		}
	}

	c := config.Config{}
	c.Addresses.Swarm = []string{fmt.Sprintf("/ip4/127.0.%d.1/tcp/4001", 0)}
	c.Identity = ident
	c.Experimental.FilestoreEnabled = true

	ds := datastore.NewMapDatastore()
	r := &repo.Mock{
		C: c,
		D: syncds.MutexWrap(ds),
		K: keystore.NewMemKeystore(),
		F: filestore.NewFileManager(ds, filepath.Dir(os.TempDir())),
	}

	node, err := core.NewNode(ctx, &core.BuildCfg{
		Repo:   r,
		Host:   mock.MockHostOption(mn),
		Online: fullIdentity,
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	})
	if err != nil {
		return nil, err
	}
	nodes[0] = node

	api, err = coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	errLink := mn.LinkAll()
	if errLink != nil {
		return nil, errLink
	}
	// bsinf := core.BootstrapConfigWithPeers(
	// 	[]pstore.PeerInfo{
	// 		nodes[0].Peerstore.PeerInfo(nodes[0].Identity),
	// 	},
	// )
	// for _, n := range nodes[1:] {
	// 	if err := n.Bootstrap(bsinf); err != nil {
	// 		return nil, err
	// 	}
	// }
	return api, nil
}
