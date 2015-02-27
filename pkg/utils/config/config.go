/*
 * Mini Object Storage, (C) 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"sync"

	"github.com/minio-io/minio/pkg/utils/helpers"
)

type Config struct {
	configPath string
	configFile string
	configLock *sync.RWMutex
	Users      map[string]User
}

type User struct {
	Name      string
	AccessKey string
	SecretKey string
}

// Initialize config directory and template config
func (c *Config) SetupConfig() error {
	confPath := path.Join(helpers.HomeDir(), ".minio")
	if err := os.MkdirAll(confPath, os.ModeDir); err != nil {
		return err
	}

	c.configPath = confPath
	c.configFile = path.Join(c.configPath, "config.json")
	if _, err := os.Stat(c.configFile); os.IsNotExist(err) {
		_, err = os.Create(c.configFile)
		if err != nil {
			return err
		}
	}

	c.configLock = new(sync.RWMutex)
	return nil
}

// Get config file location
func (c *Config) GetConfigPath() string {
	return c.configPath
}

// Verify if user exists
func (c *Config) IsUserExists(username string) bool {
	for _, user := range c.Users {
		if user.Name == username {
			return true
		}
	}
	return false
}

// Get user based on accesskey
func (c *Config) GetKey(accessKey string) User {
	value, ok := c.Users[accessKey]
	if !ok {
		return User{}
	}
	return value
}

// Get user based on username
func (c *Config) GetUser(username string) User {
	for _, user := range c.Users {
		if user.Name == username {
			return user
		}
	}
	return User{}
}

// Add a new user into existing User list
func (c *Config) AddUser(user User) {
	var currentUsers map[string]User
	if len(c.Users) == 0 {
		currentUsers = make(map[string]User)
	} else {
		currentUsers = c.Users
	}
	currentUsers[user.AccessKey] = user
	c.Users = currentUsers
}

// Write encoded json in config file
func (c *Config) WriteConfig() error {
	c.configLock.Lock()
	defer c.configLock.Unlock()

	var file *os.File
	var err error

	file, err = os.OpenFile(c.configFile, os.O_WRONLY, 0666)
	defer file.Close()
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.Encode(c.Users)
	return nil
}

// Read json config file and decode
func (c *Config) ReadConfig() error {
	c.configLock.RLock()
	defer c.configLock.RUnlock()

	var file *os.File
	var err error

	file, err = os.OpenFile(c.configFile, os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		return err
	}

	users := make(map[string]User)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&users)
	switch err {
	case io.EOF:
		return nil
	case nil:
		c.Users = users
		return nil
	default:
		return err
	}
}

/// helpers

// Load all users into memory
func Loadusers() map[string]User {
	c := Config{}
	c.SetupConfig()
	c.ReadConfig()
	return c.Users
}

// Load a given user based on accessKey
func Loadkey(accessKeyId string) User {
	c := Config{}
	c.SetupConfig()
	c.ReadConfig()
	return c.GetKey(accessKeyId)
}

// Load a given user based on username
func Loaduser(username string) User {
	c := Config{}
	c.SetupConfig()
	c.ReadConfig()
	return c.GetUser(username)
}