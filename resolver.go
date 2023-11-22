package godi

import "reflect"

type resolveStage int

const (
	resolveStageUnresolved resolveStage = 0
	resolveStagePrepare    resolveStage = 1
	reolsveStageRun        resolveStage = 2
)

type resolveMapping struct {
	instance interface{}
	fields   map[string]string
}

func getNameForInstance(typeOfInstance reflect.Type) string {
	return typeOfInstance.Name() + "." + typeOfInstance.PkgPath()
}

func mappingFromStruct(instance interface{}) (*resolveMapping, error) {
	mapping := &resolveMapping{
		instance: instance,
		fields:   make(map[string]string),
	}

	typeOfInstance := reflect.TypeOf(instance).Elem()
	if typeOfInstance.Kind() != reflect.Struct {
		return nil, ErrNotAStruct
	}

	for i := 0; i < typeOfInstance.NumField(); i++ {
		field := typeOfInstance.Field(i)
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			if tag, ok := field.Tag.Lookup("di"); ok {
				if tag != "transient" && tag != "direct" {
					return nil, ErrTagInvalid
				}

				mapping.fields[field.Name] = getNameForInstance(field.Type.Elem())
			}
		}
	}

	return mapping, nil
}

func (mapping *resolveMapping) resolve(c *Container, stage resolveStage) {
	for name, typeName := range mapping.fields {
		_, ok := reflect.TypeOf(mapping.instance).Elem().FieldByName(name)
		if !ok {
			panic("field not found! this should never happen")
		}

		if reflect.ValueOf(mapping.instance).Elem().FieldByName(name).IsNil() {
			instance := c.findInstance(typeName)
			if instance == nil {
				if stage == reolsveStageRun {
					panic(ErrInstanceNotFound)
				}
				continue
			}

			reflect.ValueOf(mapping.instance).Elem().FieldByName(name).Set(reflect.ValueOf(instance))
		}
	}
}
