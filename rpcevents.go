package phd2

type Event struct {
	Event     string  `json:"Event"`
	Timestamp float64 `json:"Timestamp"`
	Host      string  `json:"Host"`
	Inst      int     `json:"Inst"`
}

type Version struct {
	Event
	PHDVersion string `json:"PHDVersion"`
	PHDSubver  string `json:"PHDSubver"`
	MsgVersion int    `json:"MsgVersion"`
}

type CalibrationComplete struct {
	Event
	Mount string `json:"Mount"`
}

type Paused struct {
	Event
}

type AppState struct {
	Event
	State string `json:"State"`
}

type LockPositionSet struct {
	Event
	X int `json:"X"`
	Y int `json:"Y"`
}

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

type StarSelected struct {
	Event
	X int `json:"X"`
	Y int `json:"Y"`
}

type StartGuiding struct {
	Event
}

type StartCalibration struct {
	Event
	Mount string `json:"Mount"`
}

type CalibrationFailed struct {
	Event
	Reason string `json:"Reason"`
}

type CalibrationDataFlipped struct {
	Event
	Mount string `json:"Mount"`
}

type LoopingExposures struct {
	Event
	Frame int `json:"Frame"`
}

type LoopingExposuresStopped struct {
	Event
}

type SettleBegin struct {
	Event
}

type Settling struct {
	Event
	Distance   int     `json:"Distance"`
	Time       float64 `json:"Time"`
	SettleTime int     `json:"SettleTime"`
	StarLocked bool    `json:"StarLocked"`
}

type SettleDone struct {
	Event
	Status        int    `json:"Status"`
	Error         string `json:"Error"`
	TotalFrames   int    `json:"TotalFrames"`
	DroppedFrames int    `json:"DroppedFrames"`
}

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

type GuidingStopped struct {
	Event
}

type Resumed struct {
	Event
}

type GuideStep struct {
	Event
}

type GuidingDithered struct {
	Event
	DX int `json:"dx"`
	DY int `json:"dy"`
}

type LockPositionLost struct {
	Event
}

type Alert struct {
	Event
	Msg  string `json:"Msg"`
	Type string `json:"Type"`
}

type GuideParamChange struct {
	Event
	Name  string `json:"Name"`
	Value string `json:"Value"`
}
