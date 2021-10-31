package sequencefile

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
	"unsafe"

	// "io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	mmapgo "github.com/edsrzf/mmap-go"
)

// DEFAULT_BLOCK_SIZE means the max bytes of the block
const DEFAULT_BLOCK_SIZE = 1 << 17 // 128kb
// DEFAULT_RECORD_NUM_IN_BLOCK means the max record length of the block
const DEFAULT_RECORD_NUM_IN_BLOCK = 1 << 6 // 64
// MAX_FILE_SIZE means the max bytes of file
const MAX_FILE_SIZE = 1 << 24 // 16Mb

// SYNC_ESCAPE means length of sync
const SYNC_ESCAPE int32 = -1 // 0xffff_ffff_ffff_ffff
const SYNC_ESCAPE_SIZE = 4

// SYNC_HASH_SIZE means number of bytes in hash
const SYNC_HASH_SIZE = 16

// SYNC_SIZE means escape + hash
const SYNC_SIZE = SYNC_ESCAPE_SIZE + SYNC_HASH_SIZE

const DEFAULT_MAX_DATA_SIZE = DEFAULT_BLOCK_SIZE - SYNC_SIZE - 4 - 8

/*
 4byte    16byte      4byte           8byte   n1 byte 4byte           8byte   n2 byte  ...  4byte           8byte   n3 byte
----------------------------------------------------------------------------------------------------------------------------
| escape | sync hash | record length | index | value | record length | index | value | ... | record length | index | value |
----------------------------------------------------------------------------------------------------------------------------
                     ↑                               ↑                               ↑     ↑                               ↑
                     |---------- record 1 -----------|---------- record 2 -----------|     |---------- record k -----------|
↑                                                                                                                          ↑
|----------------------------------------------- block 1 ------------------------------------------------------------------|

record length = len(value) + 8
*/

type SeqFiler interface {
	Read(index int64) ([]byte, error)
	// write a byte slice, returns index or error
	Write(value []byte) (int64, error)
	Size() uint64
	Close() error
}

type Index struct {
	Filename string `json:"filename"`
	Seek     int64  `json:"seek"`
	Index    int64  `json:"index"`
}

type Meta struct {
	Sync     []byte   `json:"sync"`
	Indexs   []*Index `json:"indexs"`
	MinIndex int64    `json:"minIndex"`
	MaxIndex int64    `json:"maxIndex"`
}

type SeqFile struct {
	meta Meta
	walf *os.File
	// wal  mmapgo.MMap
	rwm sync.RWMutex

	walIndexs       []*Index
	lastSyncPos     int64
	curIndex        int64
	blockRemainSize int
	blockRemainNum  int

	closed      bool
	wg          sync.WaitGroup
	endianOrder binary.ByteOrder
}

var partitionRegex = regexp.MustCompile(`^p-.+`)
var dataPath = "data"

var ErrSeqFileClosed = errors.New("SeqFile has been closed")
var ErrWalCorrupt = errors.New("wal is corrupt")
var ErrFileCorrupt = errors.New("file is corrupt")
var ErrIndexNotFound = errors.New("index not found")
var ErrEmptyIndexs = errors.New("empty indexs")
var ErrEmptyData = errors.New("empty data")
var ErrDataOversize = errors.New(fmt.Sprintf("data over size, max size: %dbyte", DEFAULT_MAX_DATA_SIZE))

