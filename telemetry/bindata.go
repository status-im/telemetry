// Code generated for package telemetry by go-bindata DO NOT EDIT. (@generated)
// sources:
// 000001_message_type.up.sql
// 000002_bandwidth_protocol.up.sql
// 000003_index_truncate.up.sql
// 000004_envelope.table.up.sql
// doc.go
package telemetry

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

// Mode return file modify time
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

var __000001_message_typeUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\x28\x4a\x4d\x4e\xcd\x2c\x4b\x4d\xf1\x4d\x2d\x2e\x4e\x4c\x4f\x2d\x56\x70\x74\x71\x51\x70\xf6\xf7\x09\xf5\xf5\x53\xc8\x85\x88\x85\x54\x16\xa4\x2a\x84\x39\x06\x39\x7b\x38\x06\x69\x18\x99\x9a\x6a\x5a\x73\x01\x02\x00\x00\xff\xff\xf4\x14\x08\x7c\x42\x00\x00\x00")

func _000001_message_typeUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000001_message_typeUpSql,
		"000001_message_type.up.sql",
	)
}

func _000001_message_typeUpSql() (*asset, error) {
	bytes, err := _000001_message_typeUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000001_message_type.up.sql", size: 66, mode: os.FileMode(436), modTime: time.Unix(1715855770, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000002_bandwidth_protocolUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x90\x3f\x4f\xf3\x30\x10\xc6\xf7\x7c\x8a\x1b\x5b\x29\xd3\x2b\x75\xea\xe4\x37\x39\xc0\xc2\x71\x2a\xdb\x41\x74\x42\x56\x72\x82\x48\x6d\x0c\xf1\x85\x81\x4f\x8f\x92\x50\x10\x7f\x5a\x10\x03\x93\x25\x3f\x8f\xee\xee\xf7\xcb\x0c\x0a\x87\xe0\xc4\x7f\x85\x20\xcf\x40\x97\x0e\xf0\x5a\x5a\x67\xe1\xbe\x0f\x1c\xea\xb0\xb3\xec\x39\x1a\xcf\x04\x8b\x04\x00\xa0\x6d\xc0\xa2\x91\x42\xc1\xc6\xc8\x42\x98\x2d\x5c\xe2\x36\x9d\xa2\xbb\x10\x59\x36\x70\x25\x4c\x76\x21\xcc\xe2\xdf\x6a\xb5\x9c\x26\xea\x4a\xa9\xb9\x71\x18\xaa\xfd\x9e\x4e\xf5\x7a\xcf\x24\x3b\xc8\xcb\x6a\x3c\x6c\x63\x30\x93\x56\x96\xfa\x8b\x56\x39\xf0\x77\xb5\xba\x27\xcf\xd4\x08\x06\xa9\x1d\x9e\xa3\xf9\x98\x87\x2e\x72\xef\xdb\x8e\x3f\x43\xdf\x0c\x5d\xfb\x30\x10\xcc\xcf\x62\x26\x4c\xdf\x71\xa4\x6f\x0b\x96\xc9\x72\x9d\x24\x3f\x95\xea\x02\xfb\x5d\xfc\x4b\xad\x3c\x6e\x94\xdd\x11\x0f\x53\x3a\xfa\x3c\xa2\xe9\x55\x63\x2e\x1c\x9e\x56\x37\xa3\xfd\x46\x9e\x50\x0e\xcd\x8b\xbb\x9e\x6a\x6a\x1f\xa9\x29\x28\x46\x7f\x4b\x11\x44\x9e\x43\x56\xaa\xaa\xd0\xb0\x9f\xff\x6c\xfb\x44\x87\x7b\xd7\xc9\x73\x00\x00\x00\xff\xff\x88\xcb\xe5\x75\xcf\x02\x00\x00")

func _000002_bandwidth_protocolUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000002_bandwidth_protocolUpSql,
		"000002_bandwidth_protocol.up.sql",
	)
}

func _000002_bandwidth_protocolUpSql() (*asset, error) {
	bytes, err := _000002_bandwidth_protocolUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000002_bandwidth_protocol.up.sql", size: 719, mode: os.FileMode(436), modTime: time.Unix(1715855770, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000003_index_truncateUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x91\xb1\x6a\x80\x30\x10\x86\xf7\x3c\xc5\x8d\x0a\x5d\x3a\x67\x4a\x35\x83\xa0\x11\xd2\x08\xdd\x24\xe8\x91\x0a\xb6\x91\xe4\x5a\xfa\xf8\xc5\x4a\x2b\x18\xb5\x6b\xfe\xef\xbe\xfb\xb9\x18\xdd\xa9\x42\x18\x09\x46\x3c\xd5\x12\x02\x0e\x38\x7d\xe2\xd8\x60\x8c\xd6\x61\xe4\xec\x00\x2c\xc1\x93\x1f\xfc\xfc\x4c\x96\xa2\xb6\x84\xf7\x84\xf1\x64\xe7\xd4\x72\x58\x23\x9c\x0b\xe8\x2c\xe1\xc8\x19\x2b\xb4\x5c\xc1\x4a\x95\xf2\x25\xe9\xd3\x0f\x01\x57\x4e\x10\xb4\x2a\x49\xb3\xbf\x34\xe7\xb7\x9e\x7d\x61\x1f\x3e\xde\x4f\x65\x3b\x92\xfd\x20\xf9\xb1\x5a\x72\x89\x7e\x1a\xbf\x1e\x57\x53\x92\x64\xbf\x2f\xca\xbe\xe1\x03\x5c\xb6\x3c\x39\xdd\xb9\x74\xcb\x6e\xb4\x4c\xd4\x46\xea\xab\x3f\x03\x2d\x95\x68\x24\x14\x6d\xdd\x35\x0a\x5e\x7d\xa4\xaa\x04\xd3\xc2\x82\x18\xaa\x92\x5f\x4f\x6f\x8b\xff\x9d\xff\x0e\x00\x00\xff\xff\x11\x42\xcb\x4c\x56\x02\x00\x00")

