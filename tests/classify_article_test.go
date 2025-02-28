package main

import (
	"testing"

	"github.com/NickLiu-0717/crawler/service"
)

func TestClassifyArticle(t *testing.T) {
	tests := []struct {
		name       string
		inputTitle string
		expected   string
	}{
		{
			name:       "English title for politics",
			inputTitle: "Who was at the table at US-Russia talks in Saudi Arabia?",
			expected:   "politics",
		},
		{
			name:       "Chinese title for health",
			inputTitle: "研究揭清冠一號可抑制A型流感 降低發炎改善症狀 ｜ 公視新聞網 PNN",
			expected:   "health",
		},
	}
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotCatagory, err := service.ClassifyArticle(tc.inputTitle)
			if err != nil {
				t.Fatalf("Test %d %v FAIL - unexpected error: %v", i, tc.name, err)
			}
			if tc.expected != gotCatagory {
				t.Errorf("Test %d %v FAIL - wrong result: expected: %v, got: %v", i, tc.name, tc.expected, gotCatagory)
			}
		})
	}
}
