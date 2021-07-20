package notouchpure

import (
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"utils"
)

const (
	PURE_FLAG = "pure"
	VIEW_FLAG = "view"
)

type NTP struct {
	//ntp string
	jsoncontent gjson.Result
}

func NewNTP(content gjson.Result) *NTP {
	return &NTP{jsoncontent: content}
}

func (n *NTP) RunLocalVar(localvarlist []gjson.Result) []gjson.Result {
	funcNode := n.FindASTNode("name", "FunctionDefinition")
	var pureFuncNode []gjson.Result
	for _, v := range funcNode {

		//fmt.Println(v.Get("attributes").Get("stateMutability"))
		//注意这里是attribute{s}不是attribute
		if v.Get("attributes").Get("stateMutability").String() == PURE_FLAG {
			pureFuncNode = append(pureFuncNode, v)
		} else if v.Get("attributes").Get("stateMutability").String() == VIEW_FLAG {
			pureFuncNode = append(pureFuncNode, v)
		} else {
			continue
		}
	}
	//fmt.Println(pureFuncNode)
	posl := n.GetStartEndPos(pureFuncNode)
	//fmt.Println(posl)
	var npd []gjson.Result
	for _, v := range localvarlist {
		l := n.SrcToPos(v.Get("src"))
		if n.PureBool(posl, l[0], l[1]) {
			continue
		} else {
			npd = append(npd, v)
		}
	}
	//fmt.Println(npd)
	return npd
}

func (n *NTP) FindASTNode(keyin string, valuein string) []gjson.Result {
	q := utils.NewJsonQueue()
	q.In(n.jsoncontent)
	var res []gjson.Result
	for !q.Empty() {
		tempJson := q.Out()
		tempJson.ForEach(func(key, value gjson.Result) bool {
			if key.String() == keyin && value.String() == valuein {
				//fmt.Println(value.String())
				res = append(res, tempJson)
			} else if value.Type.String() == "JSON" {
				for _, temp := range value.Array() {
					if temp.Type.String() == "JSON" {
						q.In(temp)
					}
				}
			}
			return true
		})
	}
	return res
}

func (n *NTP) GetStartEndPos(node []gjson.Result) [][]int {
	var pl [][]int
	for _, v := range node {
		//fmt.Println(v.Get("src"))
		t := n.SrcToPos(v.Get("src"))
		pl = append(pl, t)
	}
	return pl
}

func (n *NTP) SrcToPos(src gjson.Result) []int {
	temp := src.String()
	temps := strings.Split(temp, ":")
	//fmt.Println(temps)

	num1, _ := strconv.Atoi(temps[0])
	num2, _ := strconv.Atoi(temps[1])
	res := []int{num1, num1 + num2}
	return res
}

func (n *NTP) PureBool(posl [][]int, sp int, ep int) bool {
	for _, v := range posl {
		if sp > v[0] && ep < v[1] {
			return true
		} else {
			continue
		}
	}
	return false
}
