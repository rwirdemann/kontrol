// Code generated by go-bindata.
// sources:
// html/assets.go
// html/index.html
// DO NOT EDIT!

package html

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

var _htmlAssetsGo = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func htmlAssetsGoBytes() ([]byte, error) {
	return bindataRead(
		_htmlAssetsGo,
		"html/assets.go",
	)
}

func htmlAssetsGo() (*asset, error) {
	bytes, err := htmlAssetsGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "html/assets.go", size: 0, mode: os.FileMode(436), modTime: time.Unix(1511341574, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _htmlIndexHtml = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x74\x91\x4d\x6b\xf3\x30\x10\x84\xef\xfe\x15\x8b\xcf\x2f\x52\xf2\x1e\x85\x30\xa4\x1f\xb4\x25\x25\x2d\x4d\x6e\xa5\x07\x57\xde\xc6\x26\x95\x56\x95\xd6\x2d\xc1\xe8\xbf\x17\xe7\x0b\xec\x26\x73\x30\x78\xe6\x99\x91\x40\xba\x9e\x16\x73\xb2\xd6\x36\x6c\xd1\x31\xcc\xc9\x71\xa0\x4f\x58\x62\xf8\x6e\x0c\x6a\x59\x4f\x8b\x4c\xfb\x80\x45\x06\x00\xa0\x23\x07\x72\xeb\xe2\x05\xbf\x5a\x8c\xac\xb4\x3c\x18\xbb\xf4\xee\x76\x05\x35\xb3\x57\x52\x76\x9d\xb8\xa7\xc8\xae\xb4\x98\x92\xea\x3a\xf1\x4c\x81\x53\x92\x9b\xfd\xbe\x2c\x8d\xa1\xd6\x71\xcc\x06\xb3\x37\x18\x4d\x68\x3c\x37\xe4\x46\xd3\x8f\x4d\x64\xa0\x0f\x38\xf6\x44\x36\xba\x4f\xf4\xe4\x22\xc2\x35\x55\x18\x47\xdd\xff\x93\x89\x82\x65\x6b\x0c\xc6\x78\xa1\x76\x45\xd5\x76\xd4\xea\x76\xdf\x5e\xf9\xec\x70\x68\xae\xe0\xf5\xe4\x0e\x99\xa3\x9e\x7e\x1c\x06\x75\x26\xe8\xf5\x50\x29\xc8\x67\x8b\xfc\xdf\xd9\x74\x51\x5a\xec\x73\xb7\x41\x58\x60\x1d\xd0\xbd\x63\x58\x5f\x80\x57\x5b\xdf\xc3\xbe\x0c\xec\x30\xe4\x7f\x98\x34\x70\xd2\x70\x44\x08\x71\xfa\x7f\xcb\xf6\xb8\x96\xfd\x23\xff\x06\x00\x00\xff\xff\x0a\xe9\xa8\xe7\x0f\x02\x00\x00")

func htmlIndexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_htmlIndexHtml,
		"html/index.html",
	)
}

func htmlIndexHtml() (*asset, error) {
	bytes, err := htmlIndexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "html/index.html", size: 527, mode: os.FileMode(436), modTime: time.Unix(1510864199, 0)}
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
	"html/assets.go": htmlAssetsGo,
	"html/index.html": htmlIndexHtml,
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
	"html": &bintree{nil, map[string]*bintree{
		"assets.go": &bintree{htmlAssetsGo, map[string]*bintree{}},
		"index.html": &bintree{htmlIndexHtml, map[string]*bintree{}},
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

