package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type slottedPage struct {
	pageID int
	data   []byte
}
type headerSlot struct {
	Location int64
	Size     int64
}

func createPage(pageID, size int) slottedPage {
	data := make([]byte, size)
	sp := slottedPage{size, data}
	sp.writeInt(0, 0)
	sp.setStartOfDataStorage(size - 4)
	return sp
}
func (sp slottedPage) add(data interface{}) int {
	location, size := sp.saveObject(data)
	index := sp.entryCount()
	sp.setEntryCount(index + 1)
	sp.saveLocation(index, location, size)
	return index
}

func (sp slottedPage) get(index int) interface{} {
	if index < 0 || index > sp.entryCount()-1 {
		return nil
	}
	location := sp.getLocation(index)
	if location.Location < 0 {
		return nil
	}
	return toData(sp.data, int(location.Location), int(location.Size))
}

func (sp slottedPage) writeInt(location, value int) {
	sp.data[location] = (byte)(value >> 24)
	sp.data[location+1] = (byte)(value >> 16)
	sp.data[location+2] = (byte)(value >> 8)
	sp.data[location+3] = (byte)(value)
}

func (sp slottedPage) entryCount() int {
	return sp.readInt(0)
}

// Golang uses logical shift on uints, arithmetic shift on ints
func (sp slottedPage) readInt(location int) int {
	return int(((uint)(sp.data[location]) << 24) + (((uint)(sp.data[location+1]) & 0xFF) << 16) + (((uint)(sp.data[location+2]) & 0xFF) << 8) + ((uint)(sp.data[location+3]) & 0xFF))
}

func (sp slottedPage) setStartOfDataStorage(startOfDataStorage int) {
	sp.writeInt(cap(sp.data)-4, startOfDataStorage)
}

func (sp slottedPage) startOfDataStorage() int {
	return sp.readInt(cap(sp.data) - 4)
}

func (sp slottedPage) freeSpaceSize() int {
	return sp.startOfDataStorage() - sp.headerSize()
}

func (sp slottedPage) headerSize() int {
	return 4 * (sp.entryCount() + 1)
}

func (sp slottedPage) setEntryCount(count int) {
	sp.writeInt(0, count)
}

func (sp slottedPage) getLocation(index int) headerSlot {
	location := headerSlot{}
	binary.Read(bytes.NewBuffer(sp.data[(index+1)*4:]), binary.BigEndian, &location)
	return location
}

func (sp slottedPage) saveLocation(index, location, size int) {
	loca := &headerSlot{int64(location), int64(size)}
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, loca)
	copy(sp.data[(index+1)*4:], buf.Bytes())
}
func (sp slottedPage) compact() {
	return
}

func (sp slottedPage) saveObject(data interface{}) (int, int) {
	arr, size := toByteArray(data)
	return sp.save(arr), size
}

func (sp slottedPage) save(b []byte) int {
	if sp.freeSpaceSize() < cap(b)+4 {
		sp.compact()
		if sp.freeSpaceSize() < cap(b)+4 {
			fmt.Println("Error: Not Enough Capacity !!!")
			return -1
		}
	}
	location := sp.startOfDataStorage() - cap(b)
	copy(sp.data[location:], b[:])
	sp.setStartOfDataStorage(location)
	return location
}

func toByteArray(data interface{}) ([]byte, int) {
	dataObject := fmt.Sprint(data)
	arr := []byte(dataObject)
	return arr, len(arr)
}

func toData(b []byte, offset, size int) string {
	println(offset, "    ", size)
	return fmt.Sprintf("%s", b[offset:offset+size])
}

func test(data interface{}) string {
	dataObject := fmt.Sprint(data)
	return dataObject
}
