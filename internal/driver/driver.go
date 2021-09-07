package driver

import (
	"fmt"
	"sync"

	dsModels "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/egpwg/bme280-driver/pkg/device"
	bDriver "github.com/egpwg/bme280-driver/pkg/driver"
	"github.com/egpwg/bme280-driver/pkg/driver/i2c"
)

type Driver struct {
	lc       logger.LoggingClient
	asyncCh  chan<- *dsModels.AsyncValues
	deviceCh chan<- []dsModels.DiscoveredDevice
	drive    *device.Bme280
	bus      i2c.Bus
}

var driver *Driver
var once sync.Once

// NewProtocolDriver initializes the singleton Driver and
// returns it to the caller
func NewProtocolDriver() dsModels.ProtocolDriver {
	once.Do(func() {
		driver = new(Driver)
	})
	return driver
}

// Initialize performs protocol-specific initialization for the device service.
func (s *Driver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues, deviceCh chan<- []dsModels.DiscoveredDevice) error {

	s.lc = lc
	s.asyncCh = asyncCh
	s.deviceCh = deviceCh

	info, err := bDriver.GetDriverInfo()
	if err != nil {
		return err
	}

	var name string
	for _, v := range info {
		if v.Driver == "i2c-driver" {
			name = v.Bus[0]
		}
	}

	bus, err := i2c.Open(name)
	if err != nil {
		return err
	}
	s.bus = bus
	// defer bus.Close()

	bme280, err := device.NewDevice(bus)
	if err != nil {
		return err
	}
	s.drive = bme280

	// 3 is indoor mode
	err = s.drive.SetUserMode(device.UserMode(3))
	if err != nil {
		return err
	}

	return nil
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *Driver) HandleReadCommands(deviceName string,
	protocols map[string]contract.ProtocolProperties,
	reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {

	res = make([]*dsModels.CommandValue, len(reqs))

	for i, req := range reqs {
		switch req.DeviceResourceName {
		case "AllSensors":
			// get sensor data and add to res
			rList, err := s.getAllSensorData()
			if err != nil {
				s.lc.Error(err.Error())
				return res, err
			}
			cv, err := dsModels.NewFloat32ArrayValue(req.DeviceResourceName, 0, rList)
			if err != nil {
				s.lc.Error(err.Error())
				return res, err
			}
			res[i] = cv
		}
	}
	return res, nil
}

// HandleWriteCommands passes a slice of CommandRequest struct each representing
// a ResourceOperation for a specific device resource.
// Since the commands are actuation commands, params provide parameters for the individual
// command.
func (s *Driver) HandleWriteCommands(deviceName string,
	protocols map[string]contract.ProtocolProperties,
	reqs []dsModels.CommandRequest,
	params []*dsModels.CommandValue) error {

	return nil
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (s *Driver) Stop(force bool) error {
	// Then Logging Client might not be initialized
	if s.lc != nil {
		s.lc.Debug(fmt.Sprintf("Driver.Stop called: force=%v", force))
	}
	return nil
}

// AddDevice is a callback function that is invoked
// when a new Device associated with this Device Service is added
func (s *Driver) AddDevice(deviceName string,
	protocols map[string]contract.ProtocolProperties,
	adminState contract.AdminState) error {
	s.lc.Debug(fmt.Sprintf("a new Device is added: %s", deviceName))
	return nil
}

// UpdateDevice is a callback function that is invoked
// when a Device associated with this Device Service is updated
func (s *Driver) UpdateDevice(deviceName string,
	protocols map[string]contract.ProtocolProperties,
	adminState contract.AdminState) error {
	s.lc.Debug(fmt.Sprintf("Device %s is updated", deviceName))
	return nil
}

// RemoveDevice is a callback function that is invoked
// when a Device associated with this Device Service is removed
func (s *Driver) RemoveDevice(deviceName string,
	protocols map[string]contract.ProtocolProperties) error {
	s.lc.Debug(fmt.Sprintf("Device %s is removed", deviceName))
	return nil
}

// Discover triggers protocol specific device discovery, which is an asynchronous operation.
// Devices found as part of this discovery operation are written to the channel devices.
func (s *Driver) Discover() {}

func (s *Driver) getAllSensorData() ([]float32, error) {
	// # 1. Temperature
	// # 2. Pressure
	// # 3. Humidity
	r := make([]float32, 3)
	all, err := s.drive.GetSenseValue()
	if err != nil {
		return nil, err
	}
	r[0] = all.Temperature
	r[1] = all.Pressure
	r[2] = all.Humidity
	return r, nil
}
