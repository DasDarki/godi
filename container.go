package main

import (
	"reflect"
)

type Container struct {
	types      map[string]interface{}
	instances  []*wrappedInstance
	isPrepared bool
}

func NewContainer() *Container {
	return &Container{
		types:      make(map[string]interface{}),
		instances:  make([]*wrappedInstance, 0),
		isPrepared: false,
	}
}

func (c *Container) findInstance(name string) interface{} {
	if instance, ok := c.types[name]; ok {
		return instance
	}

	return nil
}

func (c *Container) addWrapperInstance(instance interface{}, start serviceCallback, stop serviceCallback) {
	mapping, err := mappingFromStruct(instance)
	if err != nil {
		panic(err)
	}

	c.instances = append(c.instances, &wrappedInstance{
		mapping: mapping,
		start:   start,
		stop:    stop,
	})
}

func wrapServiceCallback(method reflect.Value) serviceCallback {
	return func(a ...any) error {
		if method.IsValid() {
			in := make([]reflect.Value, len(a))
			for i, arg := range a {
				in[i] = reflect.ValueOf(arg)
			}

			ret := method.Call(in)[0]
			if ret.IsValid() {
				retval := ret.Interface()
				if retval != nil {
					if err, ok := retval.(error); ok {
						return err
					}
				}
			}

			return nil
		}

		return nil
	}
}

func (c *Container) Register(instance interface{}) {
	if reflect.TypeOf(instance).Kind() != reflect.Ptr {
		panic(ErrNotAPointer)
	}

	typeOfStruct := reflect.TypeOf(instance).Elem()
	if typeOfStruct.Kind() != reflect.Struct {
		panic(ErrNotAStruct)
	}

	c.types[getNameForInstance(typeOfStruct)] = instance

	startMethod := reflect.ValueOf(instance).MethodByName("Start")
	stopMethod := reflect.ValueOf(instance).MethodByName("Stop")

	start := wrapServiceCallback(startMethod)
	stop := wrapServiceCallback(stopMethod)

	c.addWrapperInstance(instance, start, stop)
}

func (c *Container) RegisterMany(instances ...interface{}) {
	for _, instance := range instances {
		c.Register(instance)
	}
}

func (c *Container) Prepare() {
	if c.isPrepared {
		return
	}

	c.isPrepared = true

	for _, instance := range c.instances {
		instance.mapping.resolve(c, resolveStagePrepare)
	}
}

func (c *Container) Run() error {
	if !c.isPrepared {
		c.Prepare()
	}

	for _, instance := range c.instances {
		instance.mapping.resolve(c, reolsveStageRun)

		err := instance.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Container) Stop() error {
	for i := len(c.instances) - 1; i >= 0; i-- {
		instance := c.instances[i]

		err := instance.Stop()
		if err != nil {
			return err
		}
	}

	return nil
}
