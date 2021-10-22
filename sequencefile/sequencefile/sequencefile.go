package sequencefile

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"

	// "io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/JackZxj/golang-demo/sequencefile/syscall"
)

// DEFAULT_BLOCK_SIZE means the max bytes of the block
const DEFAULT_BLOCK_SIZE = 1 << 17 // 128kb
// DEFAULT_RECORD_SIZE means the max record length of the block
const DEFAULT_RECORD_SIZE = 1 << 6 // 64
// MAX_FILE_SIZE means the max bytes of file
const MAX_FILE_SIZE = 1 << 24 // 16Mb

// SYNC_ESCAPE means length of sync
const SYNC_ESCAPE int32 = -1 // 0xffff_ffff_ffff_ffff
const SYNC_ESCAPE_SIZE = 4

// SYNC_HASH_SIZE means number of bytes in hash
const SYNC_HASH_SIZE = 16

// SYNC_SIZE means escape + hash
const SYNC_SIZE = SYNC_ESCAPE_SIZE + SYNC_HASH_SIZE

/*
 4byte    16byte      4byte           8byte   n1 byte 4byte           8byte   n2 byte  ...  4byte           8byte   n3 byte
----------------------------------------------------------------------------------------------------------------------------
| escape | sync hash | record length | index | value | record length | index | value | ... | record length | index | value |
----------------------------------------------------------------------------------------------------------------------------
                     ↑                               ↑                               ↑     ↑                               ↑
                     |---------- record 1 -----------|---------- record 2 -----------|     |---------- record k -----------|
↑                                                                                                                          ↑
|----------------------------------------------- block 1 ------------------------------------------------------------------|

record length = len(value)
*/

type SeqFiler interface {
	Read(index int64) ([]byte, error)
	// write a byte slice, returns index or error
	Write(value []byte) (int, error)
	Size() uint64
	Flush() error
}

type Index struct {
	Filename string `json:"filename"`
	Seek     int64  `json:"seek"`
	Index    int64  `json:"index"`
}

type KeyType string

const (
	KeyTypeInt    KeyType = "int"
	KeyTypeInt32  KeyType = "int32"
	KeyTypeInt64  KeyType = "int64"
	KeyTypeString KeyType = "string"
)

type Meta struct {
	Sync     []byte   `json:"sync"`
	Indexs   []*Index `json:"indexs"`
	MinIndex int64    `json:"minIndex"`
	MaxIndex int64    `json:"maxIndex"`
	KeyType  KeyType  `json:"keyType"`
}

// type File struct {
// 	file *os.File
// 	rwm  sync.RWMutex
// }

type SeqFile struct {
	// files   map[string]*os.File
	meta    Meta
	wal     [][]byte
	walSize int
	rwm     sync.RWMutex

	endianOrder binary.ByteOrder
}

var partitionRegex = regexp.MustCompile(`^p-.+`)
var dataPath = "data"

var ErrFileCorrupt = errors.New("file is corrupt")
var ErrIndexNotFound = errors.New("index not found")

func (s *SeqFile) Init() error {
	if IsLittleEndian() {
		s.endianOrder = binary.LittleEndian
	} else {
		s.endianOrder = binary.BigEndian
	}

	// s.files = make(map[string]*os.File)

	// pwd, err := os.Getwd()
	// if err != nil {
	// 	return fmt.Errorf("get pwd: %w", err)
	// }
	// dataPath := filepath.Join(pwd, "data")
	f, err := os.Stat(dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dataPath, 0755)
			if err != nil {
				return fmt.Errorf("create data path: %w", err)
			}
			return nil
		}
		return fmt.Errorf("stat path: %w", err)
	}
	if !f.IsDir() {
		err = os.MkdirAll(dataPath, 0755)
		if err != nil {
			return fmt.Errorf("create data path: %w", err)
		}
		return nil
	}

	m := Meta{}
	mf, err := os.Open(filepath.Join(dataPath, "meta.json"))
	if err != nil {
		return fmt.Errorf("read metadata: %w", err)
	}
	defer mf.Close()
	decoder := json.NewDecoder(mf)
	if err := decoder.Decode(&m); err != nil {
		return fmt.Errorf("decode metadata: %w", err)
	}
	s.meta = m

	// fl, err := ioutil.ReadDir(dataPath)
	// if err != nil {
	// 	return fmt.Errorf("read data path: %w", err)
	// }
	// errs := []error{}
	// for i := range fl {
	// 	name := fl[i].Name()
	// 	if fl[i].IsDir() || name == "meta.json" || !partitionRegex.MatchString(name) {
	// 		continue
	// 	}
	// 	pf, err := os.Open(filepath.Join(dataPath, name))
	// 	if err != nil {
	// 		errs = append(errs, err)
	// 		continue
	// 	}
	// 	s.files[name] = pf
	// }
	// if len(errs) > 0 {
	// 	return fmt.Errorf("open partition: %v", errs)
	// }
	return nil
}

