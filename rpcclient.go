package phd2

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type RPCClient struct {
	d    Dialer
	conn net.Conn

	reader *bufio.Reader
	writer *bufio.Writer

	events chan string
}

func NewRPCClient(d Dialer) *RPCClient {
	return &RPCClient{
		d: d,
	}
}

func (c *RPCClient) Connect(host string, port int) error {
	var err error

	c.conn, err = c.d.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	c.events = make(chan string, 10)

	c.reader = bufio.NewReader(c.conn)
	c.writer = bufio.NewWriter(c.conn)

	go c.processReadLoop()
	go c.readLoop()

	return nil
}

func (c *RPCClient) readLoop() {
	var line string
	var partial bool
	var err error

	for {
		existingLine := line
		var bytes []byte

		bytes, partial, err = c.reader.ReadLine()
		if err != nil {
			println(err.Error())
			return
		}

		line = existingLine + string(bytes)

		if !partial {
			c.events <- line
			line = ""
		}
	}
}

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

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RPCResponse struct {
	ID    int      `json:"id"`
	Error RPCError `json:"error,omitempty"`
}

type CaptureSingleFrameResponse struct {
	RPCResponse
}

type ClearCalibrationResponse struct {
	RPCResponse
}

type DitherResponse struct {
	RPCResponse
}

func (c *RPCClient) processReadLoop() {
	for line := range c.events {
		var evt Event
		err := json.Unmarshal([]byte(line), &evt)
		if err != nil {
			println(err.Error())
			continue
		}

		var resp interface{}

		switch evt.Event {
		case "Version":
			resp = &Version{}
		case "CalibrationComplete":
			resp = &CalibrationComplete{}
		case "Paused":
			resp = &Paused{}
		case "AppState":
			resp = &AppState{}
		default:
			println(fmt.Sprintf("uknown event: %s", line))
		}

		err = json.Unmarshal([]byte(line), resp)
		if err != nil {
			println(err.Error())
			continue
		}

		println(fmt.Sprintf("%#v", resp))
	}
}
