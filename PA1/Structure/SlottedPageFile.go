package Structure

import (
	"fmt"
	"io"
	"os"
)

type SlottedPageFile struct {
	name            string
	file            *os.File
	slottedPageSize int
	seeks           int
	reads           int
	writes          int
}

func CreateSlottedPageFile(name string, slottedPageSize int) SlottedPageFile {
	file, err := os.Create(name)
	fmt.Println(err)
	return SlottedPageFile{name, file, slottedPageSize, 0, 0, 0}
}

func (spf SlottedPageFile) ToString() string {
	return fmt.Sprintf("{name: %s, reads: %q, writes: %q}", spf.name, spf.reads, spf.writes)
}

func (spf SlottedPageFile) Size() int {
	file, _ := spf.file.Stat()
	println()
	return (int)(file.Size()) / spf.slottedPageSize
}
func (spf SlottedPageFile) Close() {
	spf.file.Close()
}
func (spf SlottedPageFile) Clear() {
	spf.Clear()
}
func (spf SlottedPageFile) Get(pageID int) *SlottedPage {
	if pageID < 0 {
		return nil
	}
	pos := int64(pageID) * int64(spf.slottedPageSize)
	stats, _ := spf.file.Stat()
	if pos+int64(spf.slottedPageSize) > stats.Size() {
		return nil
	}
	spf.Seek(pos)
	sp := CreatePage(pageID, spf.slottedPageSize)
	spf.file.Read(sp.data)
	spf.reads += 1
	return &sp
}
func (spf SlottedPageFile) Save(sp SlottedPage) {
	spf.Seek(int64(sp.pageID * spf.slottedPageSize))
	spf.file.Write(sp.data)
	spf.writes += 1
}
func (spf SlottedPageFile) Seek(pos int64) {
	offset, _ := spf.file.Seek(0, io.SeekCurrent)
	if pos != offset {
		spf.file.Seek(0, int(pos))
	}
}
