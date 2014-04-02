package waitout_test

import (
	"testing"

	"github.com/facebookgo/waitout"
)

func TestWait(t *testing.T) {
	t.Parallel()
	expected := []byte("42")
	w := waitout.New(expected)
	sent := make(chan struct{})
	go func() {
		defer close(sent)
		w.Write(expected)
	}()
	w.Wait()
	<-sent
}

func TestWaitSplit(t *testing.T) {
	t.Parallel()
	expected := []byte("42")
	w := waitout.New(expected)
	sent := make(chan struct{})
	go func() {
		defer close(sent)
		w.Write([]byte("4"))
		w.Write([]byte("2"))
	}()
	w.Wait()
	<-sent
}

func TestTwice(t *testing.T) {
	t.Parallel()
	expected := []byte("42")
	w := waitout.New(expected)
	sent := make(chan struct{})
	go func() {
		defer close(sent)
		w.Write(expected)
		w.Write(expected)
	}()
	w.Wait()
	<-sent
}

func TestManyInParallel(t *testing.T) {
	t.Parallel()
	const count = 100
	expected := []byte("42")
	w := waitout.New(expected)

	event := make(chan bool)
	for i := 0; i < count; i++ {
		go func() {
			if r := recover(); r != nil {
				t.Fatal(r)
			}
			w.Write(expected)
			event <- true
		}()
	}

	for i := 0; i < count; i++ {
		go func() {
			if r := recover(); r != nil {
				t.Fatal(r)
			}
			w.Wait()
			event <- true
		}()
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < count*2; i++ {
			<-event
		}
	}()

	<-done
	w.Wait()
}
