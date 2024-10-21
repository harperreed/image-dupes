package main

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"
	"time"
)

func TestProgressIncrement(t *testing.T) {
	p := &Progress{totalFiles: 10}

	p.Increment()
	if p.processedFiles != 1 {
		t.Errorf("Expected processedFiles to be 1, got %d", p.processedFiles)
	}

	p.Increment()
	if p.processedFiles != 2 {
		t.Errorf("Expected processedFiles to be 2, got %d", p.processedFiles)
	}
}

func TestProgressGetProgress(t *testing.T) {
	p := &Progress{totalFiles: 10, processedFiles: 5}

	processed, total := p.GetProgress()
	if processed != 5 || total != 10 {
		t.Errorf("Expected (5, 10), got (%d, %d)", processed, total)
	}
}

func TestProgressConcurrency(t *testing.T) {
	p := &Progress{totalFiles: 1000}
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.Increment()
		}()
	}

	wg.Wait()

	processed, total := p.GetProgress()
	if processed != 1000 || total != 1000 {
		t.Errorf("Expected (1000, 1000), got (%d, %d)", processed, total)
	}
}

func TestDisplayProgress(t *testing.T) {
	// Redirect stdout to capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	p := &Progress{totalFiles: 10, processedFiles: 0}

	// Mock time.Sleep
	oldSleep := time.Sleep
	time.Sleep = func(d time.Duration) {}
	defer func() { time.Sleep = oldSleep }()

	done := make(chan bool)
	go func() {
		displayProgress(p)
		done <- true
	}()

	// Simulate progress
	for i := 0; i < 10; i++ {
		p.Increment()
		if i == 9 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Wait for displayProgress to finish
	<-done

	// Restore stdout
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expectedOutput := "\rProcessed 9/10 images"
	if output != expectedOutput {
		t.Errorf("Expected output '%s', got '%s'", expectedOutput, output)
	}
}
