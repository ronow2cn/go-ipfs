package coreunix

import (
	"context"

	core "github.com/ipfs/go-ipfs/core"
	path "gx/ipfs/QmXVEFLxXFzdV1HKRrZA7scrnsNhjHRSygPbboTZ4H7Tij/go-path"
	resolver "gx/ipfs/QmXVEFLxXFzdV1HKRrZA7scrnsNhjHRSygPbboTZ4H7Tij/go-path/resolver"
	uio "gx/ipfs/QmcLmU3rW9jmxPf2SZuvS15Dhtji4TivmuEca58HTdPRBm/go-unixfs/io"
)

func Cat(ctx context.Context, n *core.IpfsNode, pstr string) (uio.DagReader, error) {
	r := &resolver.Resolver{
		DAG:         n.DAG,
		ResolveOnce: uio.ResolveUnixfsOnce,
	}

	dagNode, err := core.Resolve(ctx, n.Namesys, r, path.Path(pstr))
	if err != nil {
		return nil, err
	}

	return uio.NewDagReader(ctx, dagNode, n.DAG)
}
