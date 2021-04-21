package pid

type Publisher interface {
	Publish(string, PID) error
}
