package ssz

import "encoding/binary"

type ReaderSource struct {
	payload []byte
}

type ReadPos struct {
	Offset       uint32
	NextOffset   uint32
	ContainerEnd uint32
}

func (p ReadPos) Add(localOffset uint32) ReadPos {
	return ReadPos{
		Offset:       p.Offset + localOffset,
		NextOffset:   0,
		ContainerEnd: p.ContainerEnd,
	}
}

func (p ReadPos) AddWithNext(localOffset, localNextOffset uint32) ReadPos {
	return ReadPos{
		Offset:       p.Offset + localOffset,
		NextOffset:   p.Offset + localNextOffset,
		ContainerEnd: p.ContainerEnd,
	}
}

type Reader[T any] interface {
	InitReaderSSZ(pos ReadPos) T
}

type StaticReader interface {
	SizeSSZ() uint32
}

type DynamicReader interface {
	SizeSSZ(fixed bool) uint32
}

func (r *ReaderSource) offset(o uint32) uint32 {
	return binary.LittleEndian.Uint32(r.payload[o : o+4])
}

// objectEnd returns the end position of the object referred to by pos.
func (r *ReaderSource) objectEnd(pos ReadPos) uint32 {
	if pos.NextOffset == 0 {
		return pos.ContainerEnd
	}
	return r.offset(pos.NextOffset)
}

type Uint64Reader struct {
	pos ReadPos
}

func (r Uint64Reader) InitReaderSSZ(pos ReadPos) Uint64Reader {
	return Uint64Reader{pos: pos}
}

func (r Uint64Reader) Read(src *ReaderSource, v *uint64) {
	data := src.payload[r.pos.Offset : int(r.pos.Offset)+8]
	*v = binary.LittleEndian.Uint64(data)
}

type ByteArrayReader struct {
	pos  ReadPos
	Size int
}

func (r ByteArrayReader) InitReaderSSZ(pos ReadPos) ByteArrayReader {
	return ByteArrayReader{pos: pos, Size: r.Size}
}

func (r ByteArrayReader) Read(src *ReaderSource, v []byte) {
	if len(v) != r.Size {
		panic("invalid size")
	}
	copy(v, src.payload[r.pos.Offset:int(r.pos.Offset)+r.Size])
}

type DynamicBytesReader struct {
	pos ReadPos
}

func (r DynamicBytesReader) InitReaderSSZ(pos ReadPos) DynamicBytesReader {
	return DynamicBytesReader{pos: pos}
}

func (r DynamicBytesReader) Read(src *ReaderSource, v []byte) {
	start := src.offset(r.pos.Offset)
	end := src.objectEnd(r.pos)
	copy(v, src.payload[start:end])
}

type ListReader[Item Reader[Item]] struct {
	Prototype Item
	pos       ReadPos
}

func (r ListReader[Item]) InitReaderSSZ(pos ReadPos) ListReader[Item] {
	return ListReader[Item]{
		Prototype: r.Prototype,
		pos:       pos,
	}
}

func (r ListReader[Item]) Item(src *ReaderSource, n int) Item {
	// Get the size of the container.
	start := src.offset(r.pos.Offset)
	end := src.objectEnd(r.pos)
	if start == end {
		panic("out of bounds") // Len() == 0
	}

	var itemSize uint32
	var nextOffset uint32
	switch vs := any(r.Prototype).(type) {
	case StaticReader:
		itemSize = vs.SizeSSZ()
		length := int((end - start) / itemSize)
		if n >= length {
			panic("out of bounds")
		}
	case DynamicReader:
		itemSize = vs.SizeSSZ(false)
		length := int((end - start) / itemSize)
		if n >= length {
			panic("out of bounds")
		}
		if n <= length-1 {
			nextOffset = start + 4
		}
	}

	pos := ReadPos{Offset: start + uint32(n), NextOffset: nextOffset, ContainerEnd: end}
	return r.Prototype.InitReaderSSZ(pos)
}

func (r ListReader[Item]) Len(src *ReaderSource) int {
	start := src.offset(r.pos.Offset)
	end := src.objectEnd(r.pos)
	return int((end - start) / r.ItemSize())
}

func (r ListReader[Item]) ItemSize() uint32 {
	switch vs := any(r.Prototype).(type) {
	case StaticReader:
		return vs.SizeSSZ()
	case DynamicReader:
		return vs.SizeSSZ(false)
	default:
		panic("invalid item type")
	}
}
