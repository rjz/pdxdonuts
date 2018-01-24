package generate_test

import (
	"github.com/rjz/forager/generate"
	"os"
	"testing"
)

func TestDo(t *testing.T) {
	generate.RootDir = ".."
	dir := os.TempDir() + "/forager"
	err := generate.Do(dir, &generate.PageData{})
	if err != nil {
		t.Error(err)
	}
}
