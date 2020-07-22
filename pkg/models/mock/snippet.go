package mock

import (
	"time"

	"herobix.com/snipetbox/pkg/models"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An young silent pond",
	Content: "An young silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

// SnippetModel is a mock structure for testing snippet struct
type SnippetModel struct{}

// Get is a mock for testing Get Snippet function
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

// Latest is a mock for testing Latest Snippet function
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}

// Insert is a mock for testing Latest Snippet function
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	switch title {
	case "working":
		return 1, nil
	default:
		return 0, nil
	}
}
