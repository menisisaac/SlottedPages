package Structure

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type SlottedPage struct {
	pageID int
	data   []byte
}
type HeaderSlot struct {
	Location int64
	Size     int64
}

func CreatePage(pageID, size int) SlottedPage {
	data := make([]byte, size)
	sp := SlottedPage{size, data}
	sp.WriteInt(0, 0)
	sp.SetStartOfDataStorage(size - 4)
	return sp
}
func (sp SlottedPage) Add(data interface{}) int {
	location, size := sp.saveObject(data)
	if location == -1 {
		return -1
	}
	index := sp.EntryCount()
	sp.SetEntryCount(index + 1)
	sp.SaveLocation(index, location, size)
	return index
}

func (sp SlottedPage) Get(index int) interface{} {
	if index < 0 || index > sp.EntryCount()-1 {
		return nil
	}
	location := sp.GetLocation(index)
	if location.Location < 0 || location.Location > int64(len(sp.data)) {
		return nil
	}
	return toData(sp.data, int(location.Location), int(location.Size))
}
func (sp SlottedPage) Remove(index int) interface{} {
	if index < 0 || index >= sp.EntryCount() {
		println("Add error here")
		return nil
	}
	location := sp.GetLocation(index)
	if location.Location < 0 {
		return nil
	}
	data := toData(sp.data, int(location.Location), int(location.Size))
	sp.SaveLocation(index, -1, -1)
	return data

}

func (sp SlottedPage) WriteInt(location, value int) {
	sp.data[location] = (byte)(value >> 24)
	sp.data[location+1] = (byte)(value >> 16)
	sp.data[location+2] = (byte)(value >> 8)
	sp.data[location+3] = (byte)(value)
}

func (sp SlottedPage) EntryCount() int {
	return sp.ReadInt(0)
}

// Golang uses logical shift on uints, arithmetic shift on ints
func (sp SlottedPage) ReadInt(location int) int {
	return int(((uint)(sp.data[location]) << 24) + (((uint)(sp.data[location+1]) & 0xFF) << 16) + (((uint)(sp.data[location+2]) & 0xFF) << 8) + ((uint)(sp.data[location+3]) & 0xFF))
}

func (sp SlottedPage) SetStartOfDataStorage(startOfDataStorage int) {
	sp.WriteInt(cap(sp.data)-4, startOfDataStorage)
}

func (sp SlottedPage) StartOfDataStorage() int {
	return sp.ReadInt(cap(sp.data) - 4)
}

func (sp SlottedPage) FreeSpaceSize() int {
	return sp.StartOfDataStorage() - sp.headerSize()
}

func (sp SlottedPage) headerSize() int {
	var temp HeaderSlot
	return binary.Size(temp)*(sp.EntryCount()) + 4
}

func (sp SlottedPage) SetEntryCount(count int) {
	sp.WriteInt(0, count)
}

func (sp SlottedPage) GetLocation(index int) HeaderSlot {
	location := HeaderSlot{}
	binary.Read(bytes.NewBuffer(sp.data[(index+1)*binary.Size(location)+4:]), binary.BigEndian, &location)
	return location
}

func (sp SlottedPage) SaveLocation(index, location, size int) {
	loca := &HeaderSlot{int64(location), int64(size)}
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, loca)
	copy(sp.data[(index+1)*binary.Size(loca)+4:], buf.Bytes())
}
func (sp SlottedPage) Compact() {
	sp.SetStartOfDataStorage(len(sp.data) - 4)
	for i := 0; i < sp.EntryCount(); i += 1 {
		location := sp.GetLocation(i)
		if location.Location >= 0 {
			newLocation := sp.StartOfDataStorage() - int(location.Size)
			copy(sp.data[newLocation:], sp.data[location.Location:location.Location+location.Size])
			sp.SetStartOfDataStorage(newLocation)
			sp.SaveLocation(i, newLocation, int(location.Size))
		}
	}
}

func (sp SlottedPage) saveObject(data interface{}) (int, int) {
	arr, size := toByteArray(data)
	return sp.save(arr), size
}

func (sp SlottedPage) save(b []byte) int {
	if sp.FreeSpaceSize() < cap(b)+4 {
		sp.Compact()
		if sp.FreeSpaceSize() < cap(b)+4 {
			fmt.Println("Error: Not Enough Capacity !!!")
			return -1
		}
	}
	location := sp.StartOfDataStorage() - cap(b)
	copy(sp.data[location:], b[:])
	sp.SetStartOfDataStorage(location)
	return location
}

func Iterator(sp SlottedPage) func() interface{} {
	n := 0
	page := sp
	return func() interface{} {
		var currdata interface{} = nil
		n = hasNext(n, page)
		if n > -1 {
			currdata = page.Get(n)
			n += 1
		}
		return currdata
	}
}
func hasNext(i int, sp SlottedPage) int {
	if i > sp.EntryCount()-1 {
		return -1
	} else if sp.GetLocation(i).Location != -1 {
		return i
	} else {
		return hasNext(i+1, sp)
	}
}

func toByteArray(data interface{}) ([]byte, int) {
	dataObject := fmt.Sprint(data)
	arr := []byte(dataObject)
	return arr, len(arr)
}

func toData(b []byte, offset, size int) string {
	return fmt.Sprintf("%s", b[offset:offset+size])
}

func test(data interface{}) string {
	dataObject := fmt.Sprint(data)
	return dataObject
}
