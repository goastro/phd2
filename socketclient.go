package phd2

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

// SocketStatus represents the current status of PHD2.
type SocketStatus byte

const (
	// SocketStatusIdle is not paused, looping, or guiding.
	SocketStatusIdle = SocketStatus(0)
	// SocketStatusStarSelected is capture active and star selected.
	SocketStatusStarSelected = SocketStatus(1)
	// SocketStatusCalibrating is running the calibration routine.
	SocketStatusCalibrating = SocketStatus(2)
	// SocketStatusGuiding is guiding and locked onto a star.
	SocketStatusGuiding = SocketStatus(3)
	// SocketStatusStarLost is guiding but star lost.
	SocketStatusStarLost = SocketStatus(4)
	// SocketStatusPaused is paused.
	SocketStatusPaused = SocketStatus(100)
	// SocketStatusLooping is looping but no star selected.
	SocketStatusLooping = SocketStatus(101)
)

// SocketDitherAmount is the amount to dither by.
type SocketDitherAmount byte

const (
	// SocketDitherAmountTiny is +/- 0.5 x dither scale.
	SocketDitherAmountTiny = SocketDitherAmount(3)
	// SocketDitherAmountSmall is +/- 1.0 x dither scale.
	SocketDitherAmountSmall = SocketDitherAmount(4)
	// SocketDitherAmountNormal is +/- 2.0 x dither scale.
	SocketDitherAmountNormal = SocketDitherAmount(5)
	// SocketDitherAmountLarge is +/- 3.0 x dither scale.
	SocketDitherAmountLarge = SocketDitherAmount(12)
	// SocketDitherAmountHuge is +/- 5.0 x dither scale.
	SocketDitherAmountHuge = SocketDitherAmount(13)
)

// Dialer is the interfaced used to connect to the PHD2 server. net.Dialer will
// satisfy this interface.
type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

// SocketClient represents the connection to the PHD2 server. See
// https://github.com/OpenPHDGuiding/phd2/wiki/SocketServerInterface for
// documentation on the SocketClient interface.
type SocketClient struct {
	d Dialer
	c net.Conn
}

// NewSocketClient creates a new client to interface with the PHD2 server.
func NewSocketClient(d Dialer) *SocketClient {
	return &SocketClient{
		d: d,
	}
}

// Connect will use the Dialer to connect to the PHD2 server.
func (c *SocketClient) Connect(host string, port int) error {
	var err error
	c.c, err = c.d.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	return errors.Wrap(err, "error connecting to phd2")
}

// Close will close the underlying client connection.
func (c *SocketClient) Close() error {
	if c.c == nil {
		return ErrNotConnected
	}

	err := c.c.Close()
	return errors.Wrap(err, "error closing connection")
}

// Pause pauses guiding. Camera exposures continue to loop if they are already
// looping.
func (c *SocketClient) Pause() error {
	if c.c == nil {
		return ErrNotConnected
	}

	_, err := c.c.Write([]byte{1})
	if err != nil {
		return errors.Wrap(err, "error sending command")
	}

	resp := make([]byte, 1)
	i, err := c.c.Read(resp)
	if err != nil {
		return errors.Wrap(err, "error reading response")
	}

	if i != 1 || len(resp) != 1 || resp[0] != 0 {
		return errors.New("unexpected response")
	}

	return nil
}

// Resume resumes guiding if it was paused, otherwise no effect.
func (c *SocketClient) Resume() error {
	if c.c == nil {
		return ErrNotConnected
	}

	_, err := c.c.Write([]byte{2})
	if err != nil {
		return errors.Wrap(err, "error sending command")
	}

	resp := make([]byte, 1)
	i, err := c.c.Read(resp)
	if err != nil {
		return errors.Wrap(err, "error reading response")
	}

	if i != 1 || len(resp) != 1 || resp[0] != 0 {
		return errors.New("unexpected response")
	}

	return nil
}

// Stop stops looping exposures or guiding. SocketClient should poll with GetStatus
// to check that looping/guiding has actually stopped.
func (c *SocketClient) Stop() error {
	if c.c == nil {
		return ErrNotConnected
	}

	_, err := c.c.Write([]byte{18})
	if err != nil {
		return errors.Wrap(err, "error sending command")
	}

	resp := make([]byte, 1)
	i, err := c.c.Read(resp)
	if err != nil {
		return errors.Wrap(err, "error reading response")
	}

	if i != 1 || len(resp) != 1 || resp[0] != 0 {
		return errors.New("unexpected response")
	}

	return nil
}

// StartGuiding starts guiding. SocketClient should poll with GetStatus to check that
// guiding has actually started.
func (c *SocketClient) StartGuiding() error {
	if c.c == nil {
		return ErrNotConnected
	}

	_, err := c.c.Write([]byte{20})
	if err != nil {
		return errors.Wrap(err, "error sending command")
	}

	resp := make([]byte, 1)
	i, err := c.c.Read(resp)
	if err != nil {
		return errors.Wrap(err, "error reading response")
	}

	if i != 1 || len(resp) != 1 || resp[0] != 0 {
		return errors.New("unexpected response")
	}

	return nil
}

// ClearCalibration clears calibration data (force re-calibration).
func (c *SocketClient) ClearCalibration() error {
	if c.c == nil {
		return ErrNotConnected
	}

	_, err := c.c.Write([]byte{22})
	if err != nil {
		return errors.Wrap(err, "error sending command")
	}

	resp := make([]byte, 1)
	i, err := c.c.Read(resp)
	if err != nil {
		return errors.Wrap(err, "error reading response")
	}

	if i != 1 || len(resp) != 1 || resp[0] != 0 {
		return errors.New("unexpected response")
	}

	return nil
}

