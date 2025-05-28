package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	const s = "สวัสดี"

	fmt.Println("Len", len(s))

	for i := 0; i < len(s); i++ {
		fmt.Printf("%x ", s[i])
	}
	fmt.Println()

	runes := []rune(s) // string convert to rune切片，每个 rune 代表一个 unicode 字符
	fmt.Println("Runes", runes)
	fmt.Println("Rune count", len(runes))
	fmt.Println("Rune count", utf8.RuneCountInString(s)) // 用于计算字符串 s 中 Unicode 字符（rune）的数量，而不是字节数

	//rune 切片转换回字符串
	s2 := string(runes)
	fmt.Println(s2)

	for idx, runeValue := range s {
		fmt.Printf("%#U starts at byte position %d\n", runeValue, idx)
	}

	fmt.Println("\nUsing DecodeRuneInString")
	for i, w := 0, 0; i < len(s); i += w {
		runeValue, width := utf8.DecodeRuneInString(s[i:])
		fmt.Printf("%#U starts at byte position %d\n", runeValue, i)
		fmt.Println("width", width)
		w = width

		examineRune(runeValue)
	}
}

func examineRune(r rune) {
	if r == 't' {
		fmt.Println("found tee")
	} else if r == 'ส' {
		fmt.Println("found so sua")
	}
}
