package scraper

import (
	"strings"
)

// Implements Parser interface
type KvpParser struct {
}

func (p *KvpParser) Run(
	data []byte,
) map[string]string {

	values := make(map[string]string)

	for _, line := range strings.Split(string(data), "\n") {
		entries := strings.SplitN(line, ":", 2)

		// Empty line
		if len(entries) == 1 {
			continue
		}

		key := strings.TrimSpace(entries[0])
		value := strings.TrimSpace(entries[1])

		values[key] = value
	}

	return values
}
