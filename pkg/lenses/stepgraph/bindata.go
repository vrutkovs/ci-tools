// Code generated by go-bindata. (@generated) DO NOT EDIT.

//Package stepgraph generated by go-bindata.// sources:
// static/template.html
package stepgraph

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
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
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

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// ModTime return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _staticTemplateHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x58\xdf\x6f\xe3\xb8\x11\x7e\xf7\x5f\x31\x55\x5a\xd8\xb9\x8d\xe4\x24\x7b\x17\xf4\x6c\x59\xc0\xb6\xbb\x40\x1f\x92\x2e\x70\xd9\x45\x51\xdc\x1e\x0c\x5a\x1c\xd9\x4c\x28\x52\x25\x29\x7b\x5d\xc3\xff\x7b\x41\x52\x92\x25\xd9\xd9\xcb\x26\x68\x51\xe5\x21\x22\xc5\xf9\xe6\x9b\xe1\x37\xfc\xe1\xdd\x8e\x62\xc6\x04\x42\xb0\x42\x42\x51\x05\xfb\xfd\x20\xe6\x4c\x3c\x82\x42\x3e\x0b\xb4\xd9\x72\xd4\x2b\x44\x13\x80\xd9\x16\x38\x0b\x0c\x7e\x35\xe3\x54\xeb\x00\x56\x0a\xb3\x59\xf0\x50\x0a\x66\x22\xdb\x91\x0c\x62\x9d\x2a\x56\x98\xf6\xc8\x07\xb2\x26\xbe\x37\x00\xad\xd2\x59\xe0\x1b\xf3\x45\x29\x28\xc7\x28\x67\x22\x7a\xd0\x41\x12\x8f\x7d\x7f\x32\xd8\xed\x50\xd0\xfd\x7e\x30\x38\x30\x5b\x48\xba\x75\xbc\x1c\x9b\x64\x70\x86\x79\x61\xb6\xa1\xf3\x1c\xa6\x52\x18\xc2\x04\x2a\xd8\x0d\x00\x52\xc9\xa5\x9a\xc0\x19\xfe\xd9\xfe\x4d\x07\x00\x96\x45\x48\x38\x5b\x8a\x09\xa4\x28\x0c\x2a\xdb\x5b\x10\x4a\x99\x58\x86\x0b\x69\x8c\xcc\x27\x70\x75\x59\x7c\x9d\x0e\xf6\x83\x41\xb4\x62\x94\xa2\x08\x0d\x6a\xa3\x1d\xe4\x9a\x69\xb6\x60\x9c\x99\xed\xc4\xc2\x73\x52\x68\xb4\x10\x94\xe9\x82\x93\xed\x04\x84\x14\xd8\xb6\x75\x56\x27\xbe\x0a\xa9\x91\x63\x6a\xdc\xf7\x52\xa3\x0a\x7d\xfb\x30\x26\xc2\xaf\x05\x11\xb4\x8a\x25\x93\xc2\x84\x1b\x64\xcb\x95\x99\x2c\x24\xa7\xd3\xba\x4f\xb3\x7f\xe3\xe4\x2a\xfa\x09\x73\x0f\x5c\x5b\x4d\x38\xd1\x26\x94\x59\x68\xf3\xef\x20\xda\xc1\x2b\x0b\xe4\x0c\x0c\x8d\x32\xc2\x38\xd2\x4e\xca\xb2\xec\xc7\xcb\x1f\x2f\x9b\x01\x9c\x3c\x6e\x3b\xdf\x29\xfd\xf9\x67\x4a\xeb\xef\x05\xd1\xba\x07\x70\x73\x95\x65\x37\x57\xf5\x00\xfd\xc8\x8a\xe2\xc8\x05\xde\x5c\x7b\x88\x8a\x41\xc8\xc9\x56\x96\xe6\x02\xbc\xc3\xaa\xe9\x8c\x36\x8c\x9a\x95\x9d\x99\xcb\x3f\xd9\xc8\x17\x52\x51\x54\x61\x3d\x03\xed\xb9\x38\x82\x03\x43\xfb\x88\xc6\x33\xf1\x28\x13\xb8\x6c\x89\xc0\xb5\x6a\x8c\x52\x61\x28\x48\x8e\x8d\xbd\x6d\xf8\x20\x4a\xa5\x6d\x14\x85\x64\x5e\x45\x2e\x4e\x4f\x75\xc5\x0c\x86\xba\x20\x29\xda\xc9\x54\x39\xe1\xf0\x07\x96\x17\x52\x19\x22\x7c\xce\xc7\x3f\xc0\x3f\x10\x88\x42\x40\xb1\x24\x4b\xa4\xc0\x04\x10\x10\xb8\x46\x15\xa2\xb0\x3c\x60\x43\x14\xc8\x0c\x52\xa2\x53\x42\x11\x50\xa7\x84\x13\xc3\xa4\x00\xb2\x24\x4c\x68\x03\x77\xef\x6f\xe1\x87\xf1\xe0\xac\x8a\xd6\xd8\xca\x80\x04\x8c\x9a\xac\xe4\x1a\xd5\x05\x9c\x79\xd2\xfd\x0f\x3e\x76\x92\x3e\x2e\x95\x2c\x05\x0d\xab\x09\x29\x85\x46\xd3\x67\x6a\xc8\x82\x63\x3f\x9f\x0e\xcf\xa8\x26\x45\x56\x58\x2f\x86\xee\xcc\x4b\x83\xec\x89\xbf\x08\xb7\x43\x0b\x28\x5b\x37\xb3\x57\x77\x38\xb4\xba\xe8\x39\x66\x66\x02\xd7\xae\xe4\x0f\xbd\xae\x3e\x0e\xdd\x9d\x39\x2d\x14\x86\x1b\x45\x8a\xa6\x04\x33\x92\x33\xbe\x9d\x40\x2e\x85\x74\x63\xbe\xbd\xa8\x74\xf8\xb5\xc4\x59\xb5\x3b\xec\x1a\xeb\x9f\x6a\x6b\xa2\x94\xdc\x84\x2c\x95\x7e\x69\x59\xa3\x32\x2c\x25\xbc\x2e\xec\x9c\x51\xca\x5d\x19\xc4\xe3\x6a\x85\x8c\xeb\xd5\x34\x2b\x45\xea\x15\x44\xe9\x3d\xba\xd7\x0f\xd5\x6a\xa1\x47\xe7\x0e\x0e\x60\x4d\x14\xd4\x6b\x88\x86\x19\x50\x99\x96\x39\x0a\x13\xfd\xab\x44\xb5\xbd\x77\x8b\x94\x54\xef\x38\x1f\x0d\x8d\x8a\xb4\x87\x09\x6b\x8b\xe1\xf9\xb4\x41\x99\x73\x29\x8b\xf9\x15\xcc\xa0\x71\x3c\xaa\xc7\xd5\xde\xec\x53\xf7\x45\x52\xa4\x9c\xa5\x8f\x1d\x83\xf6\xc0\x1a\xd8\xcb\x64\x76\x30\x2c\x88\x42\x61\x3e\x70\x74\x44\x05\x7e\xad\xdf\xef\xd9\x82\x33\xb1\x9c\x1e\x41\xb8\xfc\xb5\x10\x3a\xc1\x8d\x86\xac\x8e\xa3\x7e\x58\x06\x23\xe7\x36\x4a\x39\xd1\xfa\x96\x69\x13\x55\x9b\x8d\x1e\x0d\xdb\xdb\xc4\xf0\xbc\x4f\xd9\x3e\x7d\x5b\x85\xb9\x5c\x63\xdf\x72\x7a\x64\x67\x79\x46\x4c\x08\x54\x9f\xac\x3a\x66\x30\xf4\x94\xe7\x1c\xb5\x1e\x76\xc7\xef\x3b\x2d\xe4\x1a\x9f\x41\x84\x50\xfa\x0a\x16\xb9\x54\xf8\x4d\x16\xba\xd8\x2e\xad\x2f\x97\x2c\x14\xe6\x73\x41\x89\x41\x3a\x6a\x39\xd9\xfb\xd7\xea\x5f\x26\x15\x8c\x9c\x7a\x18\xcc\xe0\xf2\x02\xe6\x04\x66\xf0\x4e\x29\xb2\x8d\x32\x25\xf3\x46\x41\xfa\x7c\x6a\xc7\xc4\x30\x27\x11\x47\xb1\x34\x2b\xdb\x7e\xf3\xa6\x9d\xfd\xb6\x96\x61\x06\x73\xf2\xeb\x9c\xfd\x76\xf0\x5c\x09\xf4\x20\xca\x8a\xc9\x60\xdf\x29\x95\x4f\xa8\xcd\xe9\x3a\x51\x72\xf3\x3b\x25\xd2\xdd\x46\x5a\xbb\xc8\x71\xa1\x5c\x77\x74\xaf\xe4\xa6\x1d\x89\x92\x9b\x67\x97\x87\xf6\xa2\x87\x99\xb3\xfa\x8e\x6a\xb0\xc3\x9f\x53\x08\x95\x83\x6f\x94\xc2\xe9\x22\x38\xb6\xeb\x96\xc1\xff\xa4\x00\x8e\x49\xb4\x4a\xe0\xff\x5d\xfc\x56\x70\xcf\xd4\xbd\x92\x9b\x27\x25\x7f\xed\xe4\x75\x5a\xed\xf7\x86\xca\xd2\x7c\x2c\x50\xf4\xd5\x6e\x2f\x02\xbf\x23\x77\x12\xc9\x02\x45\xa8\x1d\xc6\xb1\xc4\xdf\x76\xb4\x6b\xf1\xda\xac\x6d\xfb\xa4\xc8\xb1\x2f\x26\x8c\x0a\x85\x6b\x14\xe6\x3d\x66\xa4\xe4\x66\x74\x7e\xac\x69\xe3\x27\xcb\x81\x1e\xd7\x80\x9f\xd0\xbf\x7d\xba\xbb\x3d\xb6\x5c\x70\xb9\x80\x19\x08\xdc\xc0\x5f\xb8\x5c\x8c\x7e\x0d\xbe\x08\xff\x35\xb6\xd7\xa2\xa4\x6e\x01\xc4\x86\x19\x8e\xc9\xad\x5c\xea\x78\xec\xdf\x9b\xa1\xe3\xce\xd8\xd8\x6d\x59\x6e\x43\x9e\x7d\x09\x8e\x4f\x32\x67\x6f\x2f\xed\xdf\xb4\x3e\x1d\xbb\x03\xc7\xf4\xa9\xe3\xc5\x13\xe7\x91\x2f\x41\x12\xc0\x1b\x1f\xf9\x1b\x08\xe2\xb1\x75\x9a\x04\xbf\x5d\xc0\xce\xdd\xc2\x26\x30\x74\xd7\xb0\x95\xc9\xf9\x10\xf6\xbd\xa4\x6d\x98\xa0\x76\x95\x29\x50\x8c\x3e\xff\x72\x1b\xa5\x0a\x89\xc1\x8f\x8b\x07\x4c\xcd\xe7\x5f\x6e\x47\x36\x2d\xe7\x2f\x97\xae\x53\xcf\x33\xb5\xeb\xae\x9c\x4f\x89\xf7\xad\x17\xce\x09\xf5\x72\x49\xa8\x2d\xaf\x0a\xed\x78\xe9\x9e\xd6\xfd\x3d\x91\x1f\xfa\x8f\x4e\x45\xf6\x24\x55\x65\x86\x50\xfa\xc1\x8a\xce\x2e\x1a\xd6\x70\x34\x7c\xff\xf1\xee\xaf\xbe\xb4\x6f\x9d\xeb\xe1\x45\xc5\xe1\x7c\x3a\x68\x5f\x62\xff\x28\xca\x1c\x26\x33\xe0\x28\x20\xb2\x77\x57\x7b\x08\x65\xb4\xba\x2f\x1f\x6e\xad\x41\x32\xb0\xa2\xb2\x07\xe3\xd6\x67\xd7\x0e\xc0\x2d\x58\xb3\x20\xa7\x3c\xa4\xc4\x10\xdf\x0d\xb6\xf9\xa0\xfb\x3d\x7a\x45\xa8\xdc\x84\xe1\x35\x2d\x1c\xe6\x6e\xc7\x32\x58\x1a\x70\x44\x2e\xf7\x7e\x99\x8a\x8d\x72\x5e\xfc\xed\x2d\x34\xd5\x9d\xbf\x76\xe4\x9b\xd0\x3f\xe1\x39\x3c\x6f\x4e\x4f\x73\x9a\xcf\x53\xe4\x3c\x0c\x85\x14\xa1\x28\x73\x54\x2c\x3d\xec\xc2\xde\x59\x60\x75\xae\x0b\x22\x66\xc1\x55\x90\xc4\xab\x9b\xc4\x27\x69\xbf\x07\x6d\xb0\xd0\xf1\x78\x75\x93\xc4\x63\x43\x5f\xe5\x2c\x48\x62\xd6\x8e\xb0\xe9\xaf\xa1\xec\xd2\x1e\x2e\x4a\x63\xa4\x80\x9c\x18\x54\x8c\x70\x77\xb4\xd6\xd0\x3a\x66\xd7\x57\xf5\x20\x69\x2d\xfc\xf1\x98\xb5\x08\xc6\x63\xa3\xaa\x37\x7f\x38\x6d\xe7\xd5\xfd\x62\xd1\x64\xb5\x75\xd4\x6a\x52\xb9\xdb\x29\x22\x96\xe8\xb4\x71\x58\x5b\x54\xd2\xaa\xcf\xef\x4a\x80\x45\x77\xa7\x8c\x20\xd9\xed\xa2\x7b\x83\xc5\xdf\x49\x8e\xfb\x7d\x3b\xa3\xdf\x0b\xea\xa0\xde\x97\xca\xdd\x41\xfb\x50\x87\xf8\x9f\xa2\x5e\x4f\xf7\x75\x90\x74\x56\x9d\x96\xda\x3b\x54\xc3\x9c\x08\x96\xb9\x2c\xbd\x4a\xf9\x5d\x5f\xaf\x12\xfc\x0b\x27\xa3\xa5\xba\x9e\xe6\xef\xea\x10\x8f\xe4\xfe\x5a\x57\xb5\xf0\x4f\xa7\xf4\xbf\x58\x08\x2d\xe2\xe3\xae\x0a\x7a\xc5\xf1\x04\xb3\x67\x14\xcb\xe1\x69\xca\xa6\x49\xe4\x3f\xdf\xdd\xdd\xee\xf7\x56\x7f\x2f\xad\x97\xb8\x50\x58\xed\xd3\xc1\xe9\x5d\xd6\x15\x42\x55\x00\x2e\xc8\xfa\x47\xca\xa3\xf8\xdd\xf6\x3b\xe8\x75\x5a\x0e\x9d\xea\x78\xb2\x90\xda\xb8\x2d\xb4\x43\x77\x83\x16\x8f\x29\x5b\x1f\x7e\x2d\xfd\x4f\x00\x00\x00\xff\xff\x83\x51\x9c\xe7\xc7\x15\x00\x00")

func staticTemplateHtmlBytes() ([]byte, error) {
	return bindataRead(
		_staticTemplateHtml,
		"static/template.html",
	)
}

func staticTemplateHtml() (*asset, error) {
	bytes, err := staticTemplateHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "static/template.html", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
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
			return nil, fmt.Errorf("Asset %s can't read by error: %w", name, err)
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
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %w", name, err)
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
	"static/template.html": staticTemplateHtml,
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
	"static": {nil, map[string]*bintree{
		"template.html": {staticTemplateHtml, map[string]*bintree{}},
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