package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const dataDir = "data/raw"

var sources = map[string]struct {
	URL         string
	Description string
	Format      string
}{
	"unihan-irgs": {
		URL:         "https://www.unicode.org/Public/UCD/latest/ucd/Unihan_IRGSources.txt",
		Description: "Unicode Unihan IRG Sources (简繁映射、部首、笔画)",
		Format:      "tsv",
	},
	"unihan-readings": {
		URL:         "https://www.unicode.org/Public/UCD/latest/ucd/Unihan_Readings.txt",
		Description: "Unicode Unihan Readings (拼音、释义)",
		Format:      "tsv",
	},
	"unihan-other-mappings": {
		URL:         "https://www.unicode.org/Public/UCD/latest/ucd/Unihan_OtherMappings.txt",
		Description: "Unicode Unihan Other Mappings",
		Format:      "tsv",
	},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: fetchdata <command> [args]")
		fmt.Println("Commands:")
		fmt.Println("  list              List available data sources")
		fmt.Println("  fetch [name]      Download data source(s)")
		fmt.Println("  fetch-all         Download all data sources")
		fmt.Println("  status            Show download status")
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "list":
		listSources()
	case "fetch":
		if len(os.Args) < 3 {
			fmt.Println("Usage: fetchdata fetch <name>")
			os.Exit(1)
		}
		fetchSource(os.Args[2])
	case "fetch-all":
		fetchAll()
	case "status":
		showStatus()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func listSources() {
	fmt.Println("Available data sources:")
	fmt.Println()
	for name, src := range sources {
		fmt.Printf("  %-25s %s\n", name, src.Description)
		fmt.Printf("  %-25s %s\n", "", src.URL)
		fmt.Printf("  %-25s Format: %s\n", "", src.Format)
		fmt.Println()
	}
}

func fetchSource(name string) {
	src, ok := sources[name]
	if !ok {
		fmt.Printf("Unknown source: %s\n", name)
		fmt.Println("Use 'list' to see available sources")
		os.Exit(1)
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Printf("Error creating data dir: %v\n", err)
		os.Exit(1)
	}

	ext := ".txt"
	if src.Format == "json" {
		ext = ".json"
	}
	outputPath := filepath.Join(dataDir, name+ext)

	fmt.Printf("Fetching %s...\n", name)
	fmt.Printf("  URL: %s\n", src.URL)

	resp, err := http.Get(src.URL)
	if err != nil {
		fmt.Printf("Error fetching: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: HTTP %d\n", resp.StatusCode)
		os.Exit(1)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	written, err := io.Copy(f, resp.Body)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("  Saved: %s (%d bytes)\n", outputPath, written)
}

func fetchAll() {
	for name := range sources {
		fetchSource(name)
		fmt.Println()
	}
}

func showStatus() {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Data status:")
	fmt.Println()

	for name, src := range sources {
		ext := ".txt"
		if src.Format == "json" {
			ext = ".json"
		}
		outputPath := filepath.Join(dataDir, name+ext)

		info, err := os.Stat(outputPath)
		if err != nil {
			fmt.Printf("  %-25s NOT DOWNLOADED\n", name)
			continue
		}

		fmt.Printf("  %-25s %s (%d bytes)\n", name, info.ModTime().Format("2006-01-02 15:04"), info.Size())
	}
}

func generateKangxiStarterJSON() {
	starter := []map[string]interface{}{
		{
			"char":    "一",
			"radical":  "一",
			"strokes":  1,
			"wu_xing":  "土",
			"pinyin":   "yī",
			"meaning":  "数之始也",
		},
		{
			"char":    "丁",
			"radical":  "一",
			"strokes":  2,
			"wu_xing":  "火",
			"pinyin":   "dīng",
			"meaning":  "天干第四位",
		},
	}

	data, _ := json.MarshalIndent(starter, "", "  ")
	fmt.Println(string(data))
}

func generateWuxingStarterJSON() {
	starter := []map[string]interface{}{
		{"char": "木", "wu_xing": "木", "source": "self"},
		{"char": "林", "wu_xing": "木", "source": "self"},
		{"char": "森", "wu_xing": "木", "source": "self"},
		{"char": "火", "wu_xing": "火", "source": "self"},
		{"char": "炎", "wu_xing": "火", "source": "self"},
		{"char": "土", "wu_xing": "土", "source": "self"},
		{"char": "金", "wu_xing": "金", "source": "self"},
		{"char": "水", "wu_xing": "水", "source": "self"},
	}

	data, _ := json.MarshalIndent(starter, "", "  ")
	fmt.Println(string(data))
}

func generatePinyinStarterJSON() {
	starter := []map[string]interface{}{
		{"char": "中", "pinyin": []string{"zhōng", "zhòng"}, "source": "pinyin"},
		{"char": "华", "pinyin": []string{"huá", "huà"}, "source": "pinyin"},
		{"char": "明", "pinyin": []string{"míng"}, "source": "pinyin"},
	}

	data, _ := json.MarshalIndent(starter, "", "  ")
	fmt.Println(string(data))
}

func init() {
	_ = strings.TrimSpace
	_ = generateKangxiStarterJSON
	_ = generateWuxingStarterJSON
	_ = generatePinyinStarterJSON
}
