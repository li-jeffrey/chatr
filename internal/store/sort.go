package store

import "sort"

type Ordering func(s1, s2 *Submission) bool

func (order Ordering) Reversed() Ordering {
	return func(s1, s2 *Submission) bool {
		return !order(s1, s2)
	}
}

func (order Ordering) by(subs []Submission) {
	s := &subSorter{
		subs, order,
	}
	sort.Sort(s)
}

type subSorter struct {
	subs  []Submission
	order Ordering
}

func (s *subSorter) Len() int {
	return len(s.subs)
}

func (s *subSorter) Less(i, j int) bool {
	return s.order(&s.subs[i], &s.subs[j])
}

func (s *subSorter) Swap(i, j int) {
	s.subs[i], s.subs[j] = s.subs[j], s.subs[i]
}
