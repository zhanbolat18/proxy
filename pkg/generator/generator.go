package generator

type Generator interface {
	Id() (string, error)
}
