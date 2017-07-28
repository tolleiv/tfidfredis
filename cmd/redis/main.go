package main

import (
	"log"
	"github.com/tolleiv/tfidfredis"
)

func main() {

	indexes := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}

	r, _ := tfidfredis.NewRepo("localhost:6379")

	for _, i := range indexes {
		r.Del(i)
	}
	r.Add("http://www.example.org", "A")
	r.Add("http://www.example.org", "A")
	r.Add("http://www.example.org", "A")
	r.Add("http://example.org", "A")
	r.Add("http://www.example.org", "A")
	r.Add("http://www.example.org", "B")
	r.Add("http://www.example.org", "B")
	r.Add("http://example.org", "B")
	r.Add("http://example.org", "B")
	r.Add("http://www.example.org", "C")
	r.Add("http://www.example.org", "C")
	r.Add("http://www.example.org", "C")
	r.Add("http://example.org", "C")
	r.Add("http://example.org", "C")
	r.Add("http://example.org", "C")
	r.Add("http://de.example.org", "C")
	r.Add("http://example.org", "C")
	r.Add("http://de.example.org", "D")
	r.Add("http://it.example.org", "D")
	r.Add("http://example.org", "D")
	r.Add("http://www.example.org", "E")
	r.Add("http://de.example.org", "E")
	r.Add("http://it.example.org", "E")
	r.Add("http://www.example.org", "F")
	for _, i := range indexes {
		r.Add("http://wellknown.example.org", i)
	}

	for _, t := range []string{"http://www.example.org", "http://example.org"} {
		for _, i := range indexes {
			c := r.Containing(t, indexes)
			tf := r.Tf(t, i)
			idf := r.Idf(t, indexes)
			tfidf := r.Tfidf(t, i, indexes)
			log.Printf("Result (%30s => %s) [%f]: TF %.4f IDF %.4f TFIDF %.4f \n", t, i, c, tf, idf, tfidf)
		}
	}

	log.Print("--------------")
	for _, i := range indexes {
		terms := r.TrendingTerms(i, indexes)
		for t, _ := range terms {
			log.Printf("%s: %30s %f", i, terms[t].Term, terms[t].Score)
		}
	}
	log.Print("--------------")
}
