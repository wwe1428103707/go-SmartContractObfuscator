package utils

import "github.com/tidwall/gjson"

type JsonQueue struct {
	element []gjson.Result
}

func NewJsonQueue() *JsonQueue {
	return &JsonQueue{
		element: []gjson.Result{},
	}
}

func (q *JsonQueue) Empty() bool {
	if len(q.element) == 0 {
		return true
	} else {
		return false
	}
}

func (q *JsonQueue) Size() int {
	return len(q.element)
}

func (q *JsonQueue) In(v gjson.Result) {
	q.element = append(q.element, v)
}

func (q *JsonQueue) Out() gjson.Result {
	if q.Empty() {
		return gjson.Result{}
	}
	temp := q.element[0]
	q.element = q.element[1:]
	return temp
}
