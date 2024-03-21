package queue

import (
	"encoding/binary"
	"errors"
	"regexp"
	"sync"

	"github.com/cockroachdb/pebble"
)

var (
	// ErrIsEmpty is returned when queue is empty
	ErrIsEmpty = errors.New("queue: is empty")

	// ErrIDOutOfBounds is returned when queue is is out of bounds
	ErrIDOutOfBounds = errors.New("queue: ID is out of bounds")

	// ErrInvalidName is returned when queue name is not valid
	ErrInvalidName = errors.New("queue: name is not alphanumeric")

	// ErrNameTooLong means that queue name is longer then allowed limit
	ErrNameTooLong = errors.New("queue: name is too long")

	// ErrInvalidHeadValue is returned when there is an attempt
	// to assign invalid queue head value
	ErrInvalidHeadValue = errors.New("queue: head can not be less then zero")

	// ErrSharedFlush means that there was an attempt to flush shared queue
	ErrSharedFlush = errors.New("queue: can't flush shared queue")
)

var validQueueNameRegex = regexp.MustCompile(`[^a-zA-Z0-9_\-\:]+`)

// Consumer represents a queue consumer
type Consumer interface {
}

// make sure Queue implements Consumer interface
var _ Consumer = (*Queue)(nil)

// Queue represents a persistent FIFO structure
// that stores the data in leveldb
type Queue struct {
	sync.RWMutex
	Name     string
	DataDir  string
	db       *pebble.DB
	opts     *Options
	head     uint64
	tail     uint64
	isOpened bool
	isShared bool
}

// Options represents queue options
type Options struct {
	KeyPrefix []byte
}

// Item represents a queue item
type Item struct {
	ID    uint64
	Key   []byte
	Value []byte
}

// Open creates a queue and opens underlying leveldb database
func Open(name string, dataDir string, opts *Options) (*Queue, error) {
	q := &Queue{
		Name:     name,
		DataDir:  dataDir,
		db:       &pebble.DB{},
		opts:     opts,
		head:     0,
		tail:     0,
		isOpened: false,
		isShared: false,
	}
	return q, q.open()
}

// Path returns leveldb database file path
func (q *Queue) Path() string {
	return q.DataDir + "/" + q.Name
}

func (q *Queue) open() error {
	if validQueueNameRegex.MatchString(q.Name) {
		return ErrInvalidName
	}

	if len(q.Name) > 100 {
		return ErrNameTooLong
	}

	if !q.isShared {
		var err error
		q.db, err = pebble.Open(
			q.Path(),
			nil,
		)
		if err != nil {
			return err
		}
	}
	q.isOpened = true
	return q.initialize()
}

// Enqueue adds new value to the queue
func (q *Queue) Enqueue(value []byte) error {
	q.Lock()
	defer q.Unlock()

	err := q.db.Set(q.dbKey(q.tail+1), value, nil)
	if err == nil {
		q.tail++
	}
	return err
}

func (q *Queue) dbKey(id uint64) []byte {
	if len(q.opts.KeyPrefix) == 0 {
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, id)
		return key
	}
	key := make([]byte, len(q.opts.KeyPrefix)+8)
	copy(key[0:len(q.opts.KeyPrefix)], q.opts.KeyPrefix)
	binary.BigEndian.PutUint64(key[len(q.opts.KeyPrefix):], id)
	return key
}

// BytesLimit returns the limit that satisfy the given prefix.
// This only applicable for the standard 'bytes comparer'.
func BytesLimit(prefix []byte) []byte {
	var limit []byte
	for i := len(prefix) - 1; i >= 0; i-- {
		c := prefix[i]
		if c < 0xff {
			limit = make([]byte, i+1)
			copy(limit, prefix)
			limit[i] = c + 1
			break
		}
	}
	// [prefix, limit)
	return limit
}

// Convert queue name to uint64
func (q *Queue) dbKeyToID(key []byte) uint64 {
	return binary.BigEndian.Uint64(key[len(q.opts.KeyPrefix):])
}

func (q *Queue) initialize() error {
	snapshot := q.db.NewSnapshot()
	defer snapshot.Close()

	iter, err := snapshot.NewIter(&pebble.IterOptions{
		LowerBound: []byte(q.opts.KeyPrefix),
		UpperBound: BytesLimit(q.opts.KeyPrefix),
	})

	if err != nil {
		return err
	}
	defer iter.Close()

	if iter.First() {
		q.head = q.dbKeyToID(iter.Key()) - 1
	} else {
		q.head = 0
	}

	if iter.Last() && iter.Valid() {
		q.tail = q.dbKeyToID(iter.Key())
	} else {
		q.tail = 0
	}

	return iter.Error()
}
