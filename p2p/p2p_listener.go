package p2p

import (
	"errors"
	"sync"

	net "gx/ipfs/QmQSbtGXCyNrj34LWL8EgXyNNYDZ8r3SwQcpW5pPxVhLnM/go-libp2p-net"
	peer "gx/ipfs/QmQsErDt8Qgw1XrsXf2BpEzDgGWtB1YLsTAARBup5b6B9W/go-libp2p-peer"
	ma "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	p2phost "gx/ipfs/QmfH9FKYv3Jp1xiyL8sPchGBUBg6JA6XviwajAo3qgnT3B/go-libp2p-host"
)

// Listener listens for connections and proxies them to a target
type P2PListener interface {
	Protocol() protocol.ID
	ListenAddress() ma.Multiaddr
	TargetAddress() ma.Multiaddr

	start() error

	// Close closes the listener. Does not affect child streams
	Close() error
}

// ListenerRegistry is a collection of local application proto listeners.
type ListenersP2P struct {
	sync.RWMutex

	Listeners map[protocol.ID]ListenerLocal
	starting  map[protocol.ID]struct{}
}

func newListenerP2PRegistry(id peer.ID, host p2phost.Host) *ListenersP2P {
	reg := &ListenersP2P{
		Listeners: map[protocol.ID]ListenerLocal{},
		starting:  map[protocol.ID]struct{}{},
	}

	addr, err := ma.NewMultiaddr(maPrefix + id.Pretty())
	if err != nil {
		panic(err)
	}

	host.SetStreamHandlerMatch("/x/", func(p string) bool {
		reg.RLock()
		defer reg.RUnlock()

		for _, l := range reg.Listeners {
			if l.ListenAddress().Equal(addr) && string(l.Protocol()) == p {
				return true
			}
		}

		return false
	}, func(stream net.Stream) {
		reg.RLock()
		defer reg.RUnlock()

		for _, l := range reg.Listeners {
			if l.ListenAddress().Equal(addr) && l.Protocol() == stream.Protocol() {
				go l.(*remoteListener).handleStream(stream)
				return
			}
		}
	})

	return reg
}

// Register registers listenerInfo into this registry and starts it
func (r *ListenersP2P) Register(l ListenerLocal) error {
	r.Lock()

	k := l.Protocol()
	if _, ok := r.Listeners[k]; ok {
		r.Unlock()
		return errors.New("listener already registered")
	}

	r.Listeners[k] = l
	r.starting[k] = struct{}{}

	r.Unlock()

	err := l.start()

	r.Lock()
	defer r.Unlock()

	delete(r.starting, k)

	if err != nil {
		delete(r.Listeners, k)
		return err
	}

	return nil
}

// Deregister removes p2p listener from this registry
func (r *ListenersP2P) Deregister(k protocol.ID) (bool, error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.starting[k]; ok {
		return false, errors.New("listener didn't start yet")
	}

	_, ok := r.Listeners[k]
	delete(r.Listeners, k)
	return ok, nil
}