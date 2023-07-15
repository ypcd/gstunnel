package gstunnellib

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"
)

type RunTimeDebugInfo interface {
	AddPackingPackSizeList(string, int)
	WriteFile(FileName string)
}

type runTimeDebugInfoImp struct {
	PackingPackSizeList map[string][]int
	//chanPackingPackSizeList	chan int
	lock sync.Mutex
}

type runTimeDebugInfoImpJson struct {
	PackingPackSizeLenList map[string]int
	PackingPackSizeList    map[string][]int
}

func newRunTimeDebugInfoImpJson(di *runTimeDebugInfoImp) *runTimeDebugInfoImpJson {
	dij := runTimeDebugInfoImpJson{
		PackingPackSizeLenList: make(map[string]int, 0),
		PackingPackSizeList:    di.PackingPackSizeList}

	for k, v := range di.PackingPackSizeList {
		dij.PackingPackSizeLenList[k] = len(v)
	}
	return &dij
}

func NewRunTimeDebugInfo() RunTimeDebugInfo {
	di := runTimeDebugInfoImp{}
	di.init()
	return &di
}

func (di *runTimeDebugInfoImp) init() {
	di.PackingPackSizeList = make(map[string][]int, 0)
}

func (di *runTimeDebugInfoImp) AddPackingPackSizeList(name string, size int) {
	di.lock.Lock()
	defer di.lock.Unlock()

	di.PackingPackSizeList[name] = append(di.PackingPackSizeList[name], size)
}

func (di *runTimeDebugInfoImp) WriteFile(FileName string) {
	di.lock.Lock()
	defer di.lock.Unlock()
	f, err := os.Create(FileName)
	defer func() {
		err := f.Close()
		CheckError_panic(err)
	}()
	CheckError_panic(err)
	data, err := json.Marshal(newRunTimeDebugInfoImpJson(di))
	CheckError_panic(err)
	_, err = io.Copy(f, bytes.NewBuffer(data))
	CheckError_panic(err)
}
