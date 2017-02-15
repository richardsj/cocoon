// Code generated by go-bindata.
// sources:
// bindata.go
// cocoon.job.json
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

var _bindataGo = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func bindataGoBytes() ([]byte, error) {
	return bindataRead(
		_bindataGo,
		"bindata.go",
	)
}

func bindataGo() (*asset, error) {
	bytes, err := bindataGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "bindata.go", size: 0, mode: os.FileMode(420), modTime: time.Unix(1487174149, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _cocoonJobJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x55\x4d\x8f\xda\x30\x10\xbd\xef\xaf\xb0\xac\x1e\xcb\x57\xdb\xbd\x20\xf5\x40\x61\xbb\xa2\xe2\x4b\x7c\x9c\x2a\xb4\x9a\x4d\x66\x83\x85\x63\x47\xf6\x40\x97\x46\xf9\xef\x95\x9d\xb0\x1b\x92\xc0\x96\x0b\xce\xbc\x37\x9e\xf7\x26\x63\x27\xbd\x63\x8c\xff\xd2\xcf\xbc\xcf\xdc\x92\x31\xbe\xc4\x48\x68\xc5\xfb\x8c\x47\x52\x3f\x83\xe4\x9f\xf3\xf8\x78\xe4\x62\x81\x0e\xb4\x56\xad\x34\x1d\x8f\xb2\xec\x0c\xcd\x20\xc6\xab\xe0\xfa\x94\x78\xd0\xa2\x39\x8a\x00\xcf\xe1\x85\x11\xda\x08\x3a\xf1\x3e\xbb\xef\x16\xb1\x81\x94\x03\x9a\xab\xc0\xf1\x5f\x40\x5a\x2c\xe2\x23\x20\x08\x50\x11\x1a\xcb\xfb\xec\xb7\x0f\x32\xc6\xc3\xa0\xc7\xfd\x7a\x5b\xf0\x86\x5a\x59\x32\x20\x14\x95\x79\x69\xf1\xcf\x18\x9f\xac\xc1\x44\x48\x4e\xce\xa7\x14\x88\x4c\x7b\x8f\x46\xa1\x6c\x2b\x88\xf1\xac\x38\xef\xc2\x3b\x53\x0a\x75\x78\x2d\x63\xf3\x04\x0d\xa8\xd0\x61\xdf\x79\x11\xce\x2e\x84\xac\xc1\xee\x1f\x8d\x3e\x24\x57\x74\xdc\x6a\x58\xe1\xe4\xa0\x5c\xf1\x34\xf5\xab\x2c\xbb\xc0\xca\x2e\xd5\x41\xca\x12\xe8\x2a\x97\x8b\x5e\x16\x6e\x2a\x4e\x60\xf7\x75\x05\x79\xdf\x8d\x38\xa2\x71\xdc\x50\x07\x7b\x34\x55\x7c\xa8\xd5\x8b\x88\xde\x46\xa7\x84\x88\x18\x22\x5f\x24\x4d\xc7\x6e\x59\xdb\x9b\xb9\xfa\x71\x5c\xb4\xf1\x19\xec\xae\x4e\x00\x13\x79\x2b\xbc\x23\x75\x00\xb2\x63\x03\x23\x12\xb2\x9d\x10\x8f\xad\x10\x13\xa9\x4f\xad\x40\x2b\x85\x01\x69\xd3\xea\xb6\xbb\xed\x6f\x6d\xbb\xe3\xdb\x8b\x6d\xb2\x8a\xe4\x07\x75\x6c\xd2\x3b\x9c\x0f\xe7\xf3\xd9\xd3\xf5\x19\xaf\x73\x87\xf3\xd1\xc3\xd3\x66\x39\xc9\x5d\x0e\x7d\xce\x50\x87\xb8\x59\x4e\x3e\xca\x5a\x0f\x1e\xab\x59\x6b\x88\x3e\xca\x9a\x0c\x66\xb5\xb4\x09\x28\x97\x77\xd3\xf2\x2a\x3f\x78\xbe\x95\xdb\x0a\x36\x45\x02\xd7\x8f\x6a\xce\x44\x47\xd7\x5f\xee\x14\x5e\x7f\x0a\xe9\x37\xec\x75\x6b\x92\x0b\x74\x25\xfe\xe2\xf4\x87\xa7\xdc\x54\xb7\xc6\x38\x91\x40\xcd\xf2\x06\x86\xc4\x0b\x04\xf9\x71\xae\xea\x60\x8c\x3f\x22\x11\x9a\x95\x3e\x18\x7f\x67\xf0\x1d\x51\x62\xfb\x9d\x8e\x25\x6d\x20\xc2\x76\xa4\x75\x24\x11\x12\x61\xdb\x81\x8e\x3b\x7b\xa3\x23\x50\x1f\xcc\x4f\xd5\x90\xbf\x12\x25\x90\x38\xe2\x08\xad\xbf\x11\x2e\xe7\xb1\xd2\xfd\xaa\x87\x25\x5a\xaf\xcf\x36\xce\xdd\x62\x93\x9f\xf2\xc5\x26\xcb\xea\xad\xc4\x58\x9b\x93\xef\x62\x9a\x9e\x1f\x1a\x78\xe3\xf9\x62\xc5\xfb\xac\xfe\x2e\x66\x48\x7f\xb4\xc9\x2f\x84\xdb\xe7\x62\x24\x6c\x02\x14\xec\x16\x70\x92\x1a\x42\x3f\x13\x25\xc6\xfb\x7a\x5b\xbe\x24\xaf\x58\xbb\x6e\xeb\x7f\x2c\x39\x2d\xfb\x82\x91\x2f\x2b\x78\x93\xdd\x2b\x56\xb3\x4b\xb5\x04\x86\x16\x5a\x8a\xe0\x54\x55\x3c\x76\x5f\x97\x23\x48\xde\x67\x5f\xbb\xef\xbf\x8b\x12\x03\x22\x8c\x13\xaa\xcf\x3d\x1f\xa1\x04\xb7\xe5\x97\xfb\xe6\xcc\xa9\x0e\xfd\x7c\x86\x9e\xd7\xa8\xee\xed\x24\x36\x7e\x50\x36\x49\x08\x84\x25\xd1\x7c\x45\x10\x45\xfe\x72\xee\xd5\xe5\xba\x43\xb8\x00\x03\x52\xa2\x73\xd4\xbb\x3b\xef\x98\xdd\x65\xff\x02\x00\x00\xff\xff\x7a\x78\x4f\x58\xf2\x07\x00\x00")

func cocoonJobJsonBytes() ([]byte, error) {
	return bindataRead(
		_cocoonJobJson,
		"cocoon.job.json",
	)
}

func cocoonJobJson() (*asset, error) {
	bytes, err := cocoonJobJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cocoon.job.json", size: 2034, mode: os.FileMode(420), modTime: time.Unix(1487174068, 0)}
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
	"bindata.go": bindataGo,
	"cocoon.job.json": cocoonJobJson,
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
	"bindata.go": &bintree{bindataGo, map[string]*bintree{}},
	"cocoon.job.json": &bintree{cocoonJobJson, map[string]*bintree{}},
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

