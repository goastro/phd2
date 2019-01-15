package phd2

// https://github.com/OpenPHDGuiding/phd2/wiki/EventMonitoring#event-notification-messages

func getEvent(name string) (interface{}, bool) { // nolint: gocyclo
	switch name {
	case "Version":
		return &VersionEvent{}, true
	case "CalibrationComplete":
		return &CalibrationCompleteEvent{}, true
	case "Paused":
		return &PausedEvent{}, true
	case "AppState":
		return &AppStateEvent{}, true
	case "LockPositionSet":
		return &LockPositionSetEvent{}, true
	case "Calibrating":
		return &CalibratingEvent{}, true
	case "StarSelected":
		return &StarSelectedEvent{}, true
	case "StartGuiding":
		return &StartGuidingEvent{}, true
	case "StartCalibration":
		return &StartCalibrationEvent{}, true
	case "CalibrationFailed":
		return &CalibrationFailedEvent{}, true
	case "CalibrationDataFlipped":
		return &CalibrationDataFlippedEvent{}, true
	case "LoopingExposures":
		return &LoopingExposuresEvent{}, true
	case "LoopingExposuresStopped":
		return &LoopingExposuresStoppedEvent{}, true
	case "SettleBegin":
		return &SettleBeginEvent{}, true
	case "Settling":
		return &SettlingEvent{}, true
	case "SettleDone":
		return &SettleDoneEvent{}, true
	case "StarLost":
		return &StarLostEvent{}, true
	case "GuidingStopped":
		return &GuidingStoppedEvent{}, true
	case "GuideStep":
		return &GuideStepEvent{}, true
	case "GuidingDithered":
		return &GuidingDitheredEvent{}, true
	case "LockPositionLost":
		return &LockPositionLostEvent{}, true
	case "Alert":
		return &AlertEvent{}, true
	case "GuideParamChange":
		return &GuideParamChangeEvent{}, true
	}

	return nil, false
}

// Event contains the common attributes of all events sent by PHD2.
type Event struct {
	// Event is the name of the event.
	Event string `json:"Event"`
	// Timestamp is the timesamp of the event in seconds from the epoch,
	// including fractional seconds.
	Timestamp float64 `json:"Timestamp"`
	// Host is the hostname of the machine running PHD2.
	Host string `json:"Host"`
	// Inst is the PHD2 instance number (1-based).
	Inst int `json:"Inst"`
}

// VersionEvent describes the PHD and message protocol versions.
type VersionEvent struct {
	Event
	PHDVersion string `json:"PHDVersion"`
	PHDSubver  string `json:"PHDSubver"`
	MsgVersion int    `json:"MsgVersion"`
}

// CalibrationCompleteEvent is sent when calibration is completed successfuly.
type CalibrationCompleteEvent struct {
	Event
	Mount string `json:"Mount"`
}

// PausedEvent is sent when guiding has been paused.
type PausedEvent struct {
	Event
}

// AppStateEvent is sent in the initial connection.
type AppStateEvent struct {
	Event
	State string `json:"State"`
}

// LockPositionSetEvent is sent when the lock position has been established.
type LockPositionSetEvent struct {
	Event
	X float64 `json:"X"`
	Y float64 `json:"Y"`
}

// CalibratingEvent is sent on each calibration step.
type CalibratingEvent struct {
	Event
	Mount string `json:"Mount"`
	Dir   string `json:"dir"`
	Dist  int    `json:"dist"`
	DX    int    `json:"dx"`
	DY    int    `json:"dy"`
	Pos   []int  `json:"pos"`
	Step  int    `json:"step"`
	State string `json:"State"`
}

// StarSelectedEvent is sent when a star is selected.
type StarSelectedEvent struct {
	Event
	X int `json:"X"`
	Y int `json:"Y"`
}

// StartGuidingEvent is sent when guiding begins.
type StartGuidingEvent struct {
	Event
}

// StartCalibrationEvent is sent when calibration begins.
type StartCalibrationEvent struct {
	Event
	Mount string `json:"Mount"`
}

// CalibrationFailedEvent is sent when calibration fails.
type CalibrationFailedEvent struct {
	Event
	Reason string `json:"Reason"`
}