func (s *SeqFile) Init() error {
	if IsLittleEndian() {
		s.endianOrder = binary.LittleEndian
	} else {
		s.endianOrder = binary.BigEndian
	}

	newDataPath := func() error {
		err := os.MkdirAll(dataPath, 0755)
		if err != nil {
			return fmt.Errorf("create data path: %w", err)
		}
		s.curIndex = -1
		s.meta.MaxIndex = -1
		return s.initWal()
	}

	f, err := os.Stat(dataPath)
	if err != nil {
		// data path not exists, create it
		if os.IsNotExist(err) {
			return newDataPath()
		}
		return fmt.Errorf("stat path: %w", err)
	}
	if !f.IsDir() {
		return fmt.Errorf("cannot create directory %q: File exists", dataPath)
	}

	// if data path exists,
	// read meta.json
	m := Meta{}
	var decoder *json.Decoder
	initSync := true
	mf, err := os.Open(filepath.Join(dataPath, "meta.json"))
	if err != nil {
		if os.IsNotExist(err) {
			s.curIndex = -1
			goto INIT_WAL
		}
		return fmt.Errorf("read metadata: %w", err)
	}
	defer mf.Close()
	decoder = json.NewDecoder(mf)
	if err = decoder.Decode(&m); err != nil {
		return fmt.Errorf("decode metadata: %w", err)
	}
	s.meta = m
	initSync = false
	s.curIndex = s.meta.MaxIndex

INIT_WAL:
	// read wal
	walf, err := os.OpenFile(filepath.Join(dataPath, "wal"), os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return s.initWal()
		}
		return fmt.Errorf("open wal file: %w", err)
	}
	s.walf = walf
	wal, err := mmapgo.Map(s.walf, mmapgo.RDWR, 0)
	defer wal.Unmap()
	if err != nil {
		return fmt.Errorf("mmap wal: %w", err)
	}
	s.walIndexs = make([]*Index, 0)
	bytebuf := bytes.NewBuffer(wal)
	var curPos int64
	for bytebuf.Len() > 0 {
		// read ESCAPE or record length
		length, err := s.readInt32(bytebuf)
		if err != nil {
			return fmt.Errorf("read record: %w", err)
		}
		curPos += SYNC_ESCAPE_SIZE
		if length != SYNC_ESCAPE {
			if initSync {
				return ErrWalCorrupt
			}
			// next = index + value
			next := bytebuf.Next(int(length))
			if l := len(next); l == int(length) {
				curPos += int64(l)
				s.curIndex, err = s.readInt64(bytes.NewBuffer(next))
				if err != nil {
					return fmt.Errorf("read wal index: %w", err)
				}
				continue
			}
			return ErrWalCorrupt
		}
		sync := make([]byte, SYNC_HASH_SIZE)
		n, err := bytebuf.Read(sync)
		if err != nil || n < SYNC_HASH_SIZE {
			return ErrWalCorrupt
		}
		curPos += SYNC_HASH_SIZE
		if initSync {
			s.meta.Sync = sync
			initSync = false
		} else {
			for i, c := range s.meta.Sync {
				if c != sync[i] {
					return ErrWalCorrupt
				}
			}
		}
		s.lastSyncPos = curPos - SYNC_SIZE
		index := Index{
			Filename: "wal",
			Seek:     s.lastSyncPos,
			Index:    s.curIndex + 1,
		}
		s.walIndexs = append(s.walIndexs, &index)
	}
	return nil
}

