package main

import (
	"bytes"
	"flag"
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"time"
	"github.com/pkg/errors"
)

const (
	// movement UX
	easeFactor = 300
	movementSleep = time.Millisecond * 20

	// origin on laptop screen
	CX = 250
	CY = 250

	// math for circular movement
	step = 2 * math.Pi / 20
	h    = 150
	k    = 150
	r    = 50
)

var (
	timeoutStr string
)

func init() {
	flag.StringVar(&timeoutStr, "timeout", "", "(optional) can pass in how long you want to run this for")
}

func main() {
	flag.Parse()

	fmt.Println("Press [CTRL+C] to stop moving mouse...")

	if timeoutStr == "" {
		loopForever()
		return
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		panic(err)
	}

	runWithTimeout(timeout)
}

func runWithTimeout(timeout time.Duration) {
	go loopForever()
	select {
	case <-time.After(timeout):
		fmt.Println("finished experiment hehe :)")
	}
}

func loopForever() {
	for {
		circleAlgo()
	}
}

// Adapted from: https://www.mathopenref.com/coordcirclealgorithm.html
func circleAlgo() {
	theta := 0.0
	for theta < 2*math.Pi {
		x := h + r*math.Cos(theta)
		y := k - r*math.Sin(theta)
		time.Sleep(movementSleep)
		if err := moveMouse(x, y); err != nil {
			panic(err)
		}
		theta += step
	}
}

// Installed from: https://github.com/BlueM/cliclick
func moveMouse(x, y float64) error {
	cmd := exec.Command(
		"cliclick",
		"-e", strconv.Itoa(easeFactor),
		fmt.Sprintf("m:%d,%d", round(CX+x), round(CY+y)),
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, string(stderr.Bytes()))
	}

	return nil
}

func round(f float64) int { return int(math.Round(f)) }