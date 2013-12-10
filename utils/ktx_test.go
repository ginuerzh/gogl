// ktx_test.go
package utils

import (
	gl "github.com/chsc/gogl/gl42"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestLoadKtx(t *testing.T) {
	gl.Init()
	LoadKtx("brick.ktx", 0)
}
