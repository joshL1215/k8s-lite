package apiserver

import (
	"log"
	"sync"

	"github.com/joshL1215/k8s-lite/internal/api/models"
)

type watchManager struct {
	mu       sync.Mutex
	watchers map[string][]chan WatchEvent
}

type WatchEvent struct {
	Type string       `json:"eventType"`
	Pod  *models.Pod  `json:"pod,omitempty"`
	Node *models.Node `json:"node,omitempty"`
}

func (wm *watchManager) Publish(namespace string, event WatchEvent) {
	wm.mu.Lock()
	subscribers := append([]chan WatchEvent(nil), wm.watchers[namespace]...) // good to use this mutex pattern since the entire function could be quite slow
	wm.mu.Unlock()

	for _, ch := range subscribers {
		select {
		case ch <- event:
		default:
			log.Printf("Dropped watch event for namespace %s due to slow consumer", namespace)
		}
	}
}

func (wm *watchManager) Subscribe(namespace string) chan WatchEvent {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	ch := make(chan WatchEvent, 100)
	wm.watchers[namespace] = append(wm.watchers[namespace], ch)
	return ch
}

func (wm *watchManager) Unsubscribe(namespace string, ch chan WatchEvent) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	watchers := wm.watchers[namespace]
	for i := range watchers {
		if watchers[i] == ch {
			wm.watchers[namespace] = append(watchers[:i], watchers[i+1:]...)
			close(ch)
			break
		}
	}
}
