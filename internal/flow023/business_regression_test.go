package flow023

import (
	"path/filepath"
	"testing"
)

func Test484BusinessRegression(t *testing.T) {
	scenario, err := NewScenario(filepath.Join(t.TempDir(), "stories.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer scenario.Close()
	records, err := scenario.CreateTwoStories()
	if err != nil {
		t.Fatal(err)
	}
	first, second, err := scenario.UpdateAmounts(records[0].ID, records[1].ID, 210, 640)
	if err != nil {
		t.Fatal(err)
	}
	if first.Amount != 210 {
		t.Fatalf("first amount=%d", first.Amount)
	}
	if second.Amount != 640 {
		t.Fatalf("second amount=%d", second.Amount)
	}
}
