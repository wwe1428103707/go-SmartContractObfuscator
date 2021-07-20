package utils

import (
	"github.com/tidwall/gjson"
	"os"
)

func GetJsonContent(jsonpath string) gjson.Result {
	f, err := os.Open(jsonpath)
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

	jsoncontext := gjson.Parse(string(chunks))

	return jsoncontext
}
