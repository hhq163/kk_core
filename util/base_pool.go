package util

type BasePool interface {
	Run(i interface{})
	Shutdown()
}
