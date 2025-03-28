// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package session

import (
	"github.com/wmnsk/go-pfcp/ie"
	"github.com/xpoh/pfcpsim/logger"
)

type farBuilder struct {
	farID        uint32
	applyAction  uint8
	method       IEMethod
	teid         uint32
	downlinkIP   string
	dstInterface uint8

	zeroBasedOuterHeader bool
	isActionSet          bool
	isInterfaceSet       bool
}

const (
	FarNoFuzz           = 0
	FarWithAction       = 1
	FarWithTEID         = 2
	FarWithDstInterface = 3
	FarMax              = 4
)

// NewFARBuilder returns a farBuilder.
func NewFARBuilder() *farBuilder {
	return &farBuilder{}
}

func (b *farBuilder) FuzzIE(ieType int, arg uint) *farBuilder {
	switch ieType {
	case FarWithAction:
		logger.PfcpsimLog.Infoln("FarWithAction")
		return b.WithAction(uint8(arg))
	case FarWithTEID:
		logger.PfcpsimLog.Infoln("FarWithTEID")
		return b.WithTEID(uint32(arg))
	case FarWithDstInterface:
		logger.PfcpsimLog.Infoln("FarWithDstInterface")
		return b.WithDstInterface(uint8(arg))
	default:
	}

	return b
}

func (b *farBuilder) WithID(id uint32) *farBuilder {
	b.farID = id
	return b
}

func (b *farBuilder) WithZeroBasedOuterHeaderCreation() *farBuilder {
	b.zeroBasedOuterHeader = true
	return b
}

func (b *farBuilder) WithMethod(method IEMethod) *farBuilder {
	b.method = method
	return b
}

func (b *farBuilder) WithAction(action uint8) *farBuilder {
	b.isActionSet = true
	b.applyAction = action

	return b
}

func (b *farBuilder) WithTEID(teid uint32) *farBuilder {
	b.teid = teid
	return b
}

func (b *farBuilder) WithDstInterface(iFace uint8) *farBuilder {
	b.isInterfaceSet = true
	b.dstInterface = iFace

	return b
}

func (b *farBuilder) WithDownlinkIP(downlinkIP string) *farBuilder {
	b.downlinkIP = downlinkIP
	return b
}

func (b *farBuilder) validate() {
	if b.farID == 0 {
		logger.PfcpsimLog.Panicln("tried building FAR without setting FAR ID")
	}

	if !b.isInterfaceSet {
		logger.PfcpsimLog.Panicln("tried building FAR without setting a destination interface")
	}

	if b.applyAction == ActionDrop|ActionForward {
		logger.PfcpsimLog.Panicln("tried building FAR with actions' ActionDrop and ActionForward flags")
	}

	if !b.isActionSet {
		logger.PfcpsimLog.Panicln("tried building FAR without setting an action")
	}
}

// BuildFAR returns a downlinkFAR if MarkAsDownlink was invoked.
// Returns an UplinkFAR if MarkAsUplink was invoked.
func (b *farBuilder) BuildFAR() *ie.IE {
	if doCheck {
		b.validate()
	}

	fwdParams := ie.NewForwardingParameters(
		ie.NewDestinationInterface(b.dstInterface),
	)

	createFunc := ie.NewCreateFAR
	if b.method == Update {
		createFunc = ie.NewUpdateFAR
		fwdParams = ie.NewUpdateForwardingParameters(
			ie.NewDestinationInterface(b.dstInterface),
		)
	}

	if b.zeroBasedOuterHeader {
		fwdParams.Add(ie.NewOuterHeaderCreation(S_TAG, 0, "0.0.0.0", "", 0, 0, 0))
	} else if b.downlinkIP != "" { // TODO revisit code and improve its structure
		// TEID and DownlinkIP are provided
		fwdParams.Add(ie.NewOuterHeaderCreation(S_TAG, b.teid, b.downlinkIP, "", 0, 0, 0))
	}

	far := createFunc(
		ie.NewFARID(b.farID),
		ie.NewApplyAction(b.applyAction),
		fwdParams,
	)

	if b.method == Delete {
		return ie.NewRemoveFAR(far)
	}

	return far
}
