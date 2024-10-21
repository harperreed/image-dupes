package main

import (
	"fmt"
	"sync"
)

type Progress struct {
	mu             sync.Mutex
	totalFiles     int
	processedFiles int
}

func (p *Progress) Increment() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.processedFiles++
}

func (p *Progress) GetProgress() (int, int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.processedFiles, p.totalFiles
}

func (p *Progress) DisplayProgress() {
	processed, total := p.GetProgress()
	fmt.Printf("\rProcessed %d/%d files\n", processed, total)
}
