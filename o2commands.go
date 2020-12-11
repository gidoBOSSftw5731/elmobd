package elmobd

import "fmt"

// OxygenSensorPresence represents a command that checks oxygen sensors presence.
//
// A0..A3 -> Bank 1, sensors 1-4
// A4..A7 -> Bank 2, sensors 1-4
type OxygenSensorPresence struct {
	baseCommand
	IntCommand
}

// NewOxygenSensorPresence creates a new NewOxygenSensorsPresence with the right parameters.
func NewOxygenSensorPresence() *OxygenSensorPresence {
	return &OxygenSensorPresence{
		baseCommand{SERVICE_01_ID, 19, 1, "oxygen_sensor_presence", Ready},
		IntCommand{},
	}
}

// SetValue processes the byte array value into the right byte value.
func (cmd *OxygenSensorPresence) SetValue(result *Result) error {
	payload, err := result.PayloadAsByte()

	if err != nil {
		return err
	}

	cmd.Value = int(payload)

	return nil
}

type baseOxygenSensorCommand struct {
	baseCommand
	SensorID byte
}

type fuelTrimVoltage struct {
	Voltage           float32
	ShortTermFuelTrim float32
}

func (cmd *fuelTrimVoltage) SetValue(result *Result) error {
	stftCoef := float32(.78125)

	cmd.Voltage = float32(result.value[2]) / 200
	cmd.ShortTermFuelTrim = float32(result.value[3])*stftCoef - 100

	return nil
}

func (cmd *fuelTrimVoltage) ValueAsLit() string {
	return fmt.Sprintf("Voltage: %f, ShortTermFuelTrim: %f", cmd.Voltage, cmd.ShortTermFuelTrim)
}

// OxygenSensorFuelTrim represents a command that checks fuel trim values for oxygen sensor.
type OxygenSensorFuelTrim struct {
	baseOxygenSensorCommand
	fuelTrimVoltage
}

// NewOxygenSensorFuelTrim creates a new OxygenSensorFuelTrim with the right parameters.
func NewOxygenSensorFuelTrim(sensorID byte) *OxygenSensorFuelTrim {
	pid := OBDParameterID(19 + sensorID)
	return &OxygenSensorFuelTrim{
		baseOxygenSensorCommand{
			baseCommand{SERVICE_01_ID, pid, 2, fmt.Sprintf("oxygen_sensor%d_fuel_trim", sensorID), Ready},
			sensorID,
		},
		fuelTrimVoltage{},
	}
}

type lambdaVoltage struct {
	Lambda  float32
	Voltage float32
}

func (cmd *lambdaVoltage) SetValue(result *Result) error {
	payload, err := result.PayloadAsUInt32()

	if err != nil {
		return err
	}

	cmd.Lambda = float32(payload) // TODO fix
	cmd.Voltage = float32(0xFF00 ^ payload)

	return nil
}

func (cmd *lambdaVoltage) ValueAsLit() string {
	return fmt.Sprintf("LambdaRatio: %f, Voltage:%f", cmd.Lambda, cmd.Voltage)
}

// OxygenSensorLambdaVoltage represents a command that checks lambda values for oxygen sensor.
type OxygenSensorLambdaVoltage struct {
	baseOxygenSensorCommand
	lambdaVoltage
}

// NewOxygenSensorLambdaVoltage creates a new OxygenSensorLambdaVoltage with the right parameters.
func NewOxygenSensorLambdaVoltage(sensorID byte) *OxygenSensorLambdaVoltage {
	pid := OBDParameterID(35 + sensorID)
	return &OxygenSensorLambdaVoltage{
		baseOxygenSensorCommand{
			baseCommand{SERVICE_01_ID, pid, 2, fmt.Sprintf("oxygen_sensor%d_lambda_voltage", sensorID), Ready},
			sensorID,
		},
		lambdaVoltage{},
	}
}

type lambdaCurrent struct {
	Lambda  float32
	Current float32
}

func (cmd *lambdaCurrent) SetValue(result *Result) error {
	payload, err := result.PayloadAsUInt32()

	if err != nil {
		return err
	}

	cmd.Lambda = float32(0x00FF ^ payload) // TODO fix
	cmd.Current = float32((0xFF00 ^ payload) >> 8)

	return nil
}

func (cmd *lambdaCurrent) ValueAsLit() string {
	return fmt.Sprintf("LambdaRatio: %f, Current:%f", cmd.Lambda, cmd.Current)
}

// OxygenSensorLambdaCurrent represents a command that checks lambda values for oxygen sensor.
type OxygenSensorLambdaCurrent struct {
	baseOxygenSensorCommand
	lambdaCurrent
}

// NewOxygenSensorLambdaCurrent creates a new OxygenSensorLambdaCurrent with the right parameters.
func NewOxygenSensorLambdaCurrent(sensorID byte) *OxygenSensorLambdaCurrent {
	pid := OBDParameterID(51 + sensorID)
	return &OxygenSensorLambdaCurrent{
		baseOxygenSensorCommand{
			baseCommand{SERVICE_01_ID, pid, 2, fmt.Sprintf("oxygen_sensor%d_lambda_current", sensorID), Ready},
			sensorID,
		},
		lambdaCurrent{},
	}
}
