package main

import (
	"strings"

	"github.com/ktnyt/go-moji"
)

func containClass(tSlice []string, tClass string) bool {
	for _, class := range tSlice {
		if tClass == class {
			return true
		}
	}
	return false
}

func targetClass(kadaiName string) []string {
	kadaiName = moji.Convert(kadaiName, moji.ZE, moji.HE)

	var isJ2A bool = strings.HasPrefix(kadaiName, "J2A") || strings.HasPrefix(kadaiName, "JA") || strings.HasPrefix(kadaiName, "A")
	var isJ2B bool = strings.HasPrefix(kadaiName, "J2B") || strings.HasPrefix(kadaiName, "JB") || strings.HasPrefix(kadaiName, "B")
	var isS2 bool = strings.HasPrefix(kadaiName, "S2") || strings.HasPrefix(kadaiName, "S")

	if isJ2A && isJ2B && isS2 {
		return []string{"J2A", "J2B", "S2"}
	} else if isJ2A && isJ2B {
		return []string{"J2A", "J2B"}
	} else if isJ2A {
		return []string{"J2A"}
	} else if isJ2B {
		return []string{"J2B"}
	} else if isS2 {
		return []string{"S2"}
	} else {
		return []string{"J2A", "J2B", "S2"}
	}
}
