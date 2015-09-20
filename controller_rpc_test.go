// +build ignore

/*
 * Minio Cloud Storage, (C) 2014 Minio, Inc.
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

package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gorilla/rpc/v2/json"
	"github.com/minio/minio/pkg/auth"
	. "gopkg.in/check.v1"
)

type ControllerRPCSuite struct{}

var _ = Suite(&ControllerRPCSuite{})

var testRPCServer *httptest.Server

func (s *ControllerRPCSuite) SetUpSuite(c *C) {
	root, err := ioutil.TempDir(os.TempDir(), "api-")
	c.Assert(err, IsNil)
	auth.SetAuthConfigPath(root)

	testRPCServer = httptest.NewServer(getRPCCtrlHandler())
}

func (s *ControllerRPCSuite) TearDownSuite(c *C) {
	testRPCServer.Close()
}

func (s *ControllerRPCSuite) TestMemStats(c *C) {
	op := rpcOperation{
		Method:  "Server.MemStats",
		Request: ServerArg{},
	}
	req, err := newRPCRequest(testRPCServer.URL+"/rpc", op, http.DefaultTransport)
	c.Assert(err, IsNil)
	c.Assert(req.Get("Content-Type"), Equals, "application/json")
	resp, err := req.Do()
	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, http.StatusOK)

	var reply MemStatsRep
	c.Assert(json.DecodeClientResponse(resp.Body, &reply), IsNil)
	resp.Body.Close()
	c.Assert(reply, Not(DeepEquals), MemStatsRep{})
}

func (s *ControllerRPCSuite) TestSysInfo(c *C) {
	op := rpcOperation{
		Method:  "Server.SysInfo",
		Request: ServerArg{},
	}
	req, err := newRPCRequest(testRPCServer.URL+"/rpc", op, http.DefaultTransport)
	c.Assert(err, IsNil)
	c.Assert(req.Get("Content-Type"), Equals, "application/json")
	resp, err := req.Do()
	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, http.StatusOK)

	var reply SysInfoRep
	c.Assert(json.DecodeClientResponse(resp.Body, &reply), IsNil)
	resp.Body.Close()
	c.Assert(reply, Not(DeepEquals), SysInfoRep{})
}

func (s *ControllerRPCSuite) TestServerList(c *C) {
	op := rpcOperation{
		Method:  "Server.List",
		Request: ServerArg{},
	}
	req, err := newRPCRequest(testRPCServer.URL+"/rpc", op, http.DefaultTransport)
	c.Assert(err, IsNil)
	c.Assert(req.Get("Content-Type"), Equals, "application/json")
	resp, err := req.Do()
	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, http.StatusOK)

	var reply ServerListRep
	c.Assert(json.DecodeClientResponse(resp.Body, &reply), IsNil)
	resp.Body.Close()
	c.Assert(reply, Not(DeepEquals), ServerListRep{})
}

func (s *ControllerRPCSuite) TestServerAdd(c *C) {
	op := rpcOperation{
		Method:  "Server.Add",
		Request: ServerArg{},
	}
	req, err := newRPCRequest(testRPCServer.URL+"/rpc", op, http.DefaultTransport)
	c.Assert(err, IsNil)
	c.Assert(req.Get("Content-Type"), Equals, "application/json")
	resp, err := req.Do()
	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, http.StatusOK)

	var reply DefaultRep
	c.Assert(json.DecodeClientResponse(resp.Body, &reply), IsNil)
	resp.Body.Close()
	c.Assert(reply, Not(DeepEquals), DefaultRep{0, "Added"})
}

func (s *ControllerRPCSuite) TestAuth(c *C) {
	op := rpcOperation{
		Method:  "Auth.Generate",
		Request: AuthArgs{User: "newuser"},
	}
	req, err := newRPCRequest(testRPCServer.URL+"/rpc", op, http.DefaultTransport)
	c.Assert(err, IsNil)
	c.Assert(req.Get("Content-Type"), Equals, "application/json")
	resp, err := req.Do()
	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, http.StatusOK)

	var reply AuthReply
	c.Assert(json.DecodeClientResponse(resp.Body, &reply), IsNil)
	resp.Body.Close()
	c.Assert(reply, Not(DeepEquals), AuthReply{})
	c.Assert(len(reply.AccessKeyID), Equals, 20)
	c.Assert(len(reply.SecretAccessKey), Equals, 40)
	c.Assert(len(reply.Name), Not(Equals), 0)

	op = rpcOperation{
		Method:  "Auth.Fetch",
		Request: AuthArgs{User: "newuser"},
	}
	req, err = newRPCRequest(testRPCServer.URL+"/rpc", op, http.DefaultTransport)
	c.Assert(err, IsNil)
	c.Assert(req.Get("Content-Type"), Equals, "application/json")
	resp, err = req.Do()
	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, http.StatusOK)

	var newReply AuthReply
	c.Assert(json.DecodeClientResponse(resp.Body, &newReply), IsNil)
	resp.Body.Close()
	c.Assert(newReply, Not(DeepEquals), AuthReply{})
	c.Assert(reply.AccessKeyID, Equals, newReply.AccessKeyID)
	c.Assert(reply.SecretAccessKey, Equals, newReply.SecretAccessKey)
	c.Assert(len(reply.Name), Not(Equals), 0)

	op = rpcOperation{
		Method:  "Auth.Reset",
		Request: AuthArgs{User: "newuser"},
	}
	req, err = newRPCRequest(testRPCServer.URL+"/rpc", op, http.DefaultTransport)
	c.Assert(err, IsNil)
	c.Assert(req.Get("Content-Type"), Equals, "application/json")
	resp, err = req.Do()
	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, http.StatusOK)

	var resetReply AuthReply
	c.Assert(json.DecodeClientResponse(resp.Body, &resetReply), IsNil)
	resp.Body.Close()
	c.Assert(newReply, Not(DeepEquals), AuthReply{})
	c.Assert(reply.AccessKeyID, Not(Equals), resetReply.AccessKeyID)
	c.Assert(reply.SecretAccessKey, Not(Equals), resetReply.SecretAccessKey)
	c.Assert(len(reply.Name), Not(Equals), 0)

	// these operations should fail

	/// generating access for existing user fails
	op = rpcOperation{
		Method:  "Auth.Generate",
		Request: AuthArgs{User: "newuser"},
	}
	req, err = newRPCRequest(testRPCServer.URL+"/rpc", op, http.DefaultTransport)
	c.Assert(err, IsNil)
	c.Assert(req.Get("Content-Type"), Equals, "application/json")
	resp, err = req.Do()
	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, http.StatusBadRequest)

	/// null user provided invalid
	op = rpcOperation{
		Method:  "Auth.Generate",
		Request: AuthArgs{User: ""},
	}
	req, err = newRPCRequest(testRPCServer.URL+"/rpc", op, http.DefaultTransport)
	c.Assert(err, IsNil)
	c.Assert(req.Get("Content-Type"), Equals, "application/json")
	resp, err = req.Do()
	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, http.StatusBadRequest)
}
