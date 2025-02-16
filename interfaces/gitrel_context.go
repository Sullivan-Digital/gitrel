package interfaces

type GitRelContext interface {
	Options() CommandContext
	Git() GitContext
	Output() OutputContext
}