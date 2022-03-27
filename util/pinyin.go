package util

import (
	"github.com/mozillazg/go-pinyin"
	"strings"
)

type PinyinArgs *pinyin.Args

var PinyinFirstLetterArgs PinyinArgs = &pinyin.Args{
	Style: pinyin.FirstLetter,
	Fallback: func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	},
}

var PinyinNormalArgs PinyinArgs = &pinyin.Args{
	Style: pinyin.Normal,
	Fallback: func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	},
}

func Pinyin(s string, a PinyinArgs) string {
	if len(s) == 0 {
		return ""
	}
	var b strings.Builder
	arr := pinyin.Pinyin(s, *a)
	for _, v := range arr {
		for _, vv := range v {
			b.WriteString(vv)
		}
	}
	return b.String()
}
