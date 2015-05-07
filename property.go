package reactor

import (
	"fmt"

	"github.com/sparkymat/reactor/property"
)

type Property struct {
	propertyType property.Type
	value        string
}

func (p Property) String() string {
	switch p.propertyType {
	case property.String:
		return fmt.Sprintf("\"%v\"", p.value)
	case property.Integer:
		return p.value
	case property.Float:
		return p.value
	case property.Object:
		return fmt.Sprintf("JSON.parse(\"%v\")", p.value)
	default:
		panic(fmt.Sprintf("Unknown object type %v", p.propertyType))
	}
}
