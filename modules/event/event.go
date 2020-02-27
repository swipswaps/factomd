package event

// Detect minute change excluding minute 10 boundary
func (evt *DBHT) MinuteChanged(newEvt *DBHT) bool {

	if newEvt.Minute == 10 { // ignore min 10 - we want to trigger on move from 9->0 instead
		return false
	}

	if newEvt.Minute == evt.Minute && newEvt.DBHeight == evt.DBHeight {
		return false // no change
	} else {
		return true
	}
}
