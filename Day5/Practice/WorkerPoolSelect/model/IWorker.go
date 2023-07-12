package model

import "sync"

type Worker interface {
	Worker(job int, wg *sync.WaitGroup)
}
