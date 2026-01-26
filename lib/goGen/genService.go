package generate

import (
	"os"
)

func Write(path string, content string) {
	os.WriteFile(path, []byte(content), 0644)
}
