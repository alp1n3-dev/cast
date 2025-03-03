package options

import (
	"maps"
	"strings"

	"github.com/alp1n3-eth/cast/internal/env"
)

func ParseReplacementValues(replacementSlice []string) *map[string]string {
	replacementPair := make(map[string]string)

	for _, h := range replacementSlice {
		if strings.Contains(h, ".env") {
			kvFileMap, _ := env.ReadKVFile(h)

			maps.Copy(replacementPair, *kvFileMap)
		} else {
			targetWord, value, _ := strings.Cut(h, "=")

			if len(targetWord) >= 1 {
				replacementPair[targetWord] = value
			}
		}
	}
	return &replacementPair
}
