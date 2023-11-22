package main

func main() {
	c := NewContainer()
	c.RegisterMany(&ServerService{}, &CalculateSingleton{})

	c.Run()
}

type CalculateSingleton struct {
}

func (c *CalculateSingleton) Add(a int, b int) int {
	return a + b
}

type ServerService struct {
	Calculate *CalculateSingleton `di:"direct"`
}

func (s *ServerService) Start(a ...any) error {
	println("Start")

	println(s.Calculate.Add(1, 2))
	return nil
}

func (s *ServerService) Stop(a ...any) error {
	println("Stop")
	return nil
}
