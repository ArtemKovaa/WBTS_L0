package pkg

import (
	"fmt"
	"strings"
)

func GeneratePlaceholders(n int) string {
	ps := make([]string, n)
	for i := range n {
		ps[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(ps, ",")
}
