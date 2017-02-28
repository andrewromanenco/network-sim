package nsim

type dummyMedium struct {
}

func (dm *dummyMedium) send(frame Frame) error {
	return nil
}
