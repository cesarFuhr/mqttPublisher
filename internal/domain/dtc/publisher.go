package dtc

type Publisher interface {
	Publish(string, DTC) error
}
