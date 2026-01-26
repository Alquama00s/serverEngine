package autoconfigure

import (
	"os"
	"testing"

	autoconfigure "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure"
)

// func TestScanner(t *testing.T) {
// 	os.Chdir("../../test_server")
// 	res := autoconfigure.Scan("@service", ".")
// 	m := autoconfigure.GetAppContext().GetModuleName()
// 	for _, se := range res {
// 		t.Log(se.ToString())
// 		t.Log(se.GetImportLine(m))
// 	}
// }

func TestGen(t *testing.T) {
	os.Chdir("../../test_server")
	autoconfigure.WriteService(".")
}
