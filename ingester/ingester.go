package ingester

import (
	"bufio"
	"os"
)

var (
	reader = bufio.NewScanner(os.Stdin)
)

func IngestFromStdin(ingress chan<- string) {
	flowStarted := false
	for {
		reader.Scan()
		json := reader.Text()
		if !flowStarted && json == "BEGIN" {
			flowStarted = true
			continue
		}
		if flowStarted {
			if json == "END" {
				close(ingress)
				break
			} else {
				ingress <- json
			}
		}
	}
}
