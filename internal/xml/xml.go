package xml

import (
	"fmt"
	"github.com/beevik/etree"
)

const (
	errorCode3000 = iota + 3000
	errorCode3001 = iota + 3000
	errorCode3002 = iota + 3000
)

const (
	errorMessage         = "error code: %d - %s"
	errorElementNotFound = "element not found"
)

// SetVersion set version in a xml path
func SetVersion(filePath string, xPath string, value string) error {

	doc := etree.NewDocument()
	if err := doc.ReadFromFile(filePath); err != nil {
		return fmt.Errorf(errorMessage, errorCode3000, err)
	}

	// Suchen Sie das Element basierend auf einem Pfad.
	// Im Beispiel suchen wir nach dem ersten "version"-Element unter "project".
	element := doc.FindElement(xPath)
	if element == nil {
		// ToDo ErrorCode Erstellen
		return fmt.Errorf(errorMessage, errorCode3002, errorElementNotFound)
	}

	// Setzen Sie den Wert des gefundenen Elements.
	element.SetText(value)

	// Speichern Sie das aktualisierte Dokument zurück in eine Datei.
	doc.Indent(2)
	if err := doc.WriteToFile(filePath); err != nil {
		return fmt.Errorf(errorMessage, errorCode3001, err)
	}

	return nil
}
