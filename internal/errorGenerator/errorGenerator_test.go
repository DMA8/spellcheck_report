package errorGenerator

import (
	"log"
	"testing"
)

func OneRandomDoublingTest(t *testing.T) {
	inp := "абобоа"
	ans := (OneRandomDoubling(inp))
	if ans != inp {
		log.Println(ans)
		t.Error("aboba")
	}
}