// Deselect de-selects the currently selected guide star. If subframes are
// enabled, switch to full frames. This command should be sent before sending
// AutoFindStar to ensure a full frame is captured. For example, the following
// sequence could be used to select a guide star: Stop, Deselect, Loop,
// LoopFrameCount, AutoFindStar.
func (c *SocketClient) Deselect() error {
	if c.c == nil {
		return ErrNotConnected
	}

	_, err := c.c.Write([]byte{24})
	if err != nil {
		return errors.Wrap(err, "error sending command")
	}

	resp := make([]byte, 1)
	i, err := c.c.Read(resp)
	if err != nil {
		return errors.Wrap(err, "error reading response")
	}

	if i != 1 || len(resp) != 1 || resp[0] != 0 {
		return errors.New("unexpected response")
	}

	return nil
}

// Loop starts looping exposures. SocketClient should poll with GetStatus to see if
// looping actually started.
func (c *SocketClient) Loop() (bool, error) {
	if c.c == nil {
		return false, ErrNotConnected
	}

	_, err := c.c.Write([]byte{19})
	if err != nil {
		return false, errors.Wrap(err, "error sending command")
	}

	resp := make([]byte, 1)
	i, err := c.c.Read(resp)
	if err != nil {
		return false, errors.Wrap(err, "error reading response")
	}

	if i != 1 || len(resp) != 1 {
		return false, errors.New("unexpected response")
	}

	if resp[0] == 0 {
		return true, nil
	}

	return false, nil
}

// GetStatus gets a value describing the state of PHD.
func (c *SocketClient) GetStatus() (SocketStatus, error) {
	if c.c == nil {
		return SocketStatusIdle, ErrNotConnected
	}

	_, err := c.c.Write([]byte{17})
	if err != nil {
		return SocketStatusIdle, errors.Wrap(err, "error sending command")
	}

	resp := make([]byte, 1)
	i, err := c.c.Read(resp)
	if err != nil {
		return SocketStatusIdle, errors.Wrap(err, "error reading response")
	}

	if i != 1 || len(resp) != 1 {
		return SocketStatusIdle, errors.New("unexpected response")
	}

	return SocketStatus(resp[0]), nil
}

// Dither will dither a random amount. Returns the camera exposure time in
// seconds, but not less than 1.
func (c *SocketClient) Dither(amt SocketDitherAmount) (uint8, error) {
	if c.c == nil {
		return 0, ErrNotConnected
	}

	_, err := c.c.Write([]byte{byte(amt)})
	if err != nil {
		return 0, errors.Wrap(err, "error sending command")
	}

	resp := make([]byte, 1)
	i, err := c.c.Read(resp)
	if err != nil {
		return 0, errors.Wrap(err, "error reading response")
	}

	if i != 1 || len(resp) != 1 {
		return 0, errors.New("unexpected response")
	}

	return resp[0], nil
}

// RequestDistance requests guide error distance. Returns the current guide
// error distance in units of 1/100 pixel. Values > 255 are reported as 255.
func (c *SocketClient) RequestDistance() (uint8, error) {
	if c.c == nil {
		return 255, ErrNotConnected
	}

	_, err := c.c.Write([]byte{10})
	if err != nil {
		return 255, errors.Wrap(err, "error sending command")
	}

	resp := make([]byte, 1)
	i, err := c.c.Read(resp)
	if err != nil {
		return 255, errors.Wrap(err, "error reading response")
	}

	if i != 1 || len(resp) != 1 {
		return 255, errors.New("unexpected response")
	}

	return resp[0], nil
}

// LoopFrameCount gets the current frame counter value.	Returns 0 if not
// looping or guiding. Otherwise, the current frame counter value (capped at
// 255). The frame counter is incremented for each camera exposure when looping
// or guiding.
func (c *SocketClient) LoopFrameCount() (uint8, error) {
	if c.c == nil {
		return 0, ErrNotConnected
	}

	_, err := c.c.Write([]byte{21})
	if err != nil {
		return 0, errors.Wrap(err, "error sending command")
	}

	resp := make([]byte, 1)
	i, err := c.c.Read(resp)
	if err != nil {
		return 0, errors.Wrap(err, "error reading response")
	}

	if i != 1 || len(resp) != 1 {
		return 0, errors.New("unexpected response")
	}

	return resp[0], nil
}

// AutoFindStar auto-selects a guide star.
func (c *SocketClient) AutoFindStar() (bool, error) {
	if c.c == nil {
		return false, ErrNotConnected
	}

	_, err := c.c.Write([]byte{14})
	if err != nil {
		return false, errors.Wrap(err, "error sending command")
	}

	resp := make([]byte, 1)
	i, err := c.c.Read(resp)
	if err != nil {
		return false, errors.Wrap(err, "error reading response")
	}

	if i != 1 || len(resp) != 1 {
		return false, errors.New("unexpected response")
	}

	if resp[0] == 0 {
		return true, nil
	}

	return false, nil
}

// SetLockPosition sets the lock position to (x,y). This is not yet implemented.
func (c *SocketClient) SetLockPosition(x, y uint16) error {
	return ErrNotImplemented
}

// FlipRACalibrationData flips the RA calibration data.
func (c *SocketClient) FlipRACalibrationData() (bool, error) {
	if c.c == nil {
		return false, ErrNotConnected
	}

	_, err := c.c.Write([]byte{16})
	if err != nil {
		return false, errors.Wrap(err, "error sending command")
	}

	resp := make([]byte, 1)
	i, err := c.c.Read(resp)
	if err != nil {
		return false, errors.Wrap(err, "error reading response")
	}

	if i != 1 || len(resp) != 1 {
		return false, errors.New("unexpected response")
	}

	if resp[0] == 1 {
		return true, nil
	}

	return false, nil
}
