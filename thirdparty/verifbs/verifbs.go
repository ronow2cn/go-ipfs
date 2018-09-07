package verifbs

import (
	"gx/ipfs/QmPM5YdzX6GaKJHPgCTRnMMyNc1k4YZ7KDAdhdbpG3mkhj/go-verifcid"
	bstore "gx/ipfs/QmSz8PDzfMyZZpd8KnV3ti8LjFngwTuVstmhXfvA93cYqG/go-ipfs-blockstore"
	blocks "gx/ipfs/QmVWmocy2WJk9ZvguJrNkPAkWiW5N5zdbaXrrvshZgsEFY/go-block-format"
	cid "gx/ipfs/QmcRoKTXnq18qQRZFa4jWwWvMQkxzWRpxCwcpCCFgnLUGi/go-cid"
)

type VerifBSGC struct {
	bstore.GCBlockstore
}

func (bs *VerifBSGC) Put(b blocks.Block) error {
	if err := verifcid.ValidateCid(b.Cid()); err != nil {
		return err
	}
	return bs.GCBlockstore.Put(b)
}

func (bs *VerifBSGC) PutMany(blks []blocks.Block) error {
	for _, b := range blks {
		if err := verifcid.ValidateCid(b.Cid()); err != nil {
			return err
		}
	}
	return bs.GCBlockstore.PutMany(blks)
}

func (bs *VerifBSGC) Get(c cid.Cid) (blocks.Block, error) {
	if err := verifcid.ValidateCid(c); err != nil {
		return nil, err
	}
	return bs.GCBlockstore.Get(c)
}

type VerifBS struct {
	bstore.Blockstore
}

func (bs *VerifBS) Put(b blocks.Block) error {
	if err := verifcid.ValidateCid(b.Cid()); err != nil {
		return err
	}
	return bs.Blockstore.Put(b)
}

func (bs *VerifBS) PutMany(blks []blocks.Block) error {
	for _, b := range blks {
		if err := verifcid.ValidateCid(b.Cid()); err != nil {
			return err
		}
	}
	return bs.Blockstore.PutMany(blks)
}

func (bs *VerifBS) Get(c cid.Cid) (blocks.Block, error) {
	if err := verifcid.ValidateCid(c); err != nil {
		return nil, err
	}
	return bs.Blockstore.Get(c)
}
