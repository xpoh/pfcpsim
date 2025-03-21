// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package pfcpsim

import (
	"github.com/xpoh/pfcpsim/pkg/pfcpsim"
)

var (
	remotePeerAddress string
	upfN3Address      string

	interfaceName string

	// Emulates 5G SMF/ 4G SGW
	sim                 *pfcpsim.PFCPClient
	remotePeerConnected bool
)
