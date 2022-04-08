package worker

type Work interface {
	Do()
	Stop()
}
