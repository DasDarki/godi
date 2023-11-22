package godi

type serviceCallback func(a ...any) error

type wrappedInstance struct {
	mapping *resolveMapping
	start   serviceCallback
	stop    serviceCallback
}

func (instance *wrappedInstance) Start(a ...any) error {
	if instance.start == nil {
		return nil
	}

	return instance.start(a...)
}

func (instance *wrappedInstance) Stop(a ...any) error {
	if instance.stop == nil {
		return nil
	}

	return instance.stop(a...)
}
