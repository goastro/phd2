package phd2

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/pkg/errors"
)

type RPCClient struct {
	d    Dialer
	conn net.Conn

	reader *bufio.Reader
	writer *bufio.Writer

	// Ensures only one RPC method call is active at a time.
	methodMutex    sync.Mutex
	methodResponse chan []byte
	requestID      int

	events chan []byte
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

	c.events = make(chan []byte, 10)
	c.methodResponse = make(chan []byte, 1)

	c.reader = bufio.NewReader(c.conn)
	c.writer = bufio.NewWriter(c.conn)

	go c.processReadLoop()
	go c.readLoop()

	return nil
}

func (c *RPCClient) readLoop() {
	var line []byte
	var partial bool
	var err error

	for {
		var bytes []byte

		bytes, partial, err = c.reader.ReadLine()
		if err != nil {
			println(err.Error())
			return
		}

		line = append(line, bytes...)

		if !partial {
			c.events <- line
			line = nil
		}
	}
}

func (c *RPCClient) processReadLoop() {
	for line := range c.events {
		var evt Event
		err := json.Unmarshal(line, &evt)
		if err != nil {
			println(err.Error())
			continue
		}

		resp, ok := getEvent(evt.Event)
		if !ok {
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

type rpcRequest struct {
	Method string        `json:"method"`
	ID     int           `json:"id"`
	Params []interface{} `json:"params,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type rpcResponse struct {
	ID     int             `json:"id"`
	Error  rpcError        `json:"error,omitempty"`
	Result json.RawMessage `json:"result"`
}

func (c *RPCClient) call(name string, params []interface{}, result interface{}) (*rpcResponse, error) {
	c.methodMutex.Lock()
	defer c.methodMutex.Unlock()

	c.requestID++

	req := rpcRequest{
		Method: name,
		ID:     c.requestID,
		Params: params,
	}

	bytes, err := json.Marshal(&req)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling request")
	}

	bytes = append(bytes, '\r', '\n')

	_, err = c.writer.Write(bytes)
	if err != nil {
		return nil, errors.Wrap(err, "error writing to connection")
	}

	err = c.writer.Flush()
	if err != nil {
		return nil, errors.Wrap(err, "error flushing write")
	}

	line := <-c.methodResponse

	var resp rpcResponse

	err = json.Unmarshal(line, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling response")
	}

	if resp.ID != req.ID {
		return nil, errors.New("incorrect response received")
	}

	err = json.Unmarshal(resp.Result, result)

	return &resp, errors.Wrap(err, "error unmarshalling result")
}

// https://github.com/OpenPHDGuiding/phd2/wiki/EventMonitoring#available-methods
// https://github.com/OpenPHDGuiding/phd2/blob/master/event_server.cpp

func (c *RPCClient) GetExposure() (int, error) {
	var result int
	_, err := c.call("get_exposure", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) CaptureSingleFrame(duration int, subframe []int) error {
	var result int
	_, err := c.call("capture_single_frame", []interface{}{duration, subframe}, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) ClearCalibration(which string) error {
	var result int
	_, err := c.call("clear_calibration", []interface{}{which}, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

type Settle struct {
	Pixels         float64 `json:"pixels"`
	TimeSeconds    int     `json:"time"`
	TimeoutSeconds int     `json:"timeout"`
}

func (c *RPCClient) Dither(pixels float64, raOnly bool, settle Settle) error {
	var result int
	_, err := c.call("dither", []interface{}{pixels, raOnly, settle}, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) FindStar() ([]float64, error) {
	var result []float64
	_, err := c.call("find_star", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) FlipCalibration() error {
	var result int
	_, err := c.call("flip_calibration", nil, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GetAppState() (string, error) {
	var result string
	_, err := c.call("get_app_state", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GetCalibrated() (bool, error) {
	var result bool
	_, err := c.call("get_calibrated", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GetConnected() (bool, error) {
	var result bool
	_, err := c.call("get_connected", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GetAlgorithmParamNames(axis string) ([]string, error) {
	var result []string
	_, err := c.call("get_algo_param_names", []interface{}{axis}, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GetAlgorithmParam(axis, param string) (float64, error) {
	var result float64
	_, err := c.call("get_algo_param", []interface{}{axis, param}, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

type CalibrationData struct {
	Calibrated bool    `json:"calibrated"`
	XAngle     float64 `json:"xAngle"`
	XRate      float64 `json:"xRate"`
	XParity    string  `json:"xParity"`
	YAngle     float64 `json:"yAngle"`
	YRate      float64 `json:"yRate"`
	YParity    string  `json:"yParity"`
}

func (c *RPCClient) GetCalibrationData(which string) (CalibrationData, error) {
	var result CalibrationData
	_, err := c.call("get_calibration_data", []interface{}{which}, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

type CoolerStatus struct {
	Temperature float64 `json:"temperature"`
	CoolerOn    bool    `json:"coolerOn"`
	Setpoint    float64 `json:"setpoint"`
	Power       float64 `json:"power"`
}

func (c *RPCClient) GetCoolerStatus() (CoolerStatus, error) {
	var result CoolerStatus
	_, err := c.call("get_cooler_status", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

type Equipment struct {
	Name      string `json:"name"`
	Connected bool   `json:"connected"`
}

type CurrentEquipment struct {
	Camera   Equipment `json:"camera"`
	Mount    Equipment `json:"mount"`
	AuxMount Equipment `json:"aux_mount"`
	AO       Equipment `json:"AO"`
	Rotator  Equipment `json:"rotator"`
}

func (c *RPCClient) GetCurrentEquipment() (CurrentEquipment, error) {
	var result CurrentEquipment
	_, err := c.call("get_algo_param", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GetDecGuidMode() (string, error) {
	var result string
	_, err := c.call("get_dec_guide_mode", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GetExposureDurations() ([]int, error) {
	var result []int
	_, err := c.call("get_exposure_durations", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GetLockPosition() ([]int, error) {
	var result []int
	_, err := c.call("get_lock_position", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GetLockShiftEnabled() (bool, error) {
	var result bool
	_, err := c.call("get_lock_shift_enabled", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

type LockShiftParams struct {
	Enabled bool      `json:"enabled"`
	Rate    []float64 `json:"rate"`
	Units   string    `json:"units"`
	Axes    string    `json:"axes"`
}

func (c *RPCClient) GetLockShiftParams() (LockShiftParams, error) {
	var result LockShiftParams
	_, err := c.call("get_lock_shift_params", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GetPaused() (bool, error) {
	var result bool
	_, err := c.call("get_paused", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GetPixelScale() (float64, error) {
	var result float64
	_, err := c.call("get_pixel_scale", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

type Profile struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *RPCClient) GetProfile() (Profile, error) {
	var result Profile
	_, err := c.call("get_profile", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GetProfiles() ([]Profile, error) {
	var result []Profile
	_, err := c.call("get_profiles", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GetSearchRegion() (int, error) {
	var result int
	_, err := c.call("get_search_region", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GetSensorTemperature() (float64, error) {
	var result float64
	_, err := c.call("get_sensor_temperature", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GetStarImage() error {
	return ErrNotImplemented
}

func (c *RPCClient) GetUseSubframes() (bool, error) {
	var result bool
	_, err := c.call("get_use_subframes", nil, &result)
	return result, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) Guide(settle Settle, recalibrate bool) error {
	var result int
	_, err := c.call("guide", []interface{}{settle, recalibrate}, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) GuidePulse(amount int, direction, which string) error {
	var result int
	_, err := c.call("guide_pulse", []interface{}{amount, direction, which}, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) Loop() error {
	var result int
	_, err := c.call("loop", nil, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) SaveImage() (string, error) {
	var result struct {
		Filename string `json:"filename"`
	}
	_, err := c.call("save_image", nil, &result)
	return result.Filename, errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) SetAlgorithmParam(axis, name string, value float64) error {
	var result int
	_, err := c.call("set_algo_param", []interface{}{axis, name, value}, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) SetConnected(connect bool) error {
	var result int
	_, err := c.call("set_algo_param", []interface{}{connect}, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) SetDecGuideMode(mode string) error {
	var result int
	_, err := c.call("set_dec_guide_mode", []interface{}{mode}, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) SetExposure(length int) error {
	var result int
	_, err := c.call("set_exposure", []interface{}{length}, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) SetLockPosition(x, y float64, exact bool) error {
	var result int
	_, err := c.call("set_lock_position", []interface{}{x, y, exact}, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) SetLockShiftEnabled(enable bool) error {
	var result int
	_, err := c.call("set_lock_shift_enabled", []interface{}{enable}, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) SetLockShiftParams(params LockShiftParams) error {
	var result int
	_, err := c.call("set_lock_shift_params", []interface{}{params}, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) SetPaused(paused, full bool) error {
	var result int

	params := []interface{}{paused}
	if full {
		params = append(params, "full")
	}

	_, err := c.call("set_paused", params, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) SetProfile(id int) error {
	var result int
	_, err := c.call("set_profile", []interface{}{id}, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) Shutdown() error {
	var result int
	_, err := c.call("shutdown", nil, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}

func (c *RPCClient) StopCapture() error {
	var result int
	_, err := c.call("stop_capture", nil, &result)
	return errors.Wrap(err, "error calling jsonrpc method")
}