func (s *SeqFile) Read(index int64) ([]byte, error) {
	var (
		block, blockNext *Index
		isWAL            bool
		err              error
		f                *os.File
		mapped           mmapgo.MMap
		bytebuff         *bytes.Buffer
		value            []byte
		mapStart, mapEnd int
		offStart, offEnd int
	)

	s.rwm.RLock()
	if s.closed {
		s.rwm.RUnlock()
		return nil, ErrSeqFileClosed
	}
	s.wg.Add(1)
	defer s.wg.Done()

	// fmt.Printf("%+v\n", s)

	if index < s.meta.MinIndex || index > s.curIndex {
		err = fmt.Errorf("index too small or too large: %w", ErrIndexNotFound)
		s.rwm.RUnlock()
		return nil, err
	}
	// the target index is in wal
	if s.meta.MaxIndex < index {
		block, blockNext, err = binarySearchIndex(s.walIndexs, index)
		// fmt.Println(*block)
		isWAL = true
	} else {
		block, blockNext, err = binarySearchIndex(s.meta.Indexs, index)
	}
	if !isWAL {
		// if not wal, files can be read without a mutex
		s.rwm.RUnlock()
	}
	if err != nil {
		goto ERR
	}

	f, err = os.Open(filepath.Join(dataPath, block.Filename))
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("the index file not exist")
			goto ERR
		}
		err = fmt.Errorf("read index file: %w", err)
		goto ERR
	}

	mapStart, err = MmapSizeFloor(int(block.Seek))
	if err != nil {
		goto ERR
	}
	offStart = int(block.Seek) - mapStart
	if blockNext != nil && blockNext.Filename == block.Filename {
		mapEnd, err = MmapSizeCeil(int(blockNext.Seek))
		if err != nil {
			goto ERR
		}
		offEnd = mapEnd - int(blockNext.Seek)
	} else {
		fstat, e := f.Stat()
		if e != nil {
			err = fmt.Errorf("read index stat: %w", e)
			goto ERR
		}
		mapEnd = int(fstat.Size())
	}

	mapped, err = mmapgo.MapRegion(f, mapEnd-mapStart, mmapgo.RDONLY, 0, int64(mapStart))
	if err != nil {
		err = fmt.Errorf("failed to perform mmap: %w", err)
		goto ERR
	}
	defer mapped.Unmap()
	// fmt.Println(mapStart, mapEnd, offStart, offEnd, len(mapped), mapped)
	bytebuff = bytes.NewBuffer(mapped[offStart : len(mapped)-offEnd])
	for bytebuff.Len() != 0 {
		length, e := s.readRecordLength(bytebuff)
		if e != nil {
			err = fmt.Errorf("read record Length: %w", e)
			goto ERR
		}
		// fmt.Println("Record length", length)
		i, e := s.readInt64(bytebuff)
		if e != nil {
			err = fmt.Errorf("read index: %w", e)
			goto ERR
		}
		// fmt.Println("index", i)
		l := int(length - 8) // real record length
		v := bytebuff.Next(l)
		// fmt.Printf("value: %s\n", v)
		if i == index {
			if len(v) == l {
				value = make([]byte, l)
				copy(value, v)
				if isWAL {
					s.rwm.RUnlock()
				}
				return value, nil
			}
			err = ErrFileCorrupt
			goto ERR
		}
	}
	err = ErrIndexNotFound

ERR:
	if isWAL {
		s.rwm.RUnlock()
	}
	return nil, err
}

func (s *SeqFile) Write(value []byte) (int64, error) {
	if len(value) == 0 {
		return -1, ErrEmptyData
	}
	if len(value) > DEFAULT_MAX_DATA_SIZE {
		return -1, ErrDataOversize
	}
	recordLength := len(value) + 8
	recordLengthBytes, err := Int32ToBytes(int32(recordLength), s.endianOrder)
	if err != nil {
		return -1, fmt.Errorf("convert recordLength: %w", err)
	}
	s.rwm.Lock()
	defer s.rwm.Unlock()
	if s.closed {
		return -1, ErrSeqFileClosed
	}
	s.wg.Add(1)
	defer s.wg.Done()
	fs, err := s.walf.Stat()
	if err != nil {
		return -1, fmt.Errorf("read wal stat: %w", err)
	}
	tail := fs.Size()
	// if this block has no space for this record, create new block
	if s.blockRemainSize < recordLength+4 || s.blockRemainNum < 1 {
		// if wal file has no space for new block
		if tail+DEFAULT_BLOCK_SIZE > MAX_FILE_SIZE {
			// mv wal data-xxx-xxx
			newFile := fmt.Sprintf("p-%d-%d", s.meta.MaxIndex+1, s.curIndex)
			err := os.Rename(filepath.Join(dataPath, "wal"), filepath.Join(dataPath, newFile))
			if err != nil {
				return -1, fmt.Errorf("mv wal to new data file: %w", err)
			}
			for i := range s.walIndexs {
				s.walIndexs[i].Filename = newFile
			}
			s.meta.Indexs = append(s.meta.Indexs, s.walIndexs...)
			s.meta.MaxIndex = s.curIndex
			b, err := json.Marshal(&s.meta)
			if err != nil {
				return -1, fmt.Errorf("failed to encode metadata: %w", err)
			}
			metaPath := filepath.Join(dataPath, "meta.json")
			if err := os.WriteFile(metaPath, b, 511); err != nil {
				return -1, fmt.Errorf("failed to write metadata to %s: %w", metaPath, err)
			}
			// init wal
			s.walIndexs = make([]*Index, 0)
			err = s.initWal()
			if err != nil {
				return -1, fmt.Errorf("init wal: %w", err)
			}
		} else {
			// sync()
			s.sync()
			if err != nil {
				return -1, fmt.Errorf("sync block: %w", err)
			}
		}
	}
	// write value
	// convert index to bytes
	indexBytes, err := Int64ToBytes(s.curIndex+1, s.endianOrder)
	if err != nil {
		return -1, fmt.Errorf("convert index: %w", err)
	}
	record := append(recordLengthBytes, indexBytes...)
	record = append(record, value...)
	// fmt.Println("got record:", record, "tail:", tail)
	err = s.writeWal(tail, record)
	if err != nil {
		return -1, fmt.Errorf("write wal: %w", err)
	}
	s.curIndex++
	s.blockRemainNum--
	// fmt.Println("block status:", s.blockRemainSize, s.blockRemainNum)
	return s.curIndex, nil
}

