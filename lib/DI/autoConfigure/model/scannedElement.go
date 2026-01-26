package autoConfigModel

type ScannedElement struct {
	name        string
	doc         string
	packageName string
	path        string
	fileName    string
	elementType string
}

func newElement(name string, elementType string) *ScannedElement {
	return &ScannedElement{
		name:        name,
		elementType: elementType,
	}
}

func NewStructElement(name string) *ScannedElement {
	return newElement(name, "struct")
}

func NewFuncElement(name string) *ScannedElement {
	return newElement(name, "func")
}

func (se *ScannedElement) GetName() string {
	return se.name
}

func (se *ScannedElement) GetType() string {
	return se.elementType
}

func (se *ScannedElement) GetPackageName() string {
	return se.packageName
}

func (se *ScannedElement) GetImportLine(moduleName string) string {
	return "import " + se.packageName + " \"" + moduleName + "/" + se.path + "\""
}

func (se *ScannedElement) Doc(doc string) *ScannedElement {
	se.doc = doc
	return se
}
func (se *ScannedElement) PackageName(packageName string) *ScannedElement {
	se.packageName = packageName
	return se
}

func (se *ScannedElement) FileName(fileName string) *ScannedElement {
	se.fileName = fileName
	return se
}

func (se *ScannedElement) Path(path string) *ScannedElement {
	se.path = path
	return se
}

func (se *ScannedElement) ToString() string {
	return se.name + "\n" +
		se.doc + "\n" +
		se.packageName + "\n" +
		se.path + "\n"
}
