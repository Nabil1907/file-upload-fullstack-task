package progress

import (
	"sync"
)

// Progress represents the progress of a task.
type Progress struct {
	Processed int
	Total     int
	Status    string
}

// Tracker is a thread-safe progress tracker using sync.Map.
type Tracker struct {
	progressMap sync.Map
	listeners   sync.Map
}

// NewTracker creates a new progress tracker.
func NewTracker() *Tracker {
	return &Tracker{}
}

// UpdateProgress updates the progress for a specific task and notifies listeners.
func (t *Tracker) UpdateProgress(taskID string, processed int, total int, status string) {
	progress := Progress{
		Processed: processed,
		Total:     total,
		Status:    status,
	}

	// Store the progress
	t.progressMap.Store(taskID, progress)

	// Notify listeners
	if ch, exists := t.listeners.Load(taskID); exists {
		ch.(chan Progress) <- progress
	}
}

// GetProgress retrieves the progress for a specific task.
func (t *Tracker) GetProgress(taskID string) (Progress, bool) {
	value, exists := t.progressMap.Load(taskID)
	if !exists {
		return Progress{}, false
	}
	return value.(Progress), true
}

// DeleteProgress removes the progress for a specific task.
func (t *Tracker) DeleteProgress(taskID string) {
	t.progressMap.Delete(taskID)
	t.listeners.Delete(taskID)
}

// Subscribe creates a channel to listen for progress updates for a specific task.
func (t *Tracker) Subscribe(taskID string) chan Progress {
	ch := make(chan Progress)
	t.listeners.Store(taskID, ch)
	return ch
}

// Unsubscribe removes a listener for a specific task.
func (t *Tracker) Unsubscribe(taskID string) {
	if ch, exists := t.listeners.Load(taskID); exists {
		close(ch.(chan Progress))
		t.listeners.Delete(taskID)
	}
}
