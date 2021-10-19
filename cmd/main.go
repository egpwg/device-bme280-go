// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) EGPWG
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
	"github.com/egpwg/device-bme280-go/internal/driver"
)

const serviceName string = "bme280-device-go"

func main() {
	sd := driver.Driver{}
	startup.Bootstrap(serviceName, "0.0.1", &sd)
}
