package phd2

// https://github.com/OpenPHDGuiding/phd2/wiki/EventMonitoring#event-notification-messages

func getEvent(name string) (interface{}, bool) {
	switch name {
	case "Version":
		return &Version{}, true
	case "CalibrationComplete":
		return &CalibrationComplete{}, true
	case "Paused":
		return &Paused{}, true
	case "AppState":
		return &AppState{}, true
	case "LockPositionSet":
		return &LockPositionSet{}, true
	case "Calibrating":
		return &Calibrating{}, true
	case "StarSelected":
		return &StarSelected{}, true
	case "StartGuiding":
		return &StartGuiding{}, true
	case "StartCalibration":
		return &StartCalibration{}, true
	case "CalibrationFailed":
		return &CalibrationFailed{}, true
	case "CalibrationDataFlipped":
		return &CalibrationDataFlipped{}, true
	case "LoopingExposures":
		return &LoopingExposures{}, true
	case "LoopingExposuresStopped":
		return &LoopingExposuresStopped{}, true
	case "SettleBegin":
		return &SettleBegin{}, true
	case "Settling":
		return &Settling{}, true
	case "SettleDone":
		return &SettleDone{}, true
	case "StarLost":
		return &StarLost{}, true
	case "GuidingStopped":
		return &GuidingStopped{}, true
	case "GuideStep":
		return &GuideStep{}, true
	case "GuidingDithered":
		return &GuidingDithered{}, true
	case "LockPositionLost":
		return &LockPositionLost{}, true
	case "Alert":
		return &Alert{}, true
	case "GuideParamChange":
		return &GuideParamChange{}, true
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

// Version describes the PHD and message protocol versions.
type Version struct {
	Event
	PHDVersion string `json:"PHDVersion"`
	PHDSubver  string `json:"PHDSubver"`
	MsgVersion int    `json:"MsgVersion"`
}

// CalibrationComplete is sent when calibration is completed successfuly.
type CalibrationComplete struct {
	Event
	Mount string `json:"Mount"`
}

// Paused is sent when guiding has been paused.
type Paused struct {
	Event
}

// AppState is sent in the initial connection.
type AppState struct {
	Event
	State string `json:"State"`
}

// LockPositionSet is sent when the lock position has been established.
type LockPositionSet struct {
	Event
	X int `json:"X"`
	Y int `json:"Y"`
}

// Calibrating is sent on each calibration step.
type Calibrating struct {
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

// StarSelected is sent when a star is selected.
type StarSelected struct {
	Event
	X int `json:"X"`
	Y int `json:"Y"`
}

// StartGuiding is sent when guiding begins.
type StartGuiding struct {
	Event
}

// StartCalibration
type StartCalibration struct {
	Event
	Mount string `json:"Mount"`
}

// CalibrationFailed is sent when calibration fails.
type CalibrationFailed struct {
	Event
	Reason string `json:"Reason"`
}

// CalibrationDataFlipped is sent when calibration data is flipped.
type CalibrationDataFlipped struct {
	Event
	Mount string `json:"Mount"`
}

// LoopingExposures is sent for each exposure frame while looping exposures.
type LoopingExposures struct {
	Event
	Frame int `json:"Frame"`
}

// LoopingExposuresStopped is sent when looping exposures has stopped.
type LoopingExposuresStopped struct {
	Event
}

// SettleBegin is sent when settling begins after a dither or guide operation.
type SettleBegin struct {
	Event
}

// Settling is sent for each exposure frame after a dither or guide operation
// until guiding has settled.
type Settling struct {
	Event
	Distance   int     `json:"Distance"`
	Time       float64 `json:"Time"`
	SettleTime int     `json:"SettleTime"`
	StarLocked bool    `json:"StarLocked"`
}

// SettleDone is sent after a dither or guide operation indicating whether
// settling was achieved, or if the guider failed to settle before the time
// limit was reached, or if some other error occurred preventing guide or
// dither to complete and settle.
type SettleDone struct {
	Event
	Status        int    `json:"Status"`
	Error         string `json:"Error"`
	TotalFrames   int    `json:"TotalFrames"`
	DroppedFrames int    `json:"DroppedFrames"`
}

// StarLost is sent when a frame has been dropped due to the star being lost.
type StarLost struct {
	Event
	Frame     int     `json:"Frame"`
	Time      float64 `json:"Time"`
	StarMass  float64 `json:"StarMass"`
	SNR       float64 `json:"SNR"`
	AvgDist   float64 `json:"AvgDist"`
	ErrorCode int     `json:"ErrorCode"`
	Status    string  `json:"Status"`
}

// GuidingStopped is sent when guiding has stopped.
type GuidingStopped struct {
	Event
}

// Resumed is sent when PHD2 has been resumed after having been paused.
type Resumed struct {
	Event
}

// GuideStep corresponds to a line in the PHD Guide Log. The event is sent for
// each frame while guiding.
type GuideStep struct {
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

// GuidingDithered is sent when the lock position has been dithered.
type GuidingDithered struct {
	Event
	DX int `json:"dx"`
	DY int `json:"dy"`
}

// LockPositionLost is sent when the lock position has been lost.
type LockPositionLost struct {
	Event
}

// Alert is sent when an alert message was displayed in PHD2.
type Alert struct {
	Event
	Msg  string `json:"Msg"`
	Type string `json:"Type"`
}

// GuideParamChange is sent when a guiding parameter has been changed.
type GuideParamChange struct {
	Event
	Name  string `json:"Name"`
	Value string `json:"Value"`
}
