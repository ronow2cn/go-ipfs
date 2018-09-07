package coredag

import (
	"io"
	"io/ioutil"

	ipld "gx/ipfs/QmRurKMTJEe88d8LQHeDDqc1sBf4wZVJ8PvjVqu6gEw5ee/go-ipld-format"
	ipldcbor "gx/ipfs/QmYXN9xmo85th1SxePP6sKNkQBPzcoppTL7RRcW2LBtW15/go-ipld-cbor"
)

func cborJSONParser(r io.Reader, mhType uint64, mhLen int) ([]ipld.Node, error) {
	nd, err := ipldcbor.FromJson(r, mhType, mhLen)
	if err != nil {
		return nil, err
	}

	return []ipld.Node{nd}, nil
}

func cborRawParser(r io.Reader, mhType uint64, mhLen int) ([]ipld.Node, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	nd, err := ipldcbor.Decode(data, mhType, mhLen)
	if err != nil {
		return nil, err
	}

	return []ipld.Node{nd}, nil
}
