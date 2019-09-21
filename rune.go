package common

// IsNumOrLetter checks the specified rune is number or letter.
func (*DuduRune) IsNumOrLetter(r rune) bool {
	return ('0' <= r && '9' >= r) || Rune.IsLetter(r)
}

// IsLetter checks the specified rune is letter.
func (*DuduRune) IsLetter(r rune) bool {
	return 'a' <= r && 'z' >= r || 'A' <= r && 'Z' >= r
}