package common

import (
	"math/rand"
	"time"
)

// Ints returns a random integer array with the specified from, to and size.
func (*DuduRand) Ints(from, to, size int) []int {
	if to-from < size {
		size = to - from
	}

	var slice []int
	for i := from; i < to; i++ {
		slice = append(slice, i)
	}

	var ret []int
	for i := 0; i < size; i++ {
		idx := rand.Intn(len(slice))
		if slice != nil {
			ret = append(ret, slice[idx])
			slice = append(slice[:idx], slice[idx+1:]...)
		}

	}

	return ret
}
//随机区间值
func (*DuduRand) Int64(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Nanosecond)
	return min + rand.Int63n(max-min+1)
}

func (*DuduRand) Int(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Nanosecond)
	return min + rand.Intn(max-min+1)
}

//生成count个[start,end)结束的不重复的随机数
func (*DuduRand) IntNoRepeat(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}
	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Nanosecond)
	for len(nums) < count {
		//生成随机数
		//num := r.Intn((end - start)) + start
		num := start + rand.Intn(end-start+1)
		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}

// String returns a random string ['a', 'z'] in the specified length
func (*DuduRand) String(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	time.Sleep(time.Nanosecond)

	letter := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}

	return string(b)
}

//随机几位字符串
func (*DuduRand) StringByNum(lenString int) string{
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Nanosecond)
	var i int
	for i = 0; i < lenString; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

//随机几位字符串
func (*DuduRand) StringByStr(lenString int) string{
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Nanosecond)
	var i int
	for i = 0; i < lenString; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

////随机几位字符串
//func (*DuduRand) GetRandomString(lenString int) string{
//	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
//	bytes := []byte(str)
//	result := []byte{}
//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
//	var i int
//	for i = 0; i < lenString; i++ {
//		result = append(result, bytes[r.Intn(len(bytes))])
//	}
//	return string(result)
//}
