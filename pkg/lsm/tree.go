package lsm

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	rbtree "github.com/pedrogao/RbTree"

	"github.com/pedrogao/plib/pkg/collection"
)

type keyType string

func (n keyType) LessThan(b interface{}) bool {
	value, _ := b.(keyType)
	return n < value
}

// Tree LSM tree(og structure tree)
type Tree struct {
	appendLog   *AppendLog
	bloomFilter *collection.BloomFilter
	segments    []string
	index       *rbtree.Tree // 磁盘文件稀疏索引
	memtable    *SizedMap

	threshold         int
	sparsityFactor    int
	bfNumItems        int
	bfFalsePosProb    float64
	segmentsDirectory string
	walBasename       string
	currentSegment    string
}

// metadata 元数据
type metadata struct {
	segments       []string
	currentSegment string
	index          *rbtree.Tree
	bloomFilter    *collection.BloomFilter
}

func (m *metadata) load(bytes []byte) error {
	data := map[string]any{}
	err := jsoniter.Unmarshal(bytes, &data)
	if err != nil {
		return fmt.Errorf("json unmarshal err: %s", err)
	}

	m.segments = data["segments"].([]string)
	m.currentSegment = data["segments"].(string)
	indexMap := data["index"].(map[string]any)

	m.index = rbtree.NewTree()
	for k, v := range indexMap {
		m.index.Insert(keyType(k), v)
	}

	bloomData := data["bloom"].([]byte)
	err = m.bloomFilter.UnPack(bloomData)
	if err != nil {
		return fmt.Errorf("bloom unpack err: %s", err)
	}

	return nil
}

func (m *metadata) dump() ([]byte, error) {
	data := map[string]any{
		"segments":       m.segments,
		"currentSegment": m.currentSegment,
	}
	indexMap := map[string]any{}
	iter := m.index.Iterator()
	for iter != nil {
		k := iter.Key.(keyType)
		indexMap[string(k)] = iter.Value
		iter = iter.Next()
	}
	data["index"] = indexMap
	bloom := m.bloomFilter.Pack()
	data["bloom"] = bloom

	bytes, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal err: %s", err)
	}
	return bytes, nil
}

type indexItem struct {
	segment string
	offset  int64
	val     any
}

// NewTree Initialize a new LSM tree
// - A first segment called segment_basename
// - A segments directory called segments_directory
// - A memtable write ahead log (WAL) called wal_basename
func NewTree(segmentBasename, segmentsDirectory,
	walBasename string) (*Tree, error) {
	// create lsm tree
	tree := &Tree{
		segments:          make([]string, 0),
		index:             rbtree.NewTree(),
		memtable:          NewSizedMap(),
		threshold:         1000000,
		sparsityFactor:    100,
		bfNumItems:        1000000,
		bfFalsePosProb:    0.2,
		segmentsDirectory: segmentsDirectory,
		walBasename:       walBasename,
		currentSegment:    segmentBasename,
	}

	// create bloom filter
	bloomFilter := collection.NewBloomFilter(tree.bfNumItems, tree.bfFalsePosProb)
	tree.bloomFilter = bloomFilter

	// create the segments directory
	if _, err := os.Stat(segmentsDirectory); err != nil && os.IsNotExist(err) {
		// directory not exist
		err = os.MkdirAll(segmentsDirectory, 0600)
		if err != nil {
			return nil, fmt.Errorf("make dir: %s err: %s", segmentsDirectory, err)
		}
	}

	// create write ahead log.
	appendLog, err := NewAppendLog(tree.memtableWalPath())
	if err != nil {
		return nil, fmt.Errorf("new wal: %s err: %s", tree.memtableWalPath(), err)
	}
	tree.appendLog = appendLog

	err = tree.loadMetadata()
	if err != nil {
		return nil, err
	}
	err = tree.restoreMemtable()
	if err != nil {
		return nil, err
	}
	return tree, nil
}

func (t *Tree) Set(key, value string) error {
	entry := t.toLogEntry(key, value)
	node := t.memtable.Get(key)
	if node != nil {
		if err := t.appendLog.WriteString(entry); err != nil {
			return err
		}
		t.memtable.Set(key, value)
	}
	additionalSize := len(key) + len(value)
	if t.memtable.GetTotalSize()+additionalSize > t.threshold {
		err := t.compact()
		if err != nil {
			return fmt.Errorf("compact err: %s", err)
		}
		err = t.flushMemtableToDisk(t.currentSegmentPath())
		if err != nil {
			return fmt.Errorf("flushMemtableToDisk err: %s", err)
		}
		t.memtable = NewSizedMap()
		if err := t.appendLog.Clear(); err != nil {
			return err
		}
		t.segments = append(t.segments, t.currentSegment)
		t.currentSegment = t.incrementedSegmentName()
	}
	if err := t.appendLog.WriteString(entry); err != nil {
		return err
	}
	t.memtable.Set(key, value)
	return nil
}

