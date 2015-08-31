package pdf

import (
	"fmt"
	"strings"
)

type Annotation struct {
	S   string
	Min Point
	Max Point
}

func (a *Annotation) String() string {
	return fmt.Sprintf("min:(%.1f, %.1f) max:(%.1f, %.1f) %s", a.Min.X, a.Min.Y, a.Max.X, a.Max.Y, a.S)
}

// Annots returns all the annotations on the received page.
// FIXME: After testing, remove any if blocks of the form `if !as.Production {...}`.
func (p Page) Annots() (annots []Annotation) {
	annotsObj := p.V.Key("Annots")
AnnotLoop:
	for i := 0; i < annotsObj.Len(); i++ {
		// Select a single annotation.
		annotObj := annotsObj.Index(i)

		// Get dataObj as type dict.
		// dataObj is the top `dict` in the annotation object, containing all the annotation's fields (including other `dict`s).
		if annotObj.Kind() != Dict {
			continue AnnotLoop
		}

		// Get rectObj as type pdf.array.
		if annotObj.Key("Rect").Kind() != Array {
			continue AnnotLoop
		}
		rectObj := annotObj.Key("Rect")

		// Convert rectObj (a pdf.array) into a rect (a [4]float64).
		if rectObj.Len() != 4 {
			continue AnnotLoop
		}
		rect := make([]float64, 4)
		for i := 0; i < 4; i++ {
			fl := rectObj.Index(i).Float64()
			if fl == 0 {
				continue AnnotLoop
			}
			rect[i] = fl
		}

		// Get valueObj as type string.
		// valueObj is the actual value displayed on the PDF. In this case, we're looking for a string.
		if annotObj.Key("V").Kind() != String {
			continue AnnotLoop
		}
		valueObj := annotObj.Key("V").RawString()
		valueObj = strings.TrimSpace(valueObj)
		if len(valueObj) == 0 {
			continue AnnotLoop
		}

		a := &Annotation{
			S:   valueObj,
			Min: Point{float64(rect[0]), float64(rect[1])},
			Max: Point{float64(rect[2]), float64(rect[3])},
		}

		annots = append(annots, *a)
	}
	return
}
