package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"

	files "github.com/ipfs/go-ipfs-files"
	u "github.com/ipfs/go-ipfs-util"
	importer "github.com/ipfs/go-unixfs/importer"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"

	chunker "github.com/ipfs/go-ipfs-chunker"
	ipld "github.com/ipfs/go-ipld-format"
)

// NodeFromBufferReader returns IPLD node of the file
func NodeFromBufferReader(api coreiface.CoreAPI) ipld.Node {
	data := make([]byte, 42)
	u.NewTimeSeededRand().Read(data)
	r := bytes.NewReader(data)
	nd, err := importer.BuildDagFromReader(api.Dag(), chunker.DefaultSplitter(r))
	if err != nil {
		log.Fatal(err)
	}
	return nd
}

// NodeFromDirectory example of adding virtual folder to IPFS
func NodeFromDirectory(api coreiface.CoreAPI) ipld.Node {
	mapFiles := map[string]files.Node{
		"one": files.NewBytesFile([]byte("testfileA")),
		"two": files.NewBytesFile([]byte("testfileB")),
	}
	dir := files.NewMapDirectory(mapFiles)
	_ = dir.(files.Directory)

	// nd, err := fileAdder.AddAllAndPin(dir)
	k, err := api.Unixfs().Add(context.Background(), dir.(files.Directory), options.Unixfs.Wrap(false))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("k=", k.String())
	return nil
}

func main() {

	coreAPI, err := MakeAPISwarm(context.Background(), false)
	if err != nil {
		log.Fatal("Error is ", err)
	}

	nd := NodeFromDirectory(coreAPI)
	if nd != nil {
		fmt.Println(" NODE=", nd.Cid().String())
		fmt.Println("  DBG=", GetNodeDataString(coreAPI.Dag(), nd))
		fmt.Println("  RAW=", hex.EncodeToString(nd.RawData()))
	}
}
