package json

import (
	"strconv"
	"math/rand"
	"time"
)

/**
Go 语言里面原生支持了这种数据格式的序列化以及反序列化，内部使用反射机制实现，性能有点差，在高度依赖 json 解析的应用里，往往会成为性能瓶颈
*/
type user struct {
	Name string `json:"nameabcd"`
	Age  int    `json:"age"`
	Desc string `json:"desc"`
}

func generate() {
	var userRows []user = []user{}

	var u user
	for i := 1; i < 100; i++ {
		u.Age = 10

		is := strconv.Itoa(i)
		u.Name = "name" + is

		u.Desc = GetText(300)

		userRows = append(userRows, u)

		u = user{}
	}

	_, _ = json.Marshal(userRows)
}


var _ = fmt.Sprintf("")

func GetText(slen int) string {
	character := "abcdefghjkmnpqrstuvwxyABCDEFGHJKLMNPQRSTUVWXYZ023456789,:.;-+";
	maxlen := int32(len(character))

	var result string
	rand.Seed(time.Now().UnixNano())
	for i:= 0; i < slen; i++ {
		idx := rand.Int31n(maxlen)
		result += (string)(character[idx])
	}
	return result
}