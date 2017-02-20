// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package matchmaker

import (
	"errors"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	"k8s.io/client-go/pkg/util/uuid"
)

const (
	// gameStatusOpen means the game is open, and available to be joined
	gameStatusOpen = 0
	// gameStatusClosed means the game is closed, and currently being played
	gameStatusClosed = 1

	// redisOpenGameListKey is the redis key for
	// where the list of open games is stored
	redisOpenGameListKey = "openGameList"

	// redisGamePrefix is the prefix for the key
	// for where game session data is stored in redis
	redisGamePrefix = "game:"
)

var (
	// errGameNotFound error returned when a game
	// can't be found
	errGameNotFound = errors.New("Game not found")
)

// Game represents a game that is being/has been match-made
type Game struct {
	ID        string `json:"id" redis:"id"`
	Status    int    `json:"status" redis:"status"`
	SessionID string `json:"sessionID,omitempty" redis:"sessionID"`
	Port      int    `json:"port,omitempty" redis:"port"`
	IP        string `json:"ip,omitempty" redis:"ip"`
}

// NewGame returns a game, with a unique id
func NewGame() *Game {
	return &Game{
		Status: gameStatusOpen,
		ID:     string(uuid.NewUUID()),
	}
}

// Key returns the redis key for this Game
func (g Game) Key() string {
	return redisGamePrefix + g.ID
}

// updateGame updates the game data in redis
// with the SessionID, Port, IP and status
func updateGame(con redis.Conn, g *Game) error {
	_, err := con.Do("HMSET", g.Key(), "status", g.Status, "sessionID", g.SessionID, "port", g.Port, "ip", g.IP)

	if err != nil {
		log.Printf("[Error][game] Error updating game: %#v, %v", *g, err)
	}

	return err
}

// getGame retrieves a game from redis, and then returns it
func getGame(con redis.Conn, key string) (*Game, error) {
	var g *Game
	values, err := redis.Values(con.Do("HGETALL", key))

	if err != nil {
		log.Printf("[Error][games] Error getting hash for key %v. %v", key, err)
		return g, err
	}

	if len(values) == 0 {
		log.Printf("[Error][games] Could not find record for key: %v", key)
		return g, fmt.Errorf("Could not find game for key: %v", key)
	}

	g = &Game{}
	err = redis.ScanStruct(values, g)

	if err != nil {
		log.Printf("[Error][games] Error scanning struct: %v", err)
	}

	return g, err
}

// pushOpenGame pushes an open game onto the list of open games
func pushOpenGame(con redis.Conn, g *Game) error {
	key := g.Key()
	log.Printf("[Info][game] Pushing game onto open list: %v", key)

	err := con.Send("MULTI")
	if err != nil {
		log.Printf("[Error][games] Could not Send MULTI: %v", err)
		return err
	}

	err = con.Send("RPUSH", redisOpenGameListKey, key)
	if err != nil {
		log.Printf("[Error][games] Could not Send RPUSH: %v", err)
		return err
	}

	err = con.Send("HMSET", key, "id", g.ID, "status", g.Status)
	if err != nil {
		log.Printf("[Error][games] Could not Send HMSET: %v", err)
		return err
	}

	err = con.Send("EXPIRE", key, 60*60)
	if err != nil {
		log.Printf("[Error][games] Could not Send EXPIRE: %v", err)
		return err
	}

	_, err = con.Do("EXEC")
	if err != nil {
		log.Printf("[Error][games] Could not save session to redis: %v", err)
		return err
	}

	return nil
}

// popOpenGame pops an open game off the list, and returns it's data structure
func popOpenGame(con redis.Conn) (*Game, error) {
	log.Print("[Info][game] Attempting to pop an open game")
	key, err := redis.String(con.Do("LPOP", redisOpenGameListKey))
	if err == redis.ErrNil {
		log.Print("[Info][game] Game not found, returning")
		return nil, errGameNotFound
	}

	log.Print("[Info][game] Found game, decoding...")

	g, err := getGame(con, key)

	log.Printf("[Info][game] Returning Game: %#v", g)
	return g, err
}
