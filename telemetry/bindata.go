// Code generated by go-bindata. DO NOT EDIT.
// sources:
// 000001_message_type.up.sql (66B)
// 000002_bandwidth_protocol.up.sql (719B)
// doc.go (73B)

package telemetry

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
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
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
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

	info := bindataFileInfo{name: "000001_message_type.up.sql", size: 66, mode: os.FileMode(0664), modTime: time.Unix(1675197752, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xe2, 0x43, 0xcc, 0xef, 0xad, 0x5f, 0x44, 0x58, 0x8d, 0x47, 0x70, 0x5d, 0x23, 0x30, 0xe2, 0x1f, 0xdb, 0x4d, 0xad, 0x6e, 0xd9, 0xe7, 0x50, 0x19, 0x43, 0x1c, 0x37, 0x57, 0xea, 0xc6, 0x57, 0xab}}
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

	info := bindataFileInfo{name: "000002_bandwidth_protocol.up.sql", size: 719, mode: os.FileMode(0664), modTime: time.Unix(1675262994, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xfe, 0x83, 0x69, 0xab, 0x3e, 0xf5, 0x8d, 0x44, 0xb2, 0x6e, 0x52, 0x8d, 0x27, 0xe8, 0x95, 0x28, 0x3c, 0xea, 0x29, 0x93, 0x6d, 0xa3, 0x10, 0xde, 0x9b, 0xc8, 0xa6, 0xb9, 0x80, 0xa1, 0x3, 0x6f}}
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

	info := bindataFileInfo{name: "doc.go", size: 73, mode: os.FileMode(0664), modTime: time.Unix(1675197752, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xae, 0x4f, 0xb8, 0x11, 0x84, 0x79, 0xbb, 0x6c, 0xf, 0xed, 0xc, 0xfc, 0x18, 0x32, 0x9d, 0xf1, 0x7, 0x2c, 0x20, 0xde, 0xe9, 0x97, 0x0, 0x62, 0x9f, 0x5e, 0x24, 0xfc, 0x8e, 0xc2, 0xd9, 0x2d}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
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

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
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
	"000001_message_type.up.sql": _000001_message_typeUpSql,

	"000002_bandwidth_protocol.up.sql": _000002_bandwidth_protocolUpSql,

	"doc.go": docGo,
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
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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
	"doc.go":                           &bintree{docGo, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
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
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
