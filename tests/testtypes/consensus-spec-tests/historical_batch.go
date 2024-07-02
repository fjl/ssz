// ssz: Go Simple Serialize (SSZ) codec library
// Copyright 2024 ssz Authors
// SPDX-License-Identifier: BSD-3-Clause

package consensus_spec_tests

import (
	"github.com/karalabe/ssz"
)

type HistoricalBatch struct {
	BlockRoots [8192]Hash
	StateRoots [8192]Hash
}

func (h *HistoricalBatch) StaticSSZ() bool { return true }
func (h *HistoricalBatch) SizeSSZ() uint32 { return 2 * 8192 * 32 }
func (h *HistoricalBatch) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineArrayOfStaticBytes(codec, h.BlockRoots[:])
	ssz.DefineArrayOfStaticBytes(codec, h.StateRoots[:])
}