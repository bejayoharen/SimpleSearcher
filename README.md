# SimpleSearcher
Given a list of urls in url.txt, write a program that will fetch each page and determine whether a search term exists on the page (this search can be a really rudimentary regex - this part isn't too important).

## Constraints

* Search is case insensitive
* Should be concurrent
* But! It shouldn't have more than 20 HTTP requests at any given time.
* The results should be writted out to a file results.txt
* The goal is to avoid using thread pooling libraries like ThreadPoolExecutor or Celluloid

## Notes

* The implementation is entirely in one file (main.go) to make it easy to fetch, run and test.
* I did not write any tests, but there's plenty to test here.
* The search term is given as the sole commandline argument and the in and out files are given as consts in the code.
* For input, I used the format of the file in the example: CSV with 6 colums and the URL in the second column (all other columns are ignored). I added http:// to the beginning of each URL, and first line is skipped.
* The result file is CSV with three columns: The URL, the # of matches (or -1 on error), any error messages.
* The results may not be in the same order as the input file.
* I don't implement any sort of retries, nor do I modify any of the default timeouts or other properties.
* The complete input file is read into memory at once and the output is read once after all data is colleted. If the files were huge, it would be better to read them a few hundred at a time, but I'm guessing that wasn't what was being tested.
* I only used Go's built-in packages and datastructures.
* There's no console output except if there's an error, and it doesn't write to the file until it's done, so it might seem like it's doing nothing for a while.a Uncomment the code at line 122 if you want to see what it's doing a bit.
* The searchterm passed in as a commandline argument is just inserted naively into a regex. Obviously, this could cause problems if, eg, the searchterm has special characters.

## Compile and run

It should be easy to compile and run with:

    go build
    ./SimpleSearcher SEARCHTERM
