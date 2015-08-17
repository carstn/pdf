package pdf

import "reflect"

type Annotation struct {
	V Value
	X float64
	Y float64
}

// Annots returns all the annotations on the received page.
func (p Page) Annots() (annots map[int]string) {
	annots = make(map[int]string)
	annotsObj := p.V.Key("Annots")
	for i := 0; i < annotsObj.Len(); i++ {
		var annot string
		indivAnnot := annotsObj.Index(i)
		dictAnnot := indivAnnot.data.(dict)
		//rectAnnot := dictAnnot["Rect"].(Rect)
		//fmt.Printf("min: (%f, %f)  |  max: (%f, %f)\n", rectAnnot.Min.X, rectAnnot.Min.Y, rectAnnot.Max.X, rectAnnot.Max.Y)
		valAnnot := dictAnnot["V"]
		if reflect.ValueOf(valAnnot).Kind() != reflect.String {
			continue
		}
		annot = valAnnot.(string)

		if annot != "" {
			annots[i] = annot
		}
	}
	return
}
