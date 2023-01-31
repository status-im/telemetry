// Code generated by go-bindata. DO NOT EDIT.
// sources:
// 000001_message_type.up.sql (66B)
// 000002_bandwidth_protocol.up.sql (657B)
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

var __000002_bandwidth_protocolUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x90\x4f\x4b\xc4\x30\x10\xc5\xef\xfd\x14\x73\x6c\xa1\x27\x61\x4f\x9e\x62\x77\xd4\xc1\x9a\x2e\x93\xac\xb8\x27\x09\xdb\x80\x85\x35\xd1\x76\xfa\xfd\xa5\x0d\x2a\xfe\xd9\x2a\x1e\xf6\x94\xc3\x7b\xcc\xcb\xef\x57\x31\x2a\x8b\x60\xd5\x45\x8d\x40\x97\xa0\x1b\x0b\x78\x4f\xc6\x1a\x78\xee\xa3\xc4\x7d\x3c\x18\x71\x32\xb0\x13\x0f\x79\x06\x00\xd0\xb5\x60\x90\x49\xd5\xb0\x61\xba\x55\xbc\x83\x1b\xdc\x95\x73\xf4\x18\x07\xa1\x16\xee\x14\x57\xd7\x8a\xf3\xb3\xd5\xaa\x98\x2f\xea\x6d\x5d\xa7\xc6\xdb\x51\xed\x9e\xfc\x52\xaf\x77\xe2\x29\xc0\xba\xd9\x4e\x1f\xdb\x30\x56\x64\xa8\xd1\x3f\xb4\x9a\x51\x7e\xab\xed\x7b\xef\xc4\xb7\x4a\x80\xb4\xc5\x2b\xe4\xaf\x79\x0c\x83\xf4\xae\x0b\xf2\x1d\xfa\x61\x0c\xdd\xcb\xe8\x21\x3d\x79\x22\x2c\x3f\x71\x94\x1f\x03\x45\x56\x9c\x67\xd9\x5f\xa5\xda\x28\xee\x30\x9c\x52\xab\x4c\x8b\x14\x8e\x78\x98\xd3\xc9\xe7\x11\x4d\xef\x1a\xd7\xca\xe2\xb2\xba\x84\xf6\x0f\x79\xaf\x01\x00\x00\xff\xff\xd6\x24\x1d\x9b\x91\x02\x00\x00")

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

	info := bindataFileInfo{name: "000002_bandwidth_protocol.up.sql", size: 657, mode: os.FileMode(0664), modTime: time.Unix(1675204316, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x82, 0x33, 0x32, 0x3e, 0xa6, 0x35, 0x56, 0x34, 0xa1, 0xe7, 0x90, 0x7, 0x9b, 0x6a, 0xa1, 0x5d, 0xfe, 0xf, 0xe4, 0x26, 0xec, 0xb6, 0x20, 0x9, 0xa9, 0x61, 0x78, 0x87, 0xfd, 0x14, 0xbf, 0x1c}}
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
