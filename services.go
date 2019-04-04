package main

import (
	util "github.com/ipfs/go-ipfs-util"
	ipld "github.com/ipfs/go-ipld-format"
	mdag "github.com/ipfs/go-merkledag"
)

var rand = util.NewTimeSeededRand()

func RandNode() ipld.Node {
	nd := new(mdag.ProtoNode)
	nd.SetData(make([]byte, 32))
	rand.Read(nd.Data())
	return nd
}

// // IpfsServices was a mini-version for IpfsNode
// type IpfsServices struct {
// 	DAG        ipld.DAGService
// 	BlockStore bstore.Blockstore
// 	DB         ds.Datastore
// 	Blockserv  bserv.BlockService
// 	Pinner     pin.Pinner
// }

// // NewMemoryServices creates instance of IpfsServices for in-memory testing
// func NewMemoryServices() *IpfsServices {
// 	dataStore := ds.NewMapDatastore()
// 	db := dssync.MutexWrap(dataStore)
// 	bs := bstore.NewBlockstore(db)
// 	blockserv := bserv.New(bs, offline.Exchange(bs))
// 	dag := mdag.NewDAGService(blockserv)
// 	pinner := pin.NewPinner(db, dag, dag)

// 	// var pinning pin.Pinner = pin.NewPinner()
// 	// var blockstore bstore.GCBlockstore = bstore.NewBlockstore()
// 	return &IpfsServices{
// 		DAG:        dag,
// 		BlockStore: bs,
// 		DB:         db,
// 		Blockserv:  blockserv,
// 		Pinner:     pinner,
// 	}
// }

// func nodeFromFileAdder(api coreiface.CoreAPI) ipld.Node {
// 	data := make([]byte, 40)
// 	u.NewTimeSeededRand().Read(data)
// 	f := files.NewBytesFile(data)

// 	fileAdder, err := coreunix.NewAdder(context.Background(),
// 		api.Pin(), blockstore.NewGCLocker(), api.Dag())

// 	nd, err := fileAdder.AddAllAndPin(f)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return nd
// }

// This method of adding folder seems to be wrong
// and unreliable. On output, instead of directory node,
// we have the node of the last added file, which is confusing
// func nodeFromDirectory(api *coreiface.CoreAPI) ipld.Node {

// 	mapFiles := map[string]files.Node{
// 		"one": files.NewBytesFile([]byte("testfileA")),
// 		"two": files.NewBytesFile([]byte("testfileB")),
// 	}
// 	dir := files.NewMapDirectory(mapFiles)
// 	fileAdder, err := coreunix.NewAdder(context.Background(),
// 		srv.Pinner, blockstore.NewGCLocker(), srv.DAG)

// 	nd, err := fileAdder.AddAllAndPin(dir)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return nd
// }
