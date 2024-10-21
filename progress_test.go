package main

import (
	"bytes"
	"io"
	"os"
	"strings"
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
	p := &Progress{totalFiles: 10, processedFiles: 0}

	// Use a buffer to capture the output
	var buf bytes.Buffer

	// Use a channel to synchronize the goroutine and the main test thread
	done := make(chan bool)
	outputReady := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			p.Increment()
		}
		done <- true
	}()

	go func() {
		<-done
		// Redirect stdout to our buffer
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		p.DisplayProgress()

		w.Close()
		os.Stdout = old

		_, err := io.Copy(&buf, r)
		if err != nil {
			t.Errorf("Failed to copy output: %v", err)
		}
		outputReady <- true
	}()

	// Wait for the output to be ready
	<-outputReady

	output := buf.String()

	// Trim any leading/trailing whitespace and remove carriage returns
	output = strings.TrimSpace(strings.Replace(output, "\r", "", -1))

	expectedOutput := "Processed 10/10 files"
	if output != expectedOutput {
		t.Errorf("Expected output '%s', got '%s'", expectedOutput, output)
	}
}

func TestProgressThreadSafety(t *testing.T) {
	p := &Progress{totalFiles: 1000, processedFiles: 0}
	numGoroutines := 100
	incrementsPerGoroutine := 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				p.Increment()
				time.Sleep(time.Millisecond) // Small delay to increase chances of race condition
			}
		}()
	}

	wg.Wait()

	processed, total := p.GetProgress()
	expected := numGoroutines * incrementsPerGoroutine
	if processed != expected || total != 1000 {
		t.Errorf("Expected (%d, 1000), got (%d, %d)", expected, processed, total)
	}
}
