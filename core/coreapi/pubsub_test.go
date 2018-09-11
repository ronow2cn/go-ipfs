package coreapi_test

import (
	"context"
	"testing"
	"time"
)

func TestBasicPubSub(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nds, apis, err := makeAPISwarm(ctx, true, 2)
	if err != nil {
		t.Fatal(err)
	}

	sub, err := apis[0].PubSub().Subscribe(ctx, "testch")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		tick := time.Tick(100 * time.Millisecond)

		for {
			err = apis[1].PubSub().Publish(ctx, "testch", []byte("hello world"))
			if err != nil {
				t.Fatal(err)
			}
			select {
			case <-tick:
			case <-ctx.Done():
				return
			}
		}
	}()

	m, err := sub.Next(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if string(m.Data()) != "hello world" {
		t.Errorf("got invalid data: %s", string(m.Data()))
	}

	if m.From() != nds[1].Identity {
		t.Errorf("m.From didn't match")
	}
}
