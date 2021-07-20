package utils

import "github.com/tidwall/gjson"

type JsonStack struct {
	element []gjson.Result
	size    int
}

func (js *JsonStack) Push(v gjson.Result) {
	js.element = append(js.element, v)
	js.size++
}

func (js *JsonStack) Pop() gjson.Result {
	if js.Empty() {
		return gjson.Result{}
	}
	js.size--
	temp := js.element[js.size]
	js.element = js.element[:js.size]
	return temp
}

func (js *JsonStack) Empty() bool {
	return js.size == 0
}

func (js *JsonStack) Peek() gjson.Result {
	if js.Empty() {
		return gjson.Result{}
	}
	return js.element[js.size-1]
}

func (js *JsonStack) Size() int {
	return js.size
}

func NewJsonStack() *JsonStack {
	return &JsonStack{
		element: []gjson.Result{},
		size:    0,
	}
}
