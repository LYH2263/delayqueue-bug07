package delayqueue

import (
	"context"
	"testing"
	"time"
)

func TestBug07_EnqueueIgnoresCanceledCtx(t *testing.T) {
	b, _ := New(Options{PollInterval: time.Millisecond})
	defer b.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := b.Enqueue(ctx, "t1", "q", []byte("p")); err == nil {
		t.Fatal("Enqueue should fail on canceled ctx")
	}
	if _, ok := b.tasks["t1"]; ok {
		t.Fatal("canceled Enqueue must not insert")
	}
	_ = b.Enqueue(context.Background(), "t2", "q", []byte("p"))
	dctx, dcancel := context.WithCancel(context.Background())
	dcancel()
	_, err := b.Dequeue(dctx, "w")
	if err == nil {
		t.Fatal("Dequeue should fail on canceled ctx")
	}
}
