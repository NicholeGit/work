// test
package main

import (
	"fmt"
	"strings"
)

func splitInit(s rune) bool {
	if s == ' ' || s == '"' {
		return true
	}
	return false
}

func splitUnderline(s rune) bool {
	if s == '_' {
		return true
	}
	return false
}

func main() {
	str := "DEPOSIT_ITEMS ({\"100313130_5_8\",\"100102496_8_8\",\"100101132_8_8\",\"0_0_0\",})"

	ret := strings.FieldsFunc(str, splitInit)

	for _, v := range ret {
		if strings.Count(v, "_") == 2 {
			under := strings.FieldsFunc(v, splitUnderline)
			fmt.Println(under)
			fmt.Println(v)
		}
	}

}

//func main() {
//	file, err := os.Open("d:\\log")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer file.Close()

//	r := bufio.NewReader(file)

//	for i := 0; i < 2; i++ {
//		_, err := r.ReadBytes('\n')
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//	}

//	var delim = []byte(`","`)
//	for {
//		line, err := r.ReadBytes('\n')
//		if err != nil {
//			if err != io.EOF {
//				fmt.Println(err)
//				return
//			}
//			if len(line) == 0 {
//				return
//			}
//		}

//		index := bytes.IndexByte(line, ' ')
//		if index == -1 {
//			fmt.Println("沒有找到前綴")
//			return
//		}
//		prefix := line[:index]
//		fmt.Println("前綴", string(prefix))

//		contentIndex1 := bytes.IndexByte(line, '{')
//		if contentIndex1 == -1 {
//			fmt.Println("沒有找到'{'")
//			return
//		}

//		contentIndex2 := bytes.IndexByte(line, '}')
//		if contentIndex2 == -1 {
//			fmt.Println("沒有找到'}'")
//			return
//		}
//		if contentIndex1+1 == contentIndex2 {
//			continue
//		}

//		bytesArr := bytes.Split(line[contentIndex1+2:contentIndex2-2], delim)
//		for _, bs := range bytesArr {
//			fmt.Println(string(bs))
//		}
//	}
//}
