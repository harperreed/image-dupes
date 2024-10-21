package main

import (
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

// func displayProgress(progress *Progress) {
// 	for {
// 		time.Sleep(500 * time.Millisecond)
// 		processed, total := progress.GetProgress()
// 		if processed >= total {
// 			return
// 		}
// 		fmt.Printf("\rProcessed %d/%d images", processed, total)
// 	}
// }
