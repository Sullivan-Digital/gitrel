package interfaces

type OutputContext interface {
	Print(args ...interface{})
	Println(args ...interface{})
	Printf(format string, args ...interface{})
}
