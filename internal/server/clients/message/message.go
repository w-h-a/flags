package message

import "sync"

type Client interface {
	Send(diff Diff, waitGroup *sync.WaitGroup) error
}
