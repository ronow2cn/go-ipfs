package coredag

import (
	"io"
	"io/ioutil"
	"math"

	"gx/ipfs/QmSRe5UvVPLJ6LAtVH9gQZwEL4nbck5b5zNe4MChh3LJHk/go-merkledag"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	ipld "gx/ipfs/QmRurKMTJEe88d8LQHeDDqc1sBf4wZVJ8PvjVqu6gEw5ee/go-ipld-format"
	block "gx/ipfs/QmVWmocy2WJk9ZvguJrNkPAkWiW5N5zdbaXrrvshZgsEFY/go-block-format"
	cid "gx/ipfs/QmcRoKTXnq18qQRZFa4jWwWvMQkxzWRpxCwcpCCFgnLUGi/go-cid"
)

func rawRawParser(r io.Reader, mhType uint64, mhLen int) ([]ipld.Node, error) {
	if mhType == math.MaxUint64 {
		mhType = mh.SHA2_256
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	h, err := mh.Sum(data, mhType, mhLen)
	if err != nil {
		return nil, err
	}
	c := cid.NewCidV1(cid.Raw, h)
	blk, err := block.NewBlockWithCid(data, c)
	if err != nil {
		return nil, err
	}
	nd := &merkledag.RawNode{Block: blk}
	return []ipld.Node{nd}, nil
}