func (t *Tree) Get(key string) (string, error) {
	if got := t.memtable.Get(key); got != nil {
		return got.(string), nil
	}

	if !t.bloomFilter.Check(key) {
		return "", nil
	}

	val := t.index.Floor(keyType(key))
	if val == nil {
		return t.searchAllSegments(key)
	}

	item := val.(*indexItem)
	segment := item.segment
	offset := item.offset
	path := t.segmentPath(segment)
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open segment file err: %s", err)
	}
	_, err = file.Seek(offset, 0)
	if err != nil {
		return "", fmt.Errorf("seek segment file err: %s", err)
	}
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read segment file err: %s", err)
	}
	parts := strings.Split(line, ",")
	if len(parts) != 2 {
		return "", fmt.Errorf("segment file data err: %v", parts)
	}
	if parts[0] == key {
		return parts[1], nil
	}
	return "", fmt.Errorf("%s not found", key)
}

func (t *Tree) searchAllSegments(key string) (string, error) {
	for _, segment := range t.segments {
		val, err := t.searchSegment(key, segment)
		if err != nil {
			return "", err
		}
		if val != "" {
			return val, nil
		}
	}
	return "", nil
}

func (t *Tree) searchSegment(key, segment string) (string, error) {
	path := t.segmentPath(segment)
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open segment file err: %s", err)
	}
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", fmt.Errorf("read segment file err: %s", err)
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return "", fmt.Errorf("segment file data err: %v", parts)
		}
		if parts[0] == key {
			return parts[1], nil
		}
	}

	return "", nil
}

func (t *Tree) compact() error {
	keysOnDisk := map[string]struct{}{}
	for k := range t.memtable.inner {
		if t.bloomFilter.Check(k) {
			keysOnDisk[k] = struct{}{}
		}
	}
	return t.deleteKeysFromSegments(keysOnDisk, t.segments)
}

func (t *Tree) deleteKeysFromSegments(deletionKeys map[string]struct{},
	segments []string) error {
	for _, segment := range segments {
		segmentPath := t.segmentPath(segment)
		err := t.deleteKeysFromSegment(deletionKeys, segmentPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Tree) deleteKeysFromSegment(deletionKeys map[string]struct{},
	segmentPath string) error {
	tempPath := segmentPath + "_temp"
	input, err := os.Open(segmentPath)
	if err != nil {
		return fmt.Errorf("open segment file err: %s", err)
	}
	output, err := os.Open(tempPath)
	if err != nil {
		return fmt.Errorf("open segment temp file err: %s", err)
	}

	reader := bufio.NewReader(input)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("read memtable file err: %s", err)
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return fmt.Errorf("memtable file data err: %v", parts)
		}
		_, ok := deletionKeys[parts[0]]
		if !ok {
			_, err = output.WriteString(line)
			if err != nil {
				return fmt.Errorf("write segment temp file err: %s", err)
			}
		}
	}

	err = os.Remove(segmentPath)
	if err != nil {
		return fmt.Errorf("remove segment file err: %s", err)
	}
	err = os.Rename(tempPath, segmentPath)
	if err != nil {
		return fmt.Errorf("rename segment file err: %s", err)
	}
	return nil
}

func (t *Tree) flushMemtableToDisk(path string) error {
	sparsityCounter := t.sparsity()
	var keyOffset int64 = 0
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file: %s err: %s", path, err)
	}
	for k, v := range t.memtable.inner {
		entry := t.toLogEntry(k, v.(string))
		if sparsityCounter == 1 {
			t.index.Insert(keyType(k), indexItem{
				segment: path,
				offset:  keyOffset,
				val:     v,
			})
			sparsityCounter = t.sparsity() + 1
		}
		t.bloomFilter.Add(k)
		_, err := file.WriteString(entry)
		if err != nil {
			return fmt.Errorf("write %s err: %s", path, err)
		}
		keyOffset += int64(binary.Size(entry))
		sparsityCounter -= 1
	}
	return nil
}

func (t *Tree) loadMetadata() error {
	path := t.metadataPath()
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return nil
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read meta data err: %s", err)
	}
	meta := &metadata{}
	err = meta.load(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tree) saveMetadata() error {
	m := &metadata{
		segments:       t.segments,
		currentSegment: t.currentSegment,
		index:          t.index,
		bloomFilter:    t.bloomFilter,
	}
	bytes, err := m.dump()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(t.metadataPath(), bytes, 0666)
	if err != nil {
		return fmt.Errorf("write file err: %s", err)
	}
	return nil
}

func (t *Tree) restoreMemtable() error {
	path := t.memtableWalPath()
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open memtable file err: %s", err)
	}
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("read memtable file err: %s", err)
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return fmt.Errorf("memtable file data err: %v", parts)
		}
		t.memtable.Set(parts[0], parts[1])
	}
	return nil
}

