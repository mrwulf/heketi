//
// Copyright (c) 2015 The heketi Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package glusterfs

import (
	"bytes"
	"github.com/heketi/tests"
	"github.com/heketi/utils"
	"os"
	"testing"
)

func TestAppBadConfigData(t *testing.T) {
	data := []byte(`{ bad json }`)
	app := NewApp(bytes.NewBuffer(data))
	tests.Assert(t, app == nil)

	data = []byte(`{}`)
	app = NewApp(bytes.NewReader(data))
	tests.Assert(t, app == nil)

	data = []byte(`{
		"glusterfs" : {}
		}`)
	app = NewApp(bytes.NewReader(data))
	tests.Assert(t, app == nil)
}

func TestAppUnknownExecutorInConfig(t *testing.T) {
	data := []byte(`{
		"glusterfs" : {
			"executor" : "unknown value here"
		}
		}`)
	app := NewApp(bytes.NewReader(data))
	tests.Assert(t, app == nil)
}

func TestAppUnknownAllocatorInConfig(t *testing.T) {
	data := []byte(`{
		"glusterfs" : {
			"allocator" : "unknown value here"
		}
		}`)
	app := NewApp(bytes.NewReader(data))
	tests.Assert(t, app == nil)
}

func TestAppBadDbLocation(t *testing.T) {
	data := []byte(`{
		"glusterfs" : {
			"db" : "/badlocation"
		}
	}`)
	app := NewApp(bytes.NewReader(data))
	tests.Assert(t, app == nil)
}

func TestAppAdvsettings(t *testing.T) {

	dbfile := tests.Tempfile()
	defer os.Remove(dbfile)
	os.Setenv("HEKETI_EXECUTOR", "mock")
	defer os.Unsetenv("HEKETI_EXECUTOR")

	data := []byte(`{
		"glusterfs" : {
			"executor" : "crazyexec",
			"allocator" : "simple",
			"db" : "` + dbfile + `",
			"brick_max_size_gb" : 1024,
			"brick_min_size_gb" : 1,
			"max_bricks_per_volume" : 33
		}
	}`)

	bmax, bmin, bnum := BrickMaxSize, BrickMinSize, BrickMaxNum
	defer func() {
		BrickMaxSize, BrickMinSize, BrickMaxNum = bmax, bmin, bnum
	}()

	app := NewApp(bytes.NewReader(data))
	tests.Assert(t, app != nil)
	tests.Assert(t, app.conf.Executor == "mock")
	tests.Assert(t, BrickMaxNum == 33)
	tests.Assert(t, BrickMaxSize == 1*TB)
	tests.Assert(t, BrickMinSize == 1*GB)
}

func TestAppLogLevel(t *testing.T) {
	dbfile := tests.Tempfile()
	defer os.Remove(dbfile)

	levels := []string{
		"none",
		"critical",
		"error",
		"warning",
		"info",
		"debug",
	}

	logger.SetLevel(utils.LEVEL_DEBUG)
	for _, level := range levels {
		data := []byte(`{
			"glusterfs" : {
				"executor" : "mock",
				"allocator" : "simple",
				"db" : "` + dbfile + `",
				"loglevel" : "` + level + `"
			}
		}`)

		app := NewApp(bytes.NewReader(data))
		tests.Assert(t, app != nil, level, string(data))

		switch level {
		case "none":
			tests.Assert(t, logger.Level() == utils.LEVEL_NOLOG)
		case "critical":
			tests.Assert(t, logger.Level() == utils.LEVEL_CRITICAL)
		case "error":
			tests.Assert(t, logger.Level() == utils.LEVEL_ERROR)
		case "warning":
			tests.Assert(t, logger.Level() == utils.LEVEL_WARNING)
		case "info":
			tests.Assert(t, logger.Level() == utils.LEVEL_INFO)
		case "debug":
			tests.Assert(t, logger.Level() == utils.LEVEL_DEBUG)
		}
		app.Close()
	}

	// Test that an unknown value does not change the loglevel
	logger.SetLevel(utils.LEVEL_NOLOG)
	data := []byte(`{
			"glusterfs" : {
				"executor" : "mock",
				"allocator" : "simple",
				"db" : "` + dbfile + `",
				"loglevel" : "blah"
			}
		}`)

	app := NewApp(bytes.NewReader(data))
	tests.Assert(t, app != nil)
	tests.Assert(t, logger.Level() == utils.LEVEL_NOLOG)
}
