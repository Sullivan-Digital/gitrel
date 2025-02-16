package interfaces

type GitRelContext interface {
	Command() CommandContext
	Git() GitContext
	Output() OutputContext
}
