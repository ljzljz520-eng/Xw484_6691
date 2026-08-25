package flow023

import (
	"fmt"
	"path/filepath"

	"genealogy-story-organizer/internal/application"
	"genealogy-story-organizer/internal/domain"
	"genealogy-story-organizer/internal/store"
)

type Scenario struct {
	Service *application.Service
	Store   *store.Store
	Path    string
}

func NewScenario(path string) (*Scenario, error) {
	database, err := store.OpenWithClock(filepath.Clean(path), store.StaticClock{Value: "2025-01-02T03:04:05Z"})
	if err != nil {
		return nil, err
	}
	service, err := application.NewService(database, application.NewSequenceID(), store.StaticClock{Value: "2025-01-02T03:04:05Z"})
	if err != nil {
		_ = database.Close()
		return nil, err
	}
	return &Scenario{Service: service, Store: database, Path: filepath.Clean(path)}, nil
}

func (s *Scenario) Close() error { return s.Store.Close() }

func (s *Scenario) CreateTwoStories() ([]domain.Record, error) {
	first, err := s.Service.CreateStory("长子迁徙", "长子在战乱后迁往江南并留下家书。", "researcher", 180, []string{"迁徙"})
	if err != nil {
		return nil, err
	}
	second, err := s.Service.CreateStory("次子守业", "次子留在故乡守护祖屋和族谱。", "researcher", 420, []string{"守业"})
	if err != nil {
		return nil, err
	}
	return []domain.Record{first, second}, nil
}

func (s *Scenario) UpdateAmounts(firstID, secondID string, firstAmount, secondAmount int64) (domain.Record, domain.Record, error) {
	first, err := s.Service.UpdateAmount(firstID, firstAmount, "editor")
	if err != nil {
		return domain.Record{}, domain.Record{}, err
	}
	second, err := s.Service.UpdateAmount(secondID, secondAmount, "editor")
	if err != nil {
		return domain.Record{}, domain.Record{}, err
	}
	return first, second, nil
}

func (s *Scenario) Outcome(firstID, secondID string) (string, error) {
	first, err := s.Service.GetStory(firstID)
	if err != nil {
		return "", err
	}
	second, err := s.Service.GetStory(secondID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s=%d;%s=%d", first.ID, first.Amount, second.ID, second.Amount), nil
}
