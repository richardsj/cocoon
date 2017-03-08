// Code generated by go-bindata.
// sources:
// data/bindata.go
// data/cocoon.job.json
// DO NOT EDIT!

package data

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _dataBindataGo = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func dataBindataGoBytes() ([]byte, error) {
	return bindataRead(
		_dataBindataGo,
		"data/bindata.go",
	)
}

func dataBindataGo() (*asset, error) {
	bytes, err := dataBindataGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/bindata.go", size: 0, mode: os.FileMode(420), modTime: time.Unix(1488979752, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dataCocoonJobJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x56\x4f\x8f\xa3\x36\x14\xbf\xcf\xa7\xb0\xac\x3d\x76\x93\x99\x56\x7b\x89\xd4\x03\x09\xe9\x88\x16\x12\x94\x90\x43\x55\x45\xc8\x31\x1e\xc6\x8a\xb1\x91\x6d\xb2\x9b\x22\xbe\x7b\x65\x43\x76\x00\x93\x4c\x73\x19\xcf\x7b\xbf\xf7\xc7\xbf\xf7\xc7\xd4\x4f\x00\xc0\x3f\xc5\x09\x2e\x80\x39\x02\x00\x77\x24\xa7\x82\xc3\x05\x80\x39\x13\x27\xc4\xe0\x2f\xad\x3c\xf0\x8d\xac\xae\x03\xbf\x69\x6e\xb2\x0d\x2a\x88\x2b\x4d\xae\xa5\x95\x2a\x22\x2f\x14\x93\x9b\x38\x96\x54\x48\xaa\xaf\x70\x01\xbe\x3d\x77\x32\x8f\x31\x4f\x6f\x39\x36\xf8\x37\xc4\x14\xe9\xe4\x3e\xd2\x08\x13\xae\x89\x54\x70\x01\xfe\xb1\x42\x00\x60\x86\x5f\xa0\x3d\x1f\x3b\xdc\x4a\x70\xa5\x25\xa2\x5c\xf7\x71\x75\xf7\x17\x00\x18\x26\x48\xe6\x44\x9b\x74\xbe\xd4\x48\x6b\x39\x3b\x13\xc9\x09\x9b\x71\x54\x90\x5b\xc6\xed\xbd\x3f\x90\x8c\xf2\xea\x47\x5f\xb7\x2d\x89\x44\x3c\x33\xba\xdf\x61\x27\x6e\x06\x89\x24\x48\x9d\x5f\xa5\xa8\xca\x3b\x79\x4c\x32\xd5\x5d\xa1\xe2\x26\x6a\x5d\xdb\x53\xd3\x0c\x74\xfd\xeb\xf1\x8a\xb1\x9e\xd2\x84\xec\x47\x1b\x46\x1c\x44\xd5\x48\x9d\xbf\x3a\xa1\x5b\xa6\x25\xbd\x10\x69\x40\x99\xc0\x67\x22\xc7\xfa\x95\xe0\x6f\x34\xff\xd9\x1e\x3d\x0d\x2d\x50\x7e\xbb\x93\x39\x3a\xbe\x01\x80\x6f\x42\x62\x92\x96\x15\x63\x70\x01\xb4\xac\x88\x83\xc0\xa2\x28\x3a\x6a\x4f\x48\xbd\xbb\x2e\x38\xd1\xdf\x85\x3c\xa7\x85\xc8\x6c\xb4\x77\xa1\xb4\x8b\x42\x32\xb7\x5c\xc0\x2f\xf5\x66\x1b\x79\x7e\x1a\xad\x13\x2f\xdd\xaf\x76\x41\x9c\xec\x53\x3f\xd8\x35\xf3\x81\xc6\x5f\xc7\xe1\xf6\xef\x0e\x90\x6e\xbc\x68\xdd\xc0\xa3\xe3\xb5\x94\xf4\x42\x19\xc9\x49\x76\x2f\xfd\x52\x48\x9d\x16\xa8\x34\xb1\x8f\x03\x65\x33\x22\x72\xcd\x2f\x53\x2c\xae\xb6\xab\xed\x76\x93\x4e\x4c\x97\x0b\x5a\x6d\xfd\x75\x7a\xd8\x85\x2d\x74\x25\xb0\x10\x7c\x25\x32\x72\xd8\x85\x9f\x59\x25\xde\xeb\xd8\x2a\x41\xf9\x67\x56\xa1\xb7\x71\xcc\x42\xc4\x1f\xd9\x2d\x0f\x41\xe8\xa7\xb1\xb7\xf3\xa2\x7d\xdf\x74\x59\x51\x96\xc5\x48\xa2\x42\x3d\xb0\xf6\x83\xfd\x5f\x69\x18\x44\x41\xd2\xda\xfa\x54\x9d\xa3\xa5\x63\x30\xe6\x76\xdf\x6e\x1a\xdb\x01\x0e\xc3\x77\x47\xaf\xd3\x27\xa8\x6b\x9d\x0e\xe0\xb6\x41\x2c\xa4\x0e\xd1\x89\x98\x26\x86\x58\x70\x4e\xb0\x16\x12\x02\x00\x86\x59\x8d\x2c\x61\x44\x34\x9a\xaa\xb9\xdb\x7c\xc6\xb1\xac\xf8\xd7\x9f\xce\x67\x53\xa3\xd0\x6b\x67\x63\x30\x67\x02\x23\x36\x57\x58\xd2\x52\x2b\xf8\x90\xa1\x50\xe4\xf7\x27\x39\x42\x3f\xfe\xa0\xcc\xd2\xf7\xf2\xec\x44\xed\xb4\x7b\xfa\x2f\x89\x96\x16\xf2\x30\x52\x42\x8a\x92\x21\xdd\x16\x63\xcc\x88\x27\x35\x7d\x43\x58\x4f\x56\x0a\x00\xf8\x4a\xb4\x26\x72\x2f\x2a\x89\xdb\x61\xd7\xba\x54\x8b\xf9\x5c\xa2\xef\xb3\x9c\xea\xf7\xea\x54\x29\x22\xb1\xe0\x9a\x70\x3d\xc3\xa2\x98\x73\x2c\x32\xa2\xe6\xd8\x36\xd9\xbc\x40\x4a\x13\x79\xa3\xe4\xf3\x91\x1f\x5f\xd6\xbe\x7f\x0c\x69\x7a\x21\x3e\x51\xfa\x53\x96\xc7\xf7\xdb\x11\x65\x73\x57\x93\xa3\x1e\x1f\xda\x3d\x1f\x1f\x9a\xc6\xa5\x99\x14\x42\x5e\x2d\xc3\x75\x7d\xfb\x67\x02\x17\x6c\x63\x33\x58\x6e\x9d\x36\xed\xa2\x1c\x3f\x09\xed\xcf\xe5\xda\x84\x5c\x52\x5b\x89\xba\xb6\x27\x37\x98\x45\xf9\x57\x8e\x0a\x8a\xcd\x10\x28\xb8\x70\x5d\x5b\xf7\x00\x76\x03\xd2\x9b\x8f\x66\x02\x7a\x74\x64\x63\xd4\xe3\x2d\xea\x53\x55\x22\x8d\xdf\x63\x74\x65\x02\x99\xb5\x5c\xf7\x1d\x7c\x9c\x8f\xfd\xa7\xfd\x4e\x55\xee\x57\xe4\xff\x54\x63\xb2\x12\x83\x2a\x7c\x5c\xa5\x19\x66\xa3\x91\xd4\xb1\x60\x14\x5f\xc7\x19\x05\xe6\x9b\xe7\x82\xcc\xa6\xf9\xed\xf9\xe3\x37\x08\xe1\x69\x4d\x8a\x52\xbb\xe3\x0a\x7d\xc2\x90\x71\xf9\xeb\xb7\x69\xcb\xa8\x7b\x43\x33\x8b\x9b\xcc\xee\xb6\xb3\x9a\xc9\xcf\x9c\x43\x99\x21\x4d\x7a\x49\xc3\xbd\x46\x79\x6e\x3f\x20\x5e\xdc\x74\xcd\xee\x30\x3b\x9f\x31\xbb\x3b\x5f\x9e\x6e\x1e\x9b\xa7\xe6\xbf\x00\x00\x00\xff\xff\xa9\x48\x3b\x87\x7a\x0a\x00\x00")

func dataCocoonJobJsonBytes() ([]byte, error) {
	return bindataRead(
		_dataCocoonJobJson,
		"data/cocoon.job.json",
	)
}

func dataCocoonJobJson() (*asset, error) {
	bytes, err := dataCocoonJobJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/cocoon.job.json", size: 2682, mode: os.FileMode(420), modTime: time.Unix(1488898108, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"data/bindata.go": dataBindataGo,
	"data/cocoon.job.json": dataCocoonJobJson,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"data": &bintree{nil, map[string]*bintree{
		"bindata.go": &bintree{dataBindataGo, map[string]*bintree{}},
		"cocoon.job.json": &bintree{dataCocoonJobJson, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