// CalibrationDataFlippedEvent is sent when calibration data is flipped.
type CalibrationDataFlippedEvent struct {
	Event
	Mount string `json:"Mount"`
}

// LoopingExposuresEvent is sent for each exposure frame while looping exposures.
type LoopingExposuresEvent struct {
	Event
	Frame int `json:"Frame"`
}

// LoopingExposuresStoppedEvent is sent when looping exposures has stopped.
type LoopingExposuresStoppedEvent struct {
	Event
}

// SettleBeginEvent is sent when settling begins after a dither or guide operation.
type SettleBeginEvent struct {
	Event
}

// SettlingEvent is sent for each exposure frame after a dither or guide operation
// until guiding has settled.
type SettlingEvent struct {
	Event
	Distance   int     `json:"Distance"`
	Time       float64 `json:"Time"`
	SettleTime int     `json:"SettleTime"`
	StarLocked bool    `json:"StarLocked"`
}

// SettleDoneEvent is sent after a dither or guide operation indicating whether
// settling was achieved, or if the guider failed to settle before the time
// limit was reached, or if some other error occurred preventing guide or
// dither to complete and settle.
type SettleDoneEvent struct {
	Event
	Status        int    `json:"Status"`
	Error         string `json:"Error"`
	TotalFrames   int    `json:"TotalFrames"`
	DroppedFrames int    `json:"DroppedFrames"`
}

// StarLostEvent is sent when a frame has been dropped due to the star being lost.
type StarLostEvent struct {
	Event
	Frame     int     `json:"Frame"`
	Time      float64 `json:"Time"`
	StarMass  float64 `json:"StarMass"`
	SNR       float64 `json:"SNR"`
	AvgDist   float64 `json:"AvgDist"`
	ErrorCode int     `json:"ErrorCode"`
	Status    string  `json:"Status"`
}

// GuidingStoppedEvent is sent when guiding has stopped.
type GuidingStoppedEvent struct {
	Event
}

// ResumedEvent is sent when PHD2 has been resumed after having been paused.
type ResumedEvent struct {
	Event
}

// GuideStepEvent corresponds to a line in the PHD Guide Log. The event is sent for
// each frame while guiding.
type GuideStepEvent struct {
	Event
	// Frame is the frame number; starts at 1 each time guiding starts.
	Frame int `json:"Frame"`
	// Time is the time in seconds, including fractional seconds, since guiding
	// started.
	Time             float64 `json:"Time"`
	Mount            string  `json:"Mount"`
	DX               float64 `json:"dx"`
	DY               float64 `json:"dy"`
	RADistanceRaw    float64 `json:"RADistanceRaw"`
	DecDistanceRaw   float64 `json:"DecDistanceRaw"`
	RADistanceGuide  float64 `json:"RADistanceGuide"`
	DecDistanceGuide float64 `json:"DecDistanceGuide"`
	RADuration       int     `json:"RADuration"`
	RADirection      string  `json:"RADirection"`
	DecDuration      int     `json:"DECDuration"`
	DecDirection     string  `json:"DECDirection"`
	StarMass         float64 `json:"StarMass"`
	SNR              float64 `json:"SNR"`
	AvgDist          float64 `json:"AvgDist"`
	RALimited        bool    `json:"RALimited,omitempty"`
	DecLimited       bool    `json:"DecLimited,omitempty"`
	ErrorCode        int     `json:"ErrorCode"`
}

// GuidingDitheredEvent is sent when the lock position has been dithered.
type GuidingDitheredEvent struct {
	Event
	DX int `json:"dx"`
	DY int `json:"dy"`
}

// LockPositionLostEvent is sent when the lock position has been lost.
type LockPositionLostEvent struct {
	Event
}

// AlertEvent is sent when an alert message was displayed in PHD2.
type AlertEvent struct {
	Event
	Msg  string `json:"Msg"`
	Type string `json:"Type"`
}

// GuideParamChangeEvent is sent when a guiding parameter has been changed.
type GuideParamChangeEvent struct {
	Event
	Name  string `json:"Name"`
	Value string `json:"Value"`
}
