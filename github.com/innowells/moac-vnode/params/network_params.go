// Copyright 2017 The MOAC-core Authors
// This file is part of the MOAC-core library.
//
// The MOAC-core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The MOAC-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the MOAC-core library. If not, see <http://www.gnu.org/licenses/>.

package params

import "time"

// These are network parameters that need to be constant between clients, but
// aren't necesarilly consensus related.

const (
	// BloomBitsBlocks is the number of blocks a single bloom bit section vector
	// contains.
	BloomBitsBlocks uint64 = 4096

	TimerPingInterval time.Duration = 1

	DirectCallLimitPerBlock = 2048
	DirectCallGasLimit = 4000000
	
	SubchainMsgLimit = 10000
	ScsMsgLimit = 1000

	DirectCall   = 1
	BroadCast    = 2
	ControlMsg   = 3
	ScsShakeHand = 4
	ScsPing      = 5

	// None            = -1
	RegOpen         = 0
	RegClose        = 1
	CreateProposal  = 2
	DisputeProposal = 3
	ApproveProposal = 4
	RegAdd          = 5
	RegAsMonitor    = 6
	RegAsBackup     = 7
)

type ScsKind int64
const (
	None ScsKind = iota
	ConsensusScs
	MonitorScs
	BackupScs
	MatchSelTarget
)