func (s *SeqFile) Read(index int64) ([]byte, error) {
	if index < 1 {
		return nil, fmt.Errorf("index should be [1, n]")
	}

	s.rwm.RLock()
	if index < s.meta.MinIndex {
		return nil, ErrIndexNotFound
	}
	// the target index is in wal
	if s.meta.MaxIndex < index {
		if int64(len(s.wal))+s.meta.MaxIndex < index {
			return nil, ErrIndexNotFound
		}
		return s.wal[index-s.meta.MaxIndex], nil
	}
	s.rwm.RUnlock()

	var block, blockNext *Index
	for _, i := range s.meta.Indexs {
		if i.Index <= index {
			block = i
			continue
		}
		if i.Index > index {
			blockNext = i
			break
		}
	}

	f, err := os.Open(filepath.Join(dataPath, block.Filename))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("the index file not exist")
		}
		return nil, fmt.Errorf("read index file: %w", err)
	}
	msz := 0
	if blockNext != nil {
		msz, err = MmapSize(int(blockNext.Seek - block.Seek))
		if err != nil {
			return nil, err
		}
	} else {
		fstat, err := f.Stat()
		if err != nil {
			return nil, fmt.Errorf("read index stat: %w", err)
		}
		msz, err = MmapSize(int(fstat.Size() - block.Seek))
		if err != nil {
			return nil, err
		}
	}

	mapped, err := syscall.MmapOffset(int(f.Fd()), msz, block.Seek)
	if err != nil {
		return nil, fmt.Errorf("failed to perform mmap: %w", err)
	}
	bytebuff := bytes.NewBuffer(mapped)
	var value []byte
	for bytebuff.Len() != 0 {
		length, err := s.readRecordLength(bytebuff)
		if err != nil {
			return nil, fmt.Errorf("read record Length: %w", err)
		}
		i, err := s.readInt64(bytebuff)
		if err != nil {
			return nil, fmt.Errorf("read index: %w", err)
		}
		v := bytebuff.Next(int(length))
		if i == index {
			if len(v) == int(length) {
				copy(value, v)
				return value, nil
			}
			return nil, ErrFileCorrupt
		}
	}
	return nil, ErrIndexNotFound
}

func (s *SeqFile) Write(value []byte) (int, error) {
	s.rwm.Lock()
	defer s.rwm.Unlock()
	if len(s.wal) < DEFAULT_RECORD_SIZE && s.walSize < DEFAULT_BLOCK_SIZE {
		s.wal = append(s.wal, value)
		s.walSize += len(value) + 8 + 4
	}
	//TODO
	return -1, nil
}

// func (s *SeqFile) checkSync(bytebuff *bytes.Buffer) error {
// 	length, err := s.readInt32(bytebuff)
// 	if err != nil {
// 		return err
// 	}
// 	if length != SYNC_ESCAPE {
// 		return errors.New("cannot find a SYNC_ESCAPE")
// 	}
// 	var tmpSync [SYNC_HASH_SIZE]byte
// 	n, err := bytebuff.Read(tmpSync[:])
// 	if err != nil || n < SYNC_HASH_SIZE {
// 		return ErrFileCorrupt
// 	}
// 	for i, c := range s.meta.Sync {
// 		if c != tmpSync[i] {
// 			return ErrFileCorrupt
// 		}
// 	}
// 	return nil
// }

func (s *SeqFile) readRecordLength(bytebuff *bytes.Buffer) (int32, error) {
	length, err := s.readInt32(bytebuff)
	if err != nil {
		return -1, fmt.Errorf("read record: %w", err)
	}
	if length != SYNC_ESCAPE {
		return length, nil
	}
	var tmpSync [SYNC_HASH_SIZE]byte
	n, err := bytebuff.Read(tmpSync[:])
	if err != nil || n < SYNC_HASH_SIZE {
		return -1, ErrFileCorrupt
	}
	for i, c := range s.meta.Sync {
		if c != tmpSync[i] {
			return -1, ErrFileCorrupt
		}
	}
	length, err = s.readInt32(bytebuff)
	if err != nil {
		return -1, fmt.Errorf("read record: %w", err)
	}
	return length, nil
}

func (s *SeqFile) readInt32(bytebuff *bytes.Buffer) (int32, error) {
	var data int32
	if err := binary.Read(bytebuff, s.endianOrder, &data); err != nil {
		return 0, fmt.Errorf("read int32: %w", err)
	}
	return data, nil
}

func (s *SeqFile) readInt64(bytebuff *bytes.Buffer) (int64, error) {
	var data int64
	if err := binary.Read(bytebuff, s.endianOrder, &data); err != nil {
		return 0, fmt.Errorf("read int64: %w", err)
	}
	return data, nil
}

// func (s *SeqFile) readInt32(mapped []byte) (int32, error) {
// 	bytebuff := bytes.NewBuffer(mapped)
// 	var data int32
// 	if err := binary.Read(bytebuff, s.endianOrder, &data); err != nil {
// 		return 0, fmt.Errorf("read int32: %w", err)
// 	}
// 	return data, nil
// }

func IsLittleEndian() bool {
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	return (b == 0x04)
}

func MmapSize(size int) (int, error) {
	// Verify the requested size is not above the maximum allowed.
	if size > MAX_FILE_SIZE {
		return 0, fmt.Errorf("mmap too large")
	}

	// Ensure that the mmap size is a multiple of the page size.
	// This should always be true since we're incrementing in MBs.
	pageSize := int64(os.Getpagesize())
	sz := int64(size)
	if (sz % pageSize) != 0 {
		sz = ((sz / pageSize) + 1) * pageSize
	}
	return int(sz), nil
}
