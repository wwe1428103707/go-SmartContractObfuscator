package local2state

import (
	"fmt"
	"github.com/huandu/go-clone"

	//"github.com/emirpasic/gods/lists"
	//"github.com/emirpasic/gods/lists/arraylist"
	//"github.com/emirpasic/gods/stacks"
	"github.com/tidwall/gjson"
	"math/rand"
	"notouchpure"
	"strconv"
	"strings"
	"utils"
	//"github.com/emirpasic/gods/stacks/arraystack"
)

type L2S struct {
	content     string
	jsoncontent gjson.Result
	NTP         notouchpure.NTP
	corpusDict  gjson.Result
}

func NewL2S(content string, jsoncontent gjson.Result) *L2S {
	return &L2S{
		content:     content,
		jsoncontent: jsoncontent,
		NTP:         *notouchpure.NewNTP(jsoncontent),
		corpusDict:  getCorpus(),
		//corpusDict:  gjson.Result{},
	}
}

func getCorpus() gjson.Result {
	return utils.GetJsonContent("Corpus.txt")
}

func (l2s *L2S) PreProcess() gjson.Result {

	random := rand.Float64() + 1
	var lvl = l2s.FindLocalVar(random)
	lvl = l2s.NTP.RunLocalVar(lvl)

	_ = l2s.ProcessSameName(lvl)

	//return localvarlist
	//localvarlist = l2s.NTP.RunLocalVar(localvarlist)
	//var samenamelist = l2s.ProcessSameName(localvarlist)
	//nowcontent := l2s.content
	//nowcontentjson := l2s.strReplace(gjson.Parse(nowcontent), samenamelist)
	//return nowcontentjson

	return gjson.Result{}
}

func (l2s *L2S) FindLocalVar(random float64) []gjson.Result {
	varList := l2s.FindASTNode("name", "VariableDeclaration", random)
	var res []gjson.Result
	for _, v := range varList {
		if v.Get("attributes").Get("stateVariable").Type.String() == "False" {
			res = append(res, v)
		}
	}
	return res
}

func (l2s *L2S) ProcessSameName(localvarlist []gjson.Result) []gjson.Result {

	var rl []string
	var temp [][]string
	//var temp = make([]interface{},0)
	nd := make(map[string]int)

	for _, v := range localvarlist {
		if gjson.Get(v.String(), "attributes").Exists() && gjson.Get(v.Get("attributes").String(), "name").Exists() {
			name := v.Get("attributes").Get("name").String()
			if name == "" {
				continue
			} else {
				pos := l2s.srcToPos(v.Get("src"))
				pos[0] = pos[1] - len(name)
				id := v.Get("id").String()
				//var spos = make([]interface{},0)
				spos := []string{}
				spos = append(spos, name, string(pos[0]), string(pos[1]), id)
				temp = append(temp, spos)
			}
		}
	}
	//fmt.Println(temp)

	for _, v := range temp {
		nd[v[0]] = 0
	}
	for _, v := range temp {
		nd[v[0]] += 1
	}

	//fmt.Println(nd)

	for k, v := range nd {
		//fmt.Println(k,v)
		if v > 1 {
			rl = append(rl, l2s.Rename(k, v, temp)...)
		}
	}

	templist := clone.Clone(rl)

	for i, v := range templist {

	}
	//fmt.Println(rl)

	return []gjson.Result{}
}

func (l2s *L2S) strReplace(parse gjson.Result, samenamelist gjson.Result) gjson.Result {
	return gjson.Result{}
}

func (l2s *L2S) FindASTNode(keyin string, valuein string, random float64) []gjson.Result {
	stack := utils.NewJsonQueue()
	stack.In(l2s.jsoncontent)
	var res []gjson.Result
	for !stack.Empty() {
		tempJson := stack.Out()
		tempJson.ForEach(func(key, value gjson.Result) bool {
			if key.String() == keyin && value.String() == valuein && rand.Float64() < random {
				res = append(res, tempJson)
			} else if value.Type.String() == "JSON" {
				for _, temp := range value.Array() {
					if temp.Type.String() == "JSON" {
						stack.In(temp)
					}
				}
			}
			return true
		})
	}

	return res
}

func (l2s *L2S) srcToPos(src gjson.Result) []int {
	temp := src.String()
	temps := strings.Split(temp, ":")
	//fmt.Println(temps)

	num1, _ := strconv.Atoi(temps[0])
	num2, _ := strconv.Atoi(temps[1])
	res := []int{num1, num1 + num2}
	return res
}

func (l2s *L2S) FindSameNameState(k string, v int) []int {

	return []int{}
}

func (l2s *L2S) Rename(k string, v int, ls [][]string) []string {
	rl := []string{}
	var newname string
	for _, v := range ls {
		if v[0] == k {
			nl := l2s.corpusDict.Get("variableNaming").Array()
			rd := rand.Intn(len(nl)) - 1
			for i, name := range nl {
				if i != rd {
					continue
				} else {
					newname = l2s.ShuffleStr(name.String())
					rl = append(rl, newname)
					rl = append(rl, v...)
				}
			}
		}
	}
	fmt.Println(rl)
	return rl
}

func (l2s *L2S) ShuffleStr(name string) string {

	inRune := []rune(name)

	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})

	//fmt.Println(string(inRune))
	return string(inRune)
}
