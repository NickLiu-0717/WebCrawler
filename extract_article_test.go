package main

import (
	"strings"
	"testing"
)

func TestExtarctArticle(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expectedTitle string
		containString string
	}{
		{
			name:          "BBC news 1",
			inputURL:      "https://www.bbc.com/news/articles/cp9x97yvjp4o",
			expectedTitle: "Why Saudi Arabia is the venue of choice for Trump talks on Ukraine",
			containString: "The Saudis made clear that they would follow what they perceive",
		},
		{
			name:          "CNN news 1",
			inputURL:      "https://edition.cnn.com/2025/02/16/china/china-military-readiness-rand-report-intl-hnk-ml/index.html",
			expectedTitle: "Is China’s military really built for war? New report questions Beijing’s arms buildup",
			containString: "Simulations by US defense experts have repeatedly shown the US",
		},
		{
			name:          "FORBES news 1",
			inputURL:      "https://www.forbes.com/sites/siladityaray/2025/02/17/doge-is-seeking-access-to-critical-irs-system-that-holds-taxpayer-data-heres-what-to-know/",
			expectedTitle: "DOGE Is Seeking Access To Critical IRS System That Holds Taxpayer Data—Here’s What To Know",
			containString: "DOGE’s reported bid to gain access to the IRS systems comes just days",
		},
	}
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			htmlbody, err := getHTML(tc.inputURL)
			if err != nil {
				t.Fatal(err)
			}
			gottitle, gotcontent, err := extractArticles(htmlbody)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(gottitle, tc.expectedTitle) {
				t.Errorf("Test %d - %v FAIL: expected: %s, got: %s", i, tc.name, tc.expectedTitle, gottitle)
			}
			if !strings.Contains(gotcontent, tc.containString) {
				t.Errorf("Test %d - %v FAIL: content not contain expected string %s", i, tc.name, tc.containString)
			}
			// if i == 2 {
			// 	fmt.Println("Title:", gottitle)
			// 	fmt.Println("\nContent:\n", gotcontent)
			// }
		})
	}

}
