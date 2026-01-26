package generate

import (
	"log"
	"os"
	"strings"
	"text/template"
)

func main() {
	CreateFileFromTemplate("previlege", "internal/constants")
	for _, a := range os.Args {
		switch a {
		case "repo":
			CreateFileFromTemplate("repo", "internal/repo")

		case "service":
			CreateFileFromTemplate("service", "internal/service")

		case "controller":
			CreateFileFromTemplate("controller", "internal/controller")
		}
	}

}

func CreateFileFromTemplate(templateName string, outFolder string) {
	if len(os.Args) < 2 {
		panic("struct name required")
	}
	structName := os.Args[1]
	p, _ := os.Getwd()
	var suffix = strings.Split(p, "internal")[1]
	var suffixLen = len(strings.Split(suffix, "/"))
	for suffixLen > 0 {
		os.Chdir("../")
		suffixLen--
	}
	// Read template from file
	tmpl, err := template.ParseFiles("goGen/templates/" + templateName + ".tmpl")
	if err != nil {
		log.Println(os.Getwd())
		panic(err)
	}

	// Create output file
	outFile := strings.ToLower(structName) + "_" + templateName + ".gen.go"
	err = os.Chdir(outFolder)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	k := map[string]string{
		"StructName": structName,
		"structName": strings.ToLower(structName),
	}

	// Execute template
	err = tmpl.Execute(f, k)
	if err != nil {
		panic(err)
	}

	log.Println("Generated", outFile)
}
