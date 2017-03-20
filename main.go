package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	infile     = "./url.txt"
	outfile    = "./results.txt"
	maxThreads = 20
)

func main() {
	// parse command line arguments:
	if len(os.Args) != 2 {
		usage()
		return
	}
	searchTerm := os.Args[1]

	// open and read input file
	lines, err := readLinesFromFiles(infile)
	if err != nil {
		fmt.Println("Could not read URL File:", err)
		return
	}

	// parse all the urls from the lines of the input file (second argument in CSV):
	var urls []string
	for i, l := range lines {
		if i == 0 {
			continue // skip the first line
		}
		v := strings.Split(l, ",")
		if len(v) != 6 {
			fmt.Println("Could not parse line in URL file:", i)
			return
		}
		u := v[1]
		if strings.HasPrefix(u, "\"") && strings.HasSuffix(u, "\"") {
			u = u[1 : len(u)-1]
		}
		urls = append(urls, u)
	}

	// open the outputfile. truncate and append.
	out, err := os.Create(outfile)
	if err != nil {
		fmt.Println("Could not open output file:", err)
		return
	}
	defer out.Close()

	// now process the search
	sr := parallelSearch(urls, searchTerm, maxThreads)

	// and output the results
	bufout := bufio.NewWriter(out)
	for _, r := range sr {
		if r.Error == nil {
			bufout.WriteString("\"" + r.URL + "\"," + strconv.Itoa(r.Count) + ",\n")
		} else {
			bufout.WriteString("\"" + r.URL + "\"," + strconv.Itoa(r.Count) + "," + r.Error.Error() + "\n")
		}
	}
	err = bufout.Flush()
	if err != nil {
		fmt.Println("Problem writing to output file:", err)
		return
	}
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("\t" + os.Args[0] + " SEARCH_TERM")
}

// opens a file and returns the entire contents into an array split by lines
func readLinesFromFiles(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

type SearchResult struct {
	URL   string
	Count int
	Error error
}

// Fetches all the given urls, up to maxThreads in parallel,
// and processes them.
// FIXME: searchTerm is simply added naively to a regexp, so it should not contain any regexp characters
func parallelSearch(urls []string, searchTerm string, maxThreads int) []SearchResult {
	fanout := make(chan string, maxThreads)
	fanin := make(chan SearchResult)
	done := make(chan []SearchResult)

	regexp := regexp.MustCompile("(?i)\\b+" + searchTerm + "\\b+")

	// This thread collects data from the workers
	go func() {
		i := 0
		sr := make([]SearchResult, 0, len(urls))
		for r := range fanin {
			//fmt.Println( r.URL, r.Count, r.Error )
			i++
			sr = append(sr, r)
			if i == len(urls) {
				break
			}
		}
		done <- sr
	}()

	// create worker threads:
	for i := 0; i < maxThreads; i++ {
		go func() {
			for u := range fanout {
				fanin <- performSearch(u, regexp)
			}
		}()
	}
	// push the data to the workers:
	for _, u := range urls {
		fanout <- "http://" + u //FIXME: this is not the best place to add http://
	}
	close(fanout)

	// wait for completion and get results
	sr := <-done

	// close channels
	close(fanin)
	close(done)

	return sr
}

// fetches the given URL and finds the given regular expression on it.
func performSearch(url string, re *regexp.Regexp) SearchResult {
	response, err := http.Get(url)
	if err != nil {
		return SearchResult{
			URL:   url,
			Count: -1,
			Error: err,
		}
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return SearchResult{
			URL:   url,
			Count: -1,
			Error: errors.New("Unexpected Status Code: " + strconv.Itoa(response.StatusCode)),
		}
	}

	d, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return SearchResult{
			URL:   url,
			Count: -1,
			Error: err,
		}
	}

	count := 0
	found := re.FindAllIndex(d, -1)
	//fmt.Println( found )
	if found == nil {
		count = 0
	} else {
		count = len(found)
	}

	return SearchResult{
		URL:   url,
		Count: count,
		Error: nil,
	}
}
