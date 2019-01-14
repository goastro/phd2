package phd2

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type RPCClient struct {
	d    Dialer
	conn net.Conn

	reader *bufio.Reader
	writer *bufio.Writer

	// Ensures only one RPC method call is active at a time.
	methodMutex    sync.Mutex
	methodResponse chan string
	requestID      int

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
	c.methodResponse = make(chan string, 1)

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
			println(line)
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

type EmptyResponse struct {
	RPCResponse
}

type StringResponse struct {
	RPCResponse
	Result string `json:"result"`
}

type IntResponse struct {
	RPCResponse
	Result int `json:"result"`
}

type BooleanResponse struct {
	RPCResponse
	Result bool `json:"result"`
}

type FloatSliceResponse struct {
	RPCResponse
	Result []float64 `json:"result"`
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
			if len(evt.Event) == 0 {
				c.methodResponse <- line
			} else {
				println(fmt.Sprintf("unknown event: %s", line))
			}
			continue
		}

		err = json.Unmarshal([]byte(line), resp)
		if err != nil {
			println(err.Error())
			continue
		}

		println(fmt.Sprintf("%#v", resp))
	}
}

type RPCRequest struct {
	Method string        `json:"method"`
	ID     int           `json:"id"`
	Params []interface{} `json:"params,omitempty"`
}

func (c *RPCClient) call(name string, params []interface{}, resp interface{}) error {
	c.methodMutex.Lock()
	defer c.methodMutex.Unlock()

	c.requestID++

	req := RPCRequest{
		Method: name,
		ID:     c.requestID,
		Params: params,
	}

	bytes, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	_, err = c.writer.Write(bytes)
	if err != nil {
		return err
	}

	_, err = c.writer.WriteString("\r\n")
	if err != nil {
		return err
	}

	err = c.writer.Flush()
	if err != nil {
		return err
	}

	line := <-c.methodResponse

	println(line)

	err = json.Unmarshal([]byte(line), &resp)
	if err != nil {
		return err
	}

	return nil
}

// https://github.com/OpenPHDGuiding/phd2/wiki/EventMonitoring#available-methods

func (c *RPCClient) GetExposure() (int, error) {
	var resp IntResponse
	err := c.call("get_exposure", nil, &resp)
	if err != nil {
		return 0, err
	}

	return resp.Result, nil
}

func (c *RPCClient) CaptureSingleFrame(duration int, subframe []int) error {
	var resp IntResponse
	err := c.call("capture_single_frame", []interface{}{duration, subframe}, &resp)
	if err != nil {
		return err
	}

	return nil
}

func (c *RPCClient) ClearCalibration(which string) error {
	var resp IntResponse
	err := c.call("clear_calibration", []interface{}{which}, &resp)
	if err != nil {
		return err
	}

	return nil
}

type Settle struct {
	Pixels         float64 `json:"pixels"`
	TimeSeconds    int     `json:"time"`
	TimeoutSeconds int     `json:"timeout"`
}

func (c *RPCClient) Dither(pixels float64, raOnly bool, settle Settle) error {
	var resp IntResponse
	err := c.call("dither", []interface{}{pixels, raOnly, settle}, &resp)
	if err != nil {
		return err
	}

	return nil
}

func (c *RPCClient) FindStar() ([]float64, error) {
	var resp FloatSliceResponse
	err := c.call("find_star", nil, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

func (c *RPCClient) FlipCalibration() error {
	var resp IntResponse
	err := c.call("flip_calibration", nil, &resp)
	if err != nil {
		return err
	}

	return nil
}

func (c *RPCClient) GetAppState() (string, error) {
	var resp StringResponse
	err := c.call("get_app_state", nil, &resp)
	if err != nil {
		return "", err
	}

	return resp.Result, nil
}

func (c *RPCClient) GetCalibrated() (bool, error) {
	var resp BooleanResponse
	err := c.call("get_calibrated", nil, &resp)
	if err != nil {
		return false, err
	}

	return resp.Result, nil
}

func (c *RPCClient) GetConnected() (bool, error) {
	var resp BooleanResponse
	err := c.call("get_connected", nil, &resp)
	if err != nil {
		return false, err
	}

	return resp.Result, nil
}

func (c *RPCClient) GetAlgorithmParamNames(axis string) ([]string, error) {
	return nil, ErrNotImplemented
}

func (c *RPCClient) GetAlgorithmParam(axis, param string) (string, error) {
	return "", ErrNotImplemented
}

func (c *RPCClient) GetCalibrationData(which string) error {
	return ErrNotImplemented
}

func (c *RPCClient) GetCoolerStatus() error {
	return ErrNotImplemented
}

func (c *RPCClient) GetCurrentEquipment() error {
	return ErrNotImplemented
}

func (c *RPCClient) GetDecGuidMode() error {
	return ErrNotImplemented
}

func (c *RPCClient) GetExposureDurations() error {
	return ErrNotImplemented
}

func (c *RPCClient) GetLockPosition() error {
	return ErrNotImplemented
}

func (c *RPCClient) GetLockShiftEnabled() error {
	return ErrNotImplemented
}

func (c *RPCClient) GetLockShiftParams() error {
	return ErrNotImplemented
}

func (c *RPCClient) GetPaused() error {
	return ErrNotImplemented
}

func (c *RPCClient) GetPixelScale() error {
	return ErrNotImplemented
}

func (c *RPCClient) GetProfile() error {
	return ErrNotImplemented
}

func (c *RPCClient) GetProfiles() error {
	return ErrNotImplemented
}

func (c *RPCClient) GetSearchRegion() error {
	return ErrNotImplemented
}

func (c *RPCClient) GetSensorTemperature() error {
	return ErrNotImplemented
}

func (c *RPCClient) GetStarImage() error {
	return ErrNotImplemented
}

func (c *RPCClient) GetUseSubframes() error {
	return ErrNotImplemented
}

func (c *RPCClient) Guide() error {
	return ErrNotImplemented
}

func (c *RPCClient) GuidePulse() error {
	return ErrNotImplemented
}

func (c *RPCClient) Loop() error {
	return ErrNotImplemented
}

func (c *RPCClient) SaveImage() error {
	return ErrNotImplemented
}

func (c *RPCClient) SetAlgorithmParam() error {
	return ErrNotImplemented
}

func (c *RPCClient) SetConnected() error {
	return ErrNotImplemented
}

func (c *RPCClient) SetDecGuideMode() error {
	return ErrNotImplemented
}

func (c *RPCClient) SetExposure() error {
	return ErrNotImplemented
}

func (c *RPCClient) SetLockPosition() error {
	return ErrNotImplemented
}

func (c *RPCClient) SetLockShiftEnabled() error {
	return ErrNotImplemented
}

func (c *RPCClient) SetLockShiftParams() error {
	return ErrNotImplemented
}

func (c *RPCClient) SetPaused() error {
	return ErrNotImplemented
}

func (c *RPCClient) SetProfile() error {
	return ErrNotImplemented
}

func (c *RPCClient) Shutdown() error {
	return ErrNotImplemented
}

func (c *RPCClient) StopCapture() error {
	return ErrNotImplemented
}
