package dataflowobfuscation

import (
	"fmt"
	"github.com/tidwall/gjson"
	"local2state"
	"os"
	"strings"
	"utils"
)

type Dfo struct {
	filepath       string
	jsonpath       string
	outputfilename string
	solcontent     string
	json           gjson.Result
	middlecontract string
	middlejsonast  string
	configpath     string
	featurelist    gjson.Result
	l2s            local2state.L2S
}

func NewDfo(filepath string, jsonpath string, configpath string) *Dfo {
	return &Dfo{
		filepath:       filepath,
		jsonpath:       jsonpath,
		outputfilename: GetOutputFilePath(filepath),
		solcontent:     GetContent(filepath),
		json:           utils.GetJsonContent(jsonpath),
		middlecontract: "temp.sol",
		middlejsonast:  "temp.sol_json.ast",
		configpath:     configpath,
		featurelist:    GetConfig(configpath),
	}
}

func GetOutputFilePath(filepath string) string {
	var temp = strings.Split(filepath, ".sol")
	var newfilename = temp[0] + "_dataflow_confuse.sol"
	return newfilename
}

func GetContent(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	chunks := make([]byte, 1024, 1024)
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil {
			break
		}
		//fmt.Println(string(buf[:n]))
		chunks = append(chunks, buf[:n]...)
	}

	return string(chunks)
}

func GetConfig(configfile string) gjson.Result {
	config := utils.GetJsonContent(configfile)
	//fmt.Println(config.Get("activateFunc"))
	return config.Get("activateFunc")
}

func isActive(name string, dict gjson.Result) bool {
	result := dict.Get("#." + name)
	var temp bool
	for _, name := range result.Array() {
		var tempa = name.Array()
		temp = tempa[0].Bool()
	}
	return temp
}

func (dfo *Dfo) Run() int {
	if isActive("local2State", dfo.featurelist) {
		//fmt.Println("yes!")
		dfo.l2s = *local2state.NewL2S(dfo.solcontent, dfo.json)
		newcontent := dfo.l2s.PreProcess()
		fmt.Println(newcontent.String())
	} else {
		fmt.Println("No!")
	}
	return 1

}
