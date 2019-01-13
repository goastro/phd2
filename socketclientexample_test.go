package phd2_test

import (
	"net"
	"time"

	"github.com/goastro/phd2"
)

func ExampleSocketClient() {
	c := phd2.NewSocketClient(&net.Dialer{})

	var err error

	err = c.Connect("127.0.0.1", 4300)
	if err != nil {
		panic(err.Error())
	}

	err = c.Stop()
	if err != nil {
		panic(err.Error())
	}

	err = c.Deselect()
	if err != nil {
		panic(err.Error())
	}

	var success bool

	success, err = c.Loop()
	if err != nil {
		panic(err.Error())
	}

	if !success {
		panic("unable to start looping")
	}

	var status phd2.SocketStatus

	for status != phd2.SocketStatusLooping {
		status, err = c.GetStatus()
		if err != nil {
			panic(err.Error())
		}

		println(status.String())
		time.Sleep(100 * time.Millisecond)
	}

	var frameCount uint8

	for frameCount < 10 {
		frameCount, err = c.LoopFrameCount()
		if err != nil {
			panic(err.Error())
		}

		println(status.String())
		println(frameCount)
		time.Sleep(100 * time.Millisecond)
	}

	success, err = c.AutoFindStar()
	if err != nil {
		panic(err.Error())
	}

	if !success {
		panic("unable to find star")
	}

	for status != phd2.SocketStatusStarSelected {
		status, err = c.GetStatus()
		if err != nil {
			panic(err.Error())
		}

		println(status.String())
		time.Sleep(100 * time.Millisecond)
	}

	err = c.ClearCalibration()
	if err != nil {
		panic(err.Error())
	}

	err = c.StartGuiding()
	if err != nil {
		panic(err.Error())
	}

	for status != phd2.SocketStatusGuiding {
		status, err = c.GetStatus()
		if err != nil {
			panic(err.Error())
		}

		println(status.String())
		time.Sleep(100 * time.Millisecond)
	}

	err = c.Close()
	if err != nil {
		panic(err.Error())
	}
}
