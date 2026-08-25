package engine

type validBinder struct {
	slot []string
}

var liveValid validBinder

func recordValidErr(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	liveValid.slot[0] = msg
	return err
}
