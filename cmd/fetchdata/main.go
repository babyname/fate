package main

import (
	"archive/zip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const dataDir = "data/raw"

var sources = map[string]struct {
	URL         string
	Description string
	Format      string
	Extract     bool
}{
	"unihan": {
		URL:         "https://www.unicode.org/Public/UCD/latest/ucd/Unihan.zip",
		Description: "Unicode Unihan Database (all data, zip archive)",
		Format:      "zip",
		Extract:     true,
	},
	"hanzi-wuxing": {
		URL:         "https://raw.githubusercontent.com/mozillazg/hanzi-wuxing/master/data.json",
		Description: "mozillazg/hanzi-wuxing (汉字五行数据)",
		Format:      "json",
		Extract:     false,
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
	case "scrape-wuxing":
		scrapeWuxing()
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

	var outputPath string
	if src.Format == "zip" {
		outputPath = filepath.Join(dataDir, name+".zip")
	} else if src.Format == "json" {
		outputPath = filepath.Join(dataDir, name+".json")
	} else {
		outputPath = filepath.Join(dataDir, name+".txt")
	}

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

	if src.Extract && src.Format == "zip" {
		fmt.Printf("  Extracting zip...\n")
		extractDir := filepath.Join(dataDir, name)
		if err := os.MkdirAll(extractDir, 0755); err != nil {
			fmt.Printf("  Error creating extract dir: %v\n", err)
			return
		}
		if err := unzip(outputPath, extractDir); err != nil {
			fmt.Printf("  Error extracting: %v\n", err)
		} else {
			fmt.Printf("  Extracted to: %s\n", extractDir)
		}
	}
}

func fetchAll() {
	for name := range sources {
		fetchSource(name)
		fmt.Println()
	}
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, 0755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func showStatus() {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Data status:")
	fmt.Println()

	for name, src := range sources {
		var outputPath string
		if src.Format == "zip" {
			outputPath = filepath.Join(dataDir, name+".zip")
		} else if src.Format == "json" {
			outputPath = filepath.Join(dataDir, name+".json")
		} else {
			outputPath = filepath.Join(dataDir, name+".txt")
		}

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

func scrapeWuxing() {
	elements := []struct {
		pinyin string
		label  string
	}{
		{"jin", "金"}, {"mu", "木"}, {"shui", "水"}, {"huo", "火"}, {"tu", "土"},
	}

	wuxingData := make(map[string]string)
	charInSpan := regexp.MustCompile(`<span>([一-龥])</span>`)
	jar, _ := cookiejar.New(nil)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}
	client := &http.Client{
		Timeout:   30 * time.Second,
		Jar:       jar,
		Transport: transport,
	}

	visitReq, _ := http.NewRequest("GET", "https://www.yw11.com/zidian/", nil)
	visitReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	visitResp, err := client.Do(visitReq)
	if err != nil {
		fmt.Printf("Warning: failed to visit homepage for cookies: %v\n", err)
	} else {
		visitResp.Body.Close()
		fmt.Printf("Homepage visited, cookies acquired\n")
	}

	for _, elem := range elements {
		fmt.Printf("Scraping %s (%s)...\n", elem.label, elem.pinyin)
		for page := 1; page <= 50; page++ {
			url := fmt.Sprintf("https://www.yw11.com/zidian/wuxing/%s/", elem.pinyin)
			if page > 1 {
				url = fmt.Sprintf("https://www.yw11.com/zidian/wuxing/%s/%d/", elem.pinyin, page)
			}

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Printf("  Page %d: error creating request: %v\n", page, err)
				continue
			}
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
			req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
			req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
			req.Header.Set("Referer", "https://www.yw11.com/zidian/")
			req.Header.Set("Connection", "keep-alive")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("  Page %d: error: %v\n", page, err)
				continue
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				fmt.Printf("  Page %d: read error: %v\n", page, err)
				continue
			}

			if resp.StatusCode == http.StatusForbidden {
				fmt.Printf("  Page %d: HTTP 403, retrying after delay...\n", page)
				time.Sleep(3 * time.Second)
				req2, _ := http.NewRequest("GET", url, nil)
				req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
				req2.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
				req2.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
				req2.Header.Set("Referer", "https://www.yw11.com/zidian/")
				req2.Header.Set("Connection", "keep-alive")
				resp2, err := client.Do(req2)
				if err != nil || resp2.StatusCode != http.StatusOK {
					if resp2 != nil {
						resp2.Body.Close()
					}
					fmt.Printf("  Page %d: retry failed, stopping\n", page)
					break
				}
				resp = resp2
				body, err = io.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil {
					fmt.Printf("  Page %d: read error: %v\n", page, err)
					break
				}
			} else if resp.StatusCode != http.StatusOK {
				fmt.Printf("  Page %d: HTTP %d, stopping\n", page, resp.StatusCode)
				break
			}

			content := string(body)

			count := 0
			matches := charInSpan.FindAllStringSubmatch(content, -1)
			for _, m := range matches {
				char := strings.TrimSpace(m[1])
				if _, exists := wuxingData[char]; !exists {
					wuxingData[char] = elem.label
					count++
				}
			}

			fmt.Printf("  Page %d: %d new chars\n", page, count)
			if count == 0 {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Printf("\nTotal: %d characters with WuXing data\n", len(wuxingData))

	result := make(map[string]interface{})
	result["source"] = "yw11.com"
	result["count"] = len(wuxingData)
	result["data"] = wuxingData

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		os.Exit(1)
	}

	outputPath := filepath.Join(dataDir, "wuxing.json")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Printf("Error creating dir: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		fmt.Printf("Error writing: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Saved to: %s\n", outputPath)
}
