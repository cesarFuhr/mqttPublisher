package status

type Publisher interface {
	Publish(string, Status) error
}