func (s *SeqFile) Size() uint64 {
	// TODO
	return 0
}

func (s *SeqFile) Close() error {
	s.rwm.Lock()
	defer s.rwm.Unlock()
	s.closed = true
	s.wg.Wait()
	err := s.walf.Close()
	if err != nil {
		return fmt.Errorf("unexpect error when closing wal file")
	}
	return nil
}

func (s *SeqFile) initWal() error {
	walf, err := os.OpenFile(filepath.Join(dataPath, "wal"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open wal file: %w", err)
	}
	s.walf = walf
	if s.meta.Sync == nil {
		s.meta.Sync = RandBytes(time.Now().UnixNano(), SYNC_HASH_SIZE)
	}
	s.lastSyncPos = -100
	err = s.sync()
	if err != nil {
		return fmt.Errorf("sync data: %w", err)
	}
	return nil
}

func (s *SeqFile) sync() error {
	// var tail int64
	// if s.wal != nil {
	// 	tail = int64(len(s.wal))
	// 	if s.lastSyncPos >= tail-SYNC_SIZE {
	// 		return nil
	// 	}
	// }
	fs, err := s.walf.Stat()
	if err != nil {
		return fmt.Errorf("read wal stat: %w", err)
	}
	tail := fs.Size()
	if s.lastSyncPos >= tail-SYNC_SIZE {
		return nil
	}
	sync := make([]byte, SYNC_SIZE)
	escape, err := Int32ToBytes(SYNC_ESCAPE, s.endianOrder)
	if err != nil {
		return fmt.Errorf("convert SYNC_ESCAPE: %w", err)
	}
	copy(sync[:SYNC_ESCAPE_SIZE], escape)
	copy(sync[SYNC_ESCAPE_SIZE:], s.meta.Sync)

	if err := s.writeWal(tail, sync); err != nil {
		return fmt.Errorf("write wal: %w", err)
	}

	s.lastSyncPos = tail
	s.blockRemainSize = DEFAULT_BLOCK_SIZE - SYNC_SIZE
	s.blockRemainNum = DEFAULT_RECORD_NUM_IN_BLOCK
	index := Index{
		Filename: "wal",
		Seek:     int64(tail),
		Index:    s.curIndex + 1,
	}
	s.walIndexs = append(s.walIndexs, &index)
	return nil
}

func (s *SeqFile) writeWal(tail int64, record []byte) error {
	err := s.walf.Truncate(tail + int64(len(record)))
	if err != nil {
		return fmt.Errorf("truncate wal: %w", err)
	}
	wal, err := mmapgo.Map(s.walf, mmapgo.RDWR, 0)
	if err != nil {
		return fmt.Errorf("mmap wal: %w", err)
	}
	defer wal.Unmap()
	copy(wal[tail:], record)
	// fmt.Printf("wal length: %d, wal value: %v, len: %d\n", len(wal), wal, tail+int64(len(record)))
	wal.Flush()
	return nil
}

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

func (s *SeqFile) readInt32(bytebuf *bytes.Buffer) (int32, error) {
	var data int32
	if err := binary.Read(bytebuf, s.endianOrder, &data); err != nil {
		return 0, fmt.Errorf("read int32: %w", err)
	}
	return data, nil
}

func (s *SeqFile) readInt64(bytebuf *bytes.Buffer) (int64, error) {
	var data int64
	if err := binary.Read(bytebuf, s.endianOrder, &data); err != nil {
		return 0, fmt.Errorf("read int64: %w", err)
	}
	return data, nil
}

func Int64ToBytes(data int64, order binary.ByteOrder) ([]byte, error) {
	bytebuf := bytes.NewBuffer([]byte{})
	if err := binary.Write(bytebuf, order, data); err != nil {
		return nil, fmt.Errorf("convert int64 to bytes: %w", err)
	}
	// fmt.Println("Int64ToBytes", data, bytebuf.Bytes())
	return bytebuf.Bytes(), nil
}

func Int32ToBytes(data int32, order binary.ByteOrder) ([]byte, error) {
	bytebuf := bytes.NewBuffer([]byte{})
	if err := binary.Write(bytebuf, order, data); err != nil {
		return nil, fmt.Errorf("convert int32 to bytes: %w", err)
	}
	// fmt.Println("Int32ToBytes", data, bytebuf.Bytes())
	return bytebuf.Bytes(), nil
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

func MmapSizeCeil(size int) (int, error) {
	return mmapSize(size, false)
}

func MmapSizeFloor(size int) (int, error) {
	return mmapSize(size, true)
}

func mmapSize(size int, floor bool) (int, error) {
	// Verify the requested size is not above the maximum allowed.
	if size > MAX_FILE_SIZE {
		return 0, fmt.Errorf("mmap too large, the MAX: %d", MAX_FILE_SIZE)
	}

	// Verify the requested size is not above the maximum allowed.
	if size < 0 {
		return 0, fmt.Errorf("mmap requires non-negative size, but got: %d", size)
	}

	pageSize := int64(os.Getpagesize())
	if size == 0 {
		if floor {
			return 0, nil
		}
		return int(pageSize), nil
	}
	sz := int64(size)
	if (sz % pageSize) != 0 {
		if floor {
			sz = (sz / pageSize) * pageSize
		} else {
			sz = ((sz / pageSize) + 1) * pageSize
		}
	}
	return int(sz), nil
}

// alpha ignores 'l','o','I','O','0','1'
var alpha = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// RandBytes generates a random string of fixed size.
func RandBytes(seed int64, size int) []byte {
	buf := make([]byte, size)
	rand.Seed(seed)
	for i := 0; i < size; i++ {
		buf[i] = alpha[rand.Intn(len(alpha))]
	}
	return buf
}

func binarySearchIndex(indexs []*Index, index int64) (*Index, *Index, error) {
	var low, high, mid int64
	low = 0
	high = int64(len(indexs)) - 1
	// fmt.Println("low and high:", low, high)
	if high < 0 {
		return nil, nil, ErrEmptyIndexs
	}
	if index < indexs[0].Index {
		return nil, nil, ErrIndexNotFound
	}
	if index >= indexs[high].Index {
		return indexs[high], nil, nil
	}
	for low <= high {
		mid = (high + low) / 2
		if indexs[mid].Index == index || (high-low) < 2 {
			break
		}
		if indexs[mid].Index > index {
			high = mid
		} else {
			low = mid
		}
	}
	return indexs[mid], indexs[mid+1], nil
}
