# SimpleSearcher
Given a list of urls in url.txt, write a program that will fetch each page and determine whether a search term exists on the page (this search can be a really rudimentary regex - this part isn't too important).

## Constraints

* Search is case insensitive
* Should be concurrent
* But! It shouldn't have more than 20 HTTP requests at any given time.
* The results should be writted out to a file results.txt
* The goal is to avoid using thread pooling libraries like ThreadPoolExecutor or Celluloid
