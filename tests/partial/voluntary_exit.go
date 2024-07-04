// ssz: Go Simple Serialize (SSZ) codec library
// Copyright 2024 ssz Authors
// SPDX-License-Identifier: BSD-3-Clause

package partialtests

import "github.com/karalabe/ssz"

type VoluntaryExit struct {
	Epoch          uint64
	ValidatorIndex uint64
}

func (v *VoluntaryExit) SizeSSZ() uint32 { return 16 }

func (v *VoluntaryExit) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineUint64(codec, &v.Epoch)          // Field (0) - Epoch          - 8 bytes
	ssz.DefineUint64(codec, &v.ValidatorIndex) // Field (1) - ValidatorIndex - 8 bytes
}

type VoluntaryExitReader struct{
	pos ssz.ReadPos
}

func (v VoluntaryExitReader) SizeSSZ() uint32 {
	return 16
}

func (v VoluntaryExitReader) InitReaderSSZ(pos ssz.ReadPos) VoluntaryExitReader {
	return VoluntaryExitReader{pos}
}

func (v VoluntaryExitReader) Epoch() ssz.Uint64Reader {
	return ssz.Uint64Reader{}.InitReaderSSZ(v.pos.Add(0))
}

func (v VoluntaryExitReader) ValidatorIndex() ssz.Uint64Reader {
	return ssz.Uint64Reader{}.InitReaderSSZ(v.pos.Add(8))
}

type SignedVoluntaryExit struct {
	Exit      *VoluntaryExit `json:"message"`
	Signature [96]byte       `json:"signature" ssz-size:"96"`
}

func (v *SignedVoluntaryExit) SizeSSZ() uint32 { return 112 }

func (v *SignedVoluntaryExit) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticObject(codec, &v.Exit)       // Field (0) - Exit      - 16 bytes
	ssz.DefineStaticBytes(codec, v.Signature[:]) // Field (1) - Signature - 96 bytes
}

type SignedVoluntaryExitReader struct {
	pos ssz.ReadPos
}

func (v SignedVoluntaryExitReader) SizeSSZ() uint32 {
	return 112
}

func (v SignedVoluntaryExitReader) InitReaderSSZ(pos ssz.ReadPos) SignedVoluntaryExitReader {
	return SignedVoluntaryExitReader{pos}
}

func (v SignedVoluntaryExitReader) Exit() VoluntaryExitReader {
	return VoluntaryExitReader{}.InitReaderSSZ(v.pos.Add(0))
}

func (v SignedVoluntaryExitReader) Signature() ssz.ByteArrayReader {
	return ssz.ByteArrayReader{Size: 96}.InitReaderSSZ(v.pos.Add(16))
}