func (t *Tree) merge(segment1, segment2 string) error {
	path1 := t.segmentsDirectory + segment1
	path2 := t.segmentsDirectory + segment2
	newPath := t.segmentsDirectory + "temp"

	s0, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("open file err: %s", err)
	}

	s1, err := os.OpenFile(path1, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("open file err: %s", err)
	}

	s2, err := os.OpenFile(path2, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("open file err: %s", err)
	}

	reader1 := bufio.NewReader(s1)
	reader2 := bufio.NewReader(s2)

	for {
		line1, err := reader1.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("read segment file err: %s", err)
		}
		line2, err := reader2.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("read segment file err: %s", err)
		}

		if line1 == "" && line2 == "" {
			break
		}

		parts1 := strings.Split(line1, ",")
		if len(parts1) != 2 {
			return fmt.Errorf("segment file data err: %v", parts1)
		}

		parts2 := strings.Split(line2, ",")
		if len(parts2) != 2 {
			return fmt.Errorf("segment file data err: %v", parts2)
		}

		key1 := parts1[0]
		key2 := parts2[0]

		if key1 == "" || key1 == key2 {
			s0.WriteString(line2)
			line1, err = reader1.ReadString('\n')
			if err != nil {
				return fmt.Errorf("read file err: %s", err)
			}

			line2, err = reader2.ReadString('\n')
			if err != nil {
				return fmt.Errorf("read file err: %s", err)
			}
		} else if key2 == "" || key1 < key2 {
			s0.WriteString(line1)
			line1, err = reader1.ReadString('\n')
			if err != nil {
				return fmt.Errorf("read file err: %s", err)
			}
		} else {
			s0.WriteString(line2)
			line2, err = reader2.ReadString('\n')
			if err != nil {
				return fmt.Errorf("read file err: %s", err)
			}
		}
	}
	err = os.Remove(path1)
	if err != nil {
		return fmt.Errorf("remove file err: %s", err)
	}
	err = os.Remove(path2)
	if err != nil {
		return fmt.Errorf("remove file err: %s", err)
	}
	err = os.Rename(newPath, path1)
	if err != nil {
		return fmt.Errorf("rename file err: %s", err)
	}
	return nil
}

func (t *Tree) repopulateIndex() error {
	t.index = rbtree.NewTree()
	for _, segment := range t.segments {
		path := t.segmentPath(segment)
		counter := t.sparsity()
		bytes := 0
		file, err := os.OpenFile(path, os.O_RDONLY, 0666)
		if err != nil {
			return fmt.Errorf("open file: %s err: %s", path, err)
		}

		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return fmt.Errorf("read segment file err: %s", err)
			}
			parts := strings.Split(line, ",")
			if len(parts) != 2 {
				return fmt.Errorf("segment file data err: %v", parts)
			}
			if counter == 1 {
				t.index.Insert(keyType(parts[0]), indexItem{
					segment: segment,
					offset:  int64(bytes),
					val:     parts[1],
				})
				counter = t.sparsity() + 1
			}
			bytes += binary.Size(line)
			counter -= 1
		}
	}
	return nil
}

func (t *Tree) incrementedSegmentName() string {
	parts := strings.Split(t.currentSegment, "-")
	if len(parts) != 2 {
		panic("segment name not valid")
	}
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		panic("segment name not valid")
	}
	return fmt.Sprintf("%s-%d", parts[0], num+1)
}

func (t *Tree) setThreshold(threshold int) {
	t.threshold = threshold
}

func (t *Tree) setSparsityFactor(factor int) {
	t.sparsityFactor = factor
}

func (t *Tree) setBloomFilterNumItems(numItems int) {
	t.bfNumItems = numItems
	t.bloomFilter = collection.NewBloomFilter(t.bfNumItems, t.bfFalsePosProb)
}

func (t *Tree) setBloomFilterFalsePosProb(probability float64) {
	t.bfFalsePosProb = probability
	t.bloomFilter = collection.NewBloomFilter(t.bfNumItems, t.bfFalsePosProb)
}

func (t *Tree) sparsity() int {
	return t.threshold / t.sparsityFactor
}

func (t *Tree) toLogEntry(key, value string) string {
	return key + "," + value + "\n"
}

// Returns the path to the memtable write ahead log.
func (t *Tree) memtableWalPath() string {
	return t.segmentsDirectory + t.walBasename
}

// Returns the path to the memtable write ahead log.
func (t *Tree) currentSegmentPath() string {
	return t.segmentsDirectory + t.currentSegment
}

// Returns the path to the given segment_name.
func (t *Tree) segmentPath(segmentName string) string {
	return t.segmentsDirectory + segmentName
}

// Returns the path to the metadata backup file.
func (t *Tree) metadataPath() string {
	return t.segmentsDirectory + "database_metadata"
}
