package cache

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	// QueueSelect params: SelfID = $1, Location = $2
	QueueSelect = "SELECT id FROM profiles WHERE $2 = ANY (location) AND id != $1"
	// QueueExceptSelect params: SelfID = $1, Location = $2, PrevLocation = $3
	QueueExceptSelect = "SELECT id FROM profiles WHERE $2 = ANY (location) AND id != $1 EXCEPT SELECT id FROM profiles WHERE $3 = ANY (location) AND id != $1"
)

const Timeout = 10 * time.Second

var ErrNotFound = errors.New("no rows in result set")

type Cache struct {
	queue      []int64
	lastAccess time.Time
	location   []string
	currIndex  int
}

var cacheMap = make(map[int64]Cache)

func InvalidateCache() {
	for {
		time.Sleep(10 * time.Minute)
		for k, v := range cacheMap {
			if time.Since(v.lastAccess) > 2*time.Hour {
				delete(cacheMap, k)
			}
		}
	}
}

func GetNext(db *pgxpool.Pool, self int64) (int64, error) {
	if _, exists := cacheMap[self]; !exists {
		err := populateCache(db, self)
		if err != nil {
			return 0, err
		}
	}
	for len(cacheMap[self].queue) == 0 {
		err := fillCache(db, self)
		if err != nil {
			return 0, err
		}
	}
	selfCache := cacheMap[self]
	selfCache.lastAccess = time.Now()

	id := selfCache.queue[0]
	selfCache.queue = selfCache.queue[1:]

	cacheMap[self] = selfCache
	return id, nil
}

func fillCache(db *pgxpool.Pool, selfID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	self := cacheMap[selfID]

	if self.currIndex+1 != len(self.location) {
		self.currIndex++
	} else {
		return populateCache(db, selfID)
	}

	rows, err := db.Query(ctx, QueueExceptSelect,
		selfID,
		self.location[self.currIndex],
		self.location[self.currIndex-1],
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 && self.currIndex+1 == len(self.location) {
		return ErrNotFound
	}

	self.queue = ids
	cacheMap[selfID] = self
	return nil
}

func populateCache(db *pgxpool.Pool, self int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	rows, err := db.Query(ctx, "SELECT location FROM profiles WHERE id = $1", self)
	if err != nil {
		return err
	}
	defer rows.Close()

	var location []string
	for rows.Next() {
		err = rows.Scan(&location)
		if err != nil {
			return err
		}
	}

	if len(location) == 0 {
		return ErrNotFound
	}

	rows, err = db.Query(ctx, QueueSelect, self, location[0])
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}

	cacheMap[self] = Cache{
		queue:      ids,
		lastAccess: time.Now(),
		location:   location,
		currIndex:  0,
	}
	return nil
}
