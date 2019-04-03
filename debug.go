package main

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"strconv"

	ipld "github.com/ipfs/go-ipld-format"
	mdag "github.com/ipfs/go-merkledag"
	unixfs "github.com/ipfs/go-unixfs"
	uio "github.com/ipfs/go-unixfs/io"
	log "github.com/sirupsen/logrus"
)

// GetNodeDataString should return some string reprecenging node data
func GetNodeDataString(dagService ipld.DAGService, ipldNode ipld.Node) string {

	protoNode := ipldNode.(*mdag.ProtoNode)
	ctx := context.Background()

	prefix := ""
	canRead := false
	switch ipldNode := ipldNode.(type) {
	case *mdag.RawNode:
		size := len(ipldNode.RawData())
		prefix = "*RAW.NODE, size=" + strconv.Itoa(size)

	case *mdag.ProtoNode:
		fsNode, err := unixfs.FSNodeFromBytes(ipldNode.Data())
		if err != nil {
			prefix = "*PROTO, err=" + err.Error()
		} else {
			switch fsNode.Type() {
			case unixfs.TFile:
				size := int(fsNode.FileSize())
				prefix = "*PROTO.NODE, TFile size=" + strconv.Itoa(size)
				canRead = true

			case unixfs.TRaw:
				size := int(fsNode.FileSize())
				prefix = "*PROTO.NODE, TRaw size=" + strconv.Itoa(size)
				canRead = true

			case unixfs.TDirectory:
				prefix = "*PROTO.NODE, TDirectory, children=" + strconv.Itoa(len(ipldNode.Links()))

			case unixfs.THAMTShard:
				prefix = "*PROTO.NODE, THAMTShard"

			case unixfs.TMetadata:
				prefix = "*PROTO.NODE, TMetadata, links=" + strconv.Itoa(len(ipldNode.Links()))
			case unixfs.TSymlink:
				prefix = "*PROTO.NODE, TSymLink"
			default:
				prefix = "*PROTO.NODE, Unrecognized"
			}
		}
	default:
		prefix = "*UNKNOWN.NODE"
	}

	body := ""
	if canRead {
		ndr, err := uio.NewDagReader(ctx, protoNode, dagService)
		if err != nil {
			log.Fatal(err)
		}
		out, err := ioutil.ReadAll(ndr)
		if err != nil {
			log.Fatal(err)
		}
		// _ = out
		body = hex.EncodeToString(out)
	}
	return prefix + " " + body
}
