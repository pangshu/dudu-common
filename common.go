package common

// 在这里感谢  b3log 的 https://github.com/b3log/gulu
// 另外也很感谢beego的各位开发人员，提取了beego的模板函数
// 本人在此基础上面，进行了一些修改，主要是增加了一些额外的常用方法。

type (
	DuduAes byte
	DuduFile byte
	DuduGo byte
	DuduLog byte
	DuduTpl byte

	DuduNet byte
	DuduOS byte
	DuduPanic byte
	DuduRand byte
	DuduRet byte
	DuduRune byte
	DuduStr byte
	DuduZip byte
)

var (
	Aes DuduAes
	File DuduFile
	Go DuduGo
	Log DuduLog
	Net DuduNet
	OS DuduOS
	Panic DuduPanic
	Rand DuduRand
	Ret DuduRet
	Rune DuduRune
	Str DuduStr
	Tpl DuduTpl
	Zip DuduZip
)