func _000003_index_truncateUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000003_index_truncateUpSql,
		"000003_index_truncate.up.sql",
	)
}

func _000003_index_truncateUpSql() (*asset, error) {
	bytes, err := _000003_index_truncateUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000003_index_truncate.up.sql", size: 598, mode: os.FileMode(436), modTime: time.Unix(1716989343, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000004_envelopeTableUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x91\xb1\x6e\x83\x30\x10\x86\x77\x9e\xe2\x46\x90\x98\x2a\x65\xca\x74\x85\x6b\x63\xc5\x98\xca\x98\xaa\x99\x2a\x02\xa7\x14\xa9\xb1\xa9\x0d\x91\xfa\xf6\x95\x42\x5a\x55\xa9\x42\x27\x0f\xff\xe7\xdf\xe7\xef\x32\x4d\x68\x08\x0c\xde\x4b\x02\xf1\x00\xaa\x34\x40\x2f\xa2\x32\x15\x78\x6e\xb9\x3f\x71\x47\xf6\xc4\xef\x6e\xe0\x00\x71\x04\x00\xd0\x77\x50\x91\x16\x28\xe1\x49\x8b\x02\xf5\x0e\xb6\xb4\x4b\xcf\xd1\x91\x43\x68\x0e\xbc\x69\xc2\x1b\x3c\xa3\xce\x36\xa8\xe3\xbb\xd5\x2a\x39\xd7\xaa\x5a\xca\x19\x0b\x6c\x47\x1c\x41\x28\x43\x8f\xa4\xaf\xc2\xd6\x73\x33\x72\x77\x33\x1f\xdd\xd0\xb7\x4b\xed\xc3\xb4\xaf\xa6\xbd\xf9\x0f\xbb\x7c\xcf\x6f\xf9\xb3\x16\xf9\x12\x69\x5d\xc7\xaa\x39\xf2\xe2\xa3\xde\xb5\x1c\x42\x6f\x0f\xe4\xbd\xf3\x4b\x68\x56\xaa\xca\x68\x14\xca\xfc\x55\xfc\x3a\xd9\xfe\x63\x62\x98\x8f\x78\x16\x95\xfe\xf6\x9a\x5e\x0d\x9e\xfe\x8c\x97\x44\xc9\x3a\x8a\x50\x1a\xd2\x97\x7d\x7e\xd7\x17\xf3\xf5\x00\x98\xe7\x90\x95\xb2\x2e\xd4\x4d\x4b\xeb\xe8\x2b\x00\x00\xff\xff\x9d\x3f\xc2\xc6\x13\x02\x00\x00")

func _000004_envelopeTableUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000004_envelopeTableUpSql,
		"000004_envelope.table.up.sql",
	)
}

func _000004_envelopeTableUpSql() (*asset, error) {
	bytes, err := _000004_envelopeTableUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000004_envelope.table.up.sql", size: 531, mode: os.FileMode(436), modTime: time.Unix(1715855770, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _docGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2c\xc9\x31\x12\x84\x20\x0c\x05\xd0\x9e\x53\xfc\x0b\x90\xf4\x7b\x9b\xac\xfe\xc9\x38\x20\x41\x4c\xe3\xed\x6d\xac\xdf\xb4\xad\x99\x13\xf7\xd5\x4b\x51\xf5\xf8\x39\x07\x97\x25\xe1\x51\xff\xc7\xd8\x2d\x0d\x75\x36\x47\xb2\xf3\x64\xae\x07\x35\x20\xa2\x1f\x8a\x07\x44\xcb\x1b\x00\x00\xff\xff\xb6\x03\x50\xe0\x49\x00\x00\x00")

func docGoBytes() ([]byte, error) {
	return bindataRead(
		_docGo,
		"doc.go",
	)
}

func docGo() (*asset, error) {
	bytes, err := docGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "doc.go", size: 73, mode: os.FileMode(436), modTime: time.Unix(1715855770, 0)}
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
	"000001_message_type.up.sql":       _000001_message_typeUpSql,
	"000002_bandwidth_protocol.up.sql": _000002_bandwidth_protocolUpSql,
	"000003_index_truncate.up.sql":     _000003_index_truncateUpSql,
	"000004_envelope.table.up.sql":     _000004_envelopeTableUpSql,
	"doc.go":                           docGo,
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
	"000001_message_type.up.sql":       &bintree{_000001_message_typeUpSql, map[string]*bintree{}},
	"000002_bandwidth_protocol.up.sql": &bintree{_000002_bandwidth_protocolUpSql, map[string]*bintree{}},
	"000003_index_truncate.up.sql":     &bintree{_000003_index_truncateUpSql, map[string]*bintree{}},
	"000004_envelope.table.up.sql":     &bintree{_000004_envelopeTableUpSql, map[string]*bintree{}},
	"doc.go":                           &bintree{docGo, map[string]*bintree{}},
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
