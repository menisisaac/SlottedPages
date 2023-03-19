package Structure

import "fmt"

type fileManager struct {
	id2file         map[int]SlottedPageFile
	slottedPageSize int
}

func CreateFileManager(slottedPageSize int) fileManager {
	return fileManager{map[int]SlottedPageFile{}, slottedPageSize}
}

func (fm fileManager) firstLocation() int64 {
	return 0
}

func (fm fileManager) shutdown() {
	for _, val := range fm.id2file {
		val.Close()
	}
}

func (fm fileManager) toString() string {
	rep := ""
	for _, val := range fm.id2file {
		rep += val.ToString()
	}
	return rep
}

func (fm fileManager) add(fileID int, data interface{}) int64 {
	size := fm.size(fileID)
	var p SlottedPage
	var location int64
	if size == 0 {
		p = CreatePage(fileID, fm.slottedPageSize)
		location = concatenate(p.pageID, p.Add(data))
	} else {
		p = *fm.page(fileID, size-1)
		dataLocation := p.Add(data)
		if dataLocation == -1 {
			p = CreatePage(p.pageID+1, fm.slottedPageSize)
			location = concatenate(p.pageID, p.Add(data))
		} else {
			location = concatenate(p.pageID, dataLocation)
		}
	}
	fm.updated(p, fileID)
	return location
}

func (fm fileManager) put() interface{} {
	return 0
}

func (fm fileManager) get(fileID int, location int64) interface{} {
	if first(location) < 0 || second(location) < 0 {
		return -1
	}
	p := fm.page(fileID, first(location))
	if p == nil {
		return -1
	} else {
		return p.Get(int(second(location)))
	}
}

func (fm fileManager) remove() interface{} {
	return 0
}

func (fm fileManager) clear() {

}

func (fm fileManager) iterator() {

}

func (fm fileManager) size(fileID int) int {
	file := fm.file(fileID)
	if file == nil {
		return 0
	}
	return file.Size()
}

func (fm fileManager) page(fileID, pageID int) *SlottedPage {
	file := fm.file(fileID)
	page := file.Get(pageID)
	return page
}

func (fm fileManager) updated(p SlottedPage, fileID int) {
	file := fm.file(fileID)
	file.Save(p)
}

func concatenate(i, j int) int64 {
	first := uint64(i)
	second := uint64(j)
	return int64(first<<32) | int64(second)
}

func first(i int64) int {
	return int(i >> 32)
}

func second(i int64) int32 {
	return int32(i)
}

func (fm fileManager) file(fileID int) *SlottedPageFile {
	file, check := fm.id2file[fileID]

	if check == false {
		file = CreateSlottedPageFile(fmt.Sprintf("%q.dat", fileID), fm.slottedPageSize)
		fm.id2file[fileID] = file
	}
	return &file
}
