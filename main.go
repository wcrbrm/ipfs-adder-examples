package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"

	u "github.com/ipfs/go-ipfs-util"
	importer "github.com/ipfs/go-unixfs/importer"
	coreiface "github.com/ipfs/interface-go-ipfs-core"

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

func main() {

	coreAPI, err := MakeAPISwarm(context.Background(), false)
	if err != nil {
		log.Fatal("Error is ", err)
	}

	nd := NodeFromBufferReader(coreAPI)
	fmt.Println(" NODE=", nd.Cid().String())
	fmt.Println("  DBG=", GetNodeDataString(coreAPI.Dag(), nd))
	fmt.Println("  RAW=", hex.EncodeToString(nd.RawData()))
}
