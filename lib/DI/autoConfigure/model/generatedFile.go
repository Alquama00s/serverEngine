package autoConfigModel

type GeneratedFile struct {
	path     string
	FileName string
	Contents string
}

func GetNewGeneratedFile(path string) *GeneratedFile {
	return &GeneratedFile{
		path: path,
	}
}

func (gf *GeneratedFile) GetPath() string {
	return gf.path
}
