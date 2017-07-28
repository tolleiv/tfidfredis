package tfidfredis

import (
	"github.com/mediocregopher/radix.v2/redis"
	"log"
	"strconv"
	"math"
	"sort"
)

type Term struct {
	Term string
	Score float64
}

type Terms []Term

func (a Terms) Len() int           { return len(a) }
func (a Terms) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Terms) Less(i, j int) bool { return a[i].Score < a[j].Score }

type Repo struct {
	client *redis.Client
}

func NewRepo(addr string) (*Repo, error) {
	client, err := redis.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Repo{client: client}, nil
}

func (r *Repo) Add(term, index string) error {
	return r.client.Cmd("HINCRBY", index, term, "1").Err
}
func (r *Repo) Del(index string) error {
	return r.client.Cmd("DEL", index).Err
}

func (r *Repo) Tfidf(term, index string, indexes []string) float64 {
	return r.Tf(term, index) * r.Idf(term, indexes)
}

func (r *Repo) Tf(term, index string) float64 {
	termValue, _ := r.client.Cmd("HGET", index, term).Str()
	termCount, err := strconv.Atoi(termValue)
	if err != nil {
		//log.Fatalf("Error retrieving data (%s): %v",termValue, err)
		return 0
	}
	elems, err := r.client.Cmd("HGETALL", index).Map()
	if err != nil {
		log.Fatalf("Error retrieving data: %v", err)
	}
	sum := 0
	for i := range elems {
		value, err := strconv.Atoi(elems[i])
		if err != nil {
			log.Printf("Error retrieving increment: %v", err)
		}
		sum += value
	}
	//log.Printf("%d / %d \n", termCount, sum)
	return float64(termCount) / float64(sum)
}

func (r *Repo) Idf(term string, indexes []string) float64 {
	length := float64(len(indexes))
	cont := 1.0 + r.Containing(term, indexes)

	//log.Printf("Idf: %f / %f",len,cont )
	return math.Log(length / cont)
}

func (r *Repo) Containing(term string, indexes []string) float64 {
	count := 0
	for _, i := range indexes {
		termValue, _ := r.client.Cmd("HGET", i, term).Str()
		_, err := strconv.Atoi(termValue)
		if err == nil {
			count += 1.0
		}
	}
	return float64(count)
}

func (r *Repo) TrendingTerms(index string, indexes []string) Terms {
	elems, err := r.client.Cmd("HKEYS", index).List()
	if err != nil {
		log.Fatalf("Error retrieving data: %v", err)
	}
	result := Terms{}
	for i := range elems {
		result = append(result, Term{elems[i], r.Tfidf(elems[i], index, indexes)})
	}
	sort.Sort(sort.Reverse(result))
	return result
}