package common

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"unsafe"
)

// FromBytes converts the specified byte array to a string.
func (*DuduStr) FromBytes(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

// Bytes converts the specified str to a byte array.
func (*DuduStr) ToBytes(str string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Contains determines whether the str is in the strs.
func (*DuduStr) Contains(str string, strs []string) bool {
	for _, v := range strs {
		if v == str {
			return true
		}
	}

	return false
}

// LCS gets the longest common substring of s1 and s2.
//
// Refers to http://en.wikibooks.org/wiki/Algorithm_Implementation/Strings/Longest_common_substring.
func (*DuduStr) LCS(s1 string, s2 string) string {
	var m = make([][]int, 1+len(s1))

	for i := 0; i < len(m); i++ {
		m[i] = make([]int, 1+len(s2))
	}

	longest := 0
	xLongest := 0

	for x := 1; x < 1+len(s1); x++ {
		for y := 1; y < 1+len(s2); y++ {
			if s1[x-1] == s2[y-1] {
				m[x][y] = m[x-1][y-1] + 1
				if m[x][y] > longest {
					longest = m[x][y]
					xLongest = x
				}
			} else {
				m[x][y] = 0
			}
		}
	}

	return s1[xLongest-longest : xLongest]
}

// 字符串生成md5
func (*DuduStr) Md5Str(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

//判断是否为手机号码
func (*DuduStr) IsMobile(phone string) int32 {
	mobile := "^1[\\d]{10}$"
	reg := regexp.MustCompile(mobile)
	isMobile := reg.MatchString(phone)
	if isMobile {
		return 200
	} else {
		return 40002
	}
}


//
////仿php的in_array
//func InArray(search interface{}, array []interface{}) bool {
//	for _,v := range array{
//		if search == v{
//			return true
//		}
//	}
//	return false
//}
//
////string转interface，结合inArray使用
//func StringToInterface(str []string) []interface{} {
//	newInterface := make([]interface{}, len(str))
//	for i, v := range str {
//		newInterface[i] = v
//	}
//
//	return newInterface
//}

