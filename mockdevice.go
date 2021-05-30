package elmobd

import (
	"fmt"
	"strings"
	"time"
)

/*==============================================================================
 * External
 */

// MockResult represents the raw text output of running a raw command,
// including information used in debugging to show what input caused what
// error, how long the command took, etc.
type MockResult struct {
	input     string
	outputs   []string
	error     error
	writeTime time.Duration
	readTime  time.Duration
	totalTime time.Duration
}

// Failed checks if the result is successful or not
func (res *MockResult) Failed() bool {
	return res.error != nil
}

// GetError returns the results current error
func (res *MockResult) GetError() error {
	return res.error
}

// GetOutputs returns the outputs of the result
func (res *MockResult) GetOutputs() []string {
	return res.outputs
}

// FormatOverview formats a result as an overview of what command was run and
// how long it took.
func (res *MockResult) FormatOverview() string {
	lines := []string{
		"=======================================",
		" Mocked command \"%s\"",
		"=======================================",
	}

	return fmt.Sprintf(
		strings.Join(lines, "\n"),
		res.input,
	)
}

// MockDevice represent a mocked serial connection
type MockDevice struct {
}

// RunCommand mocks the given AT/OBD command by just returning a result for the
// mocked outputs set earlier.
func (dev *MockDevice) RunCommand(command string) RawResult {
	return &MockResult{
		input:     command,
		outputs:   mockOutputs(command),
		writeTime: 0,
		readTime:  0,
		totalTime: 0,
	}
}

// Close closes mock device
func (dev *MockDevice) Close() error {
	return nil
}

/*==============================================================================
 * Internal
 */

func mockMode1Outputs(subcmd string) []string {

	if len(subcmd) < 2 {
		return []string{"NOT SUPPORTED"}
	}

	// First 2 charaters signify the type of command
	switch strings.ToUpper(subcmd[:2]) {
	case "00":
		// PIDs supported part 1
		return []string{
			"41 00 0C 10 00 00", // Means PIDs supported: 05, 06, 0C
		}
	case "20": // PIDs supported part 2
		return []string{
			"41 20 00 00 00 00",
		}
	case "40": // PIDs supported part 3
		return []string{
			"41 40 00 00 00 00",
		}
	case "60": // PIDs supported part 4
		return []string{
			"41 60 00 00 00 00",
		}
	case "80": // PIDs supported part 5
		return []string{
			"41 80 00 00 00 00",
		}
	case "01":
		return []string{
			"41 01 FF 00 00 00",
		}
	//Coolant Temp
	case "05":
		return []string{
			"41 05 4F",
		}
	case "06":
		return []string{
			"41 06 02",
		}
	//Engine RPM
	case "0C":
		return []string{
			"41 0C 03 00",
		}
	//Fuel Level
	case "2F":
		return []string{
			"41 2F 6B",
		}
	case "OD":
		return []string{
			"41 0D 4B",
		}
	case "31":
		return []string{
			"41 31 02 0C",
		}
	//Engine Load
	case "04":
		return []string{
			"41 04 50",
		}
	//Fuel Pressure
	case "0A":
		return []string{
			"41 0A 50",
		}
	//Speed (km/h)
	case "0D":
		return []string{
			"41 0D 0F",
		}
	//Intake Air temp
	case "0F":
		return []string{
			"41 0F 3C",
		}
	}

	//default case
	return []string{"NOT SUPPORTED"}

}

func mockOutputs(cmd string) []string {
	if cmd == "ATSP0" {
		return []string{"OK"}
	} else if cmd == "AT@1" {
		return []string{"OBDII by elm329@gmail.com"}
	} else if strings.HasPrefix(cmd, "01") {
		return mockMode1Outputs(cmd[2:])
	}

	return []string{"NOT SUPPORTED"}
}
