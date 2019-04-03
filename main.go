package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"

	ds "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	files "github.com/ipfs/go-ipfs-files"
	u "github.com/ipfs/go-ipfs-util"
	util "github.com/ipfs/go-ipfs-util"
	coreunix "github.com/ipfs/go-ipfs/core/coreunix"
	pin "github.com/ipfs/go-ipfs/pin"
	mdag "github.com/ipfs/go-merkledag"
	merkledag "github.com/ipfs/go-merkledag"
	importer "github.com/ipfs/go-unixfs/importer"

	bserv "github.com/ipfs/go-blockservice"
	dssync "github.com/ipfs/go-datastore/sync"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	chunker "github.com/ipfs/go-ipfs-chunker"
	ipld "github.com/ipfs/go-ipld-format"
)

// IpfsServices is a mini-version for IpfsNode
type IpfsServices struct {
	DAG        ipld.DAGService
	BlockStore bstore.Blockstore
	DB         ds.Datastore
	Blockserv  bserv.BlockService
	Pinner     pin.Pinner
}

// NewMemoryServices creates instance of IpfsServices for in-memory testing
func NewMemoryServices() *IpfsServices {
	dataStore := ds.NewMapDatastore()
	db := dssync.MutexWrap(dataStore)
	bs := bstore.NewBlockstore(db)
	blockserv := bserv.New(bs, offline.Exchange(bs))
	dag := merkledag.NewDAGService(blockserv)
	pinner := pin.NewPinner(db, dag, dag)

	// var pinning pin.Pinner = pin.NewPinner()
	// var blockstore bstore.GCBlockstore = bstore.NewBlockstore()
	return &IpfsServices{
		DAG:        dag,
		BlockStore: bs,
		DB:         db,
		Blockserv:  blockserv,
		Pinner:     pinner,
	}
}

func nodeFromBufferReader(dagService ipld.DAGService) ipld.Node {
	data := make([]byte, 42)
	u.NewTimeSeededRand().Read(data)
	r := bytes.NewReader(data)
	nd, err := importer.BuildDagFromReader(dagService, chunker.DefaultSplitter(r))
	if err != nil {
		log.Fatal(err)
	}
	return nd
}

var rand = util.NewTimeSeededRand()

func randNode() ipld.Node {
	nd := new(mdag.ProtoNode)
	nd.SetData(make([]byte, 32))
	rand.Read(nd.Data())
	return nd
}

func nodeFromFileAdder(srv *IpfsServices) ipld.Node {
	data := make([]byte, 40)
	u.NewTimeSeededRand().Read(data)
	f := files.NewBytesFile(data)

	fileAdder, err := coreunix.NewAdder(context.Background(),
		srv.Pinner, blockstore.NewGCLocker(), srv.DAG)
	nd, err := fileAdder.AddAllAndPin(f)
	if err != nil {
		log.Fatal(err)
	}
	return nd
}

// This method of adding folder seems to be wrong
// and unreliable. On output, instead of directory node,
// we have the node of the last added file, which is confusing
func nodeFromDirectory(srv *IpfsServices) ipld.Node {

	mapFiles := map[string]files.Node{
		"one": files.NewBytesFile([]byte("testfileA")),
		"two": files.NewBytesFile([]byte("testfileB")),
	}
	dir := files.NewMapDirectory(mapFiles)
	fileAdder, err := coreunix.NewAdder(context.Background(),
		srv.Pinner, blockstore.NewGCLocker(), srv.DAG)

	nd, err := fileAdder.AddAllAndPin(dir)
	if err != nil {
		log.Fatal(err)
	}
	return nd
}

func main() {
	srv := NewMemoryServices()
	// nd := nodeFromBufferReader(srv.DAG)
	// nd := nodeFromFileAdder(srv)
	// nd := randNode()

	nd := nodeFromDirectory(srv)

	fmt.Println(" NODE=", nd.Cid().String())
	fmt.Println("  DBG=", GetNodeDataString(srv.DAG, nd))
	fmt.Println("  RAW=", hex.EncodeToString(nd.RawData()))
}
