package autoconfigure

import (
	"os"
	"testing"

	"github.com/Alquama00s/serverEngine/lib/DI"
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
	DI.InitialiseContextBuilder(".")
}
