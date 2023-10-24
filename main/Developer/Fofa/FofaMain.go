// Copyright 2016 The FOFA SDK Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Package examples shows how to use fofa sdk

package Fofa

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type QueryResponse struct {
	Error           bool            `json:"error"`
	ConsumedFPoint  int             `json:"consumed_fpoint"`
	RequiredFPoints int             `json:"required_fpoints"`
	Size            int             `json:"size"`
	Page            int             `json:"page"`
	Mode            string          `json:"mode"`
	Query           string          `json:"query"`
	Results         [][]interface{} `json:"results"`
}

// Fofa a fofa client can be used to make queries
type Fofa struct {
	email []byte
	key   []byte
	*http.Client
}

// Result represents a record of the query results
// contain domain host  ip  port title country city
type result struct {
	Domain  string `json:"domain,omitempty"`
	Host    string `json:"host,omitempty"`
	IP      string `json:"ip,omitempty"`
	Port    string `json:"port,omitempty"`
	Title   string `json:"title,omitempty"`
	Country string `json:"country,omitempty"`
	City    string `json:"city,omitempty"`
}

// User struct for fofa user
type User struct {
	Email  string `json:"email,omitempty"`
	Fcoin  int    `json:"fcoin,omitempty"`
	Vip    bool   `json:"bool,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Err    string `json:"errmsg,omitempty"`
}

// Results fofa result set
type Results []result

const (
	defaultAPIUrl = "https://fofa.info/api/v1/search/all?"
)

var (
	errFofaReplyWrongFormat = errors.New("Fofa Reply With Wrong Format")
	errFofaReplyNoData      = errors.New("No Data In Fofa Reply")
)

// Query fofa sdk functions included
func Query(email, key, fofaQuery string, page uint) []byte {
	clt := NewFofaClient([]byte(email), []byte(key))
	if clt == nil {
		fmt.Printf("fofa query error!\n")
		os.Exit(1)
	}
	ret, err := clt.QueryAsJSON(page, []byte(fofaQuery))
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(1)
	}
	return ret
}

// NewFofaClient create a fofa client
func NewFofaClient(email, key []byte) *Fofa {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Fofa{
		email: email,
		key:   key,
		Client: &http.Client{
			Transport: transCfg, // disable tls verify
		},
	}
}

// Get overwrite http.Get
func (ff *Fofa) Get(u string) ([]byte, error) {

	body, err := ff.Client.Get(u)
	if err != nil {
		return nil, err
	}
	defer body.Body.Close()
	content, err := ioutil.ReadAll(body.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// QueryAsJSON make a fofa query and return json data as result
// echo 'domain="nosec.org"' | base64 - | xargs -I{}
// curl "https://fofa.so/api/v1/search/all?email=${FOFA_EMAIL}&key=${FOFA_KEY}&qbase64={}"
// host title ip domain port country city
func (ff *Fofa) QueryAsJSON(page uint, args ...[]byte) ([]byte, error) {
	var (
		query  = []byte(nil)
		fields = []byte("protocol,host,ip")
		q      = []byte(nil)
	)
	switch {
	case len(args) == 1 || (len(args) == 2 && args[1] == nil):
		query = args[0]
	case len(args) == 2:
		query = args[0]
		fields = args[1]
	}

	q = []byte(base64.StdEncoding.EncodeToString(query))
	q = bytes.Join([][]byte{[]byte(defaultAPIUrl),
		[]byte("email="), ff.email,
		[]byte("&key="), ff.key,
		[]byte("&qbase64="), q,
		[]byte("&fields="), fields,
		[]byte("&page="), []byte(strconv.Itoa(int(page))),
		[]byte("&size="), []byte(strconv.Itoa(500)),
	}, []byte(""))
	//fmt.Printf("%s\n", q)
	content, err := ff.Get(string(q))
	if err != nil {
		return nil, err
	}
	errmsg, err := jsonparser.GetString(content, "errmsg")
	if err == nil {
		err = errors.New(errmsg)
	} else {
		err = nil
	}
	return content, err
}

// QueryAsArray make a fofa query and
// return array data as result
// echo 'domain="nosec.org"' | base64 - | xargs -I{}
// curl "https://fofa.so/api/v1/search/all?email=${FOFA_EMAIL}&key=${FOFA_KEY}&qbase64={}"
func (ff *Fofa) QueryAsArray(page uint, args ...[]byte) (result Results, err error) {

	var content []byte

	content, err = ff.QueryAsJSON(page, args...)
	if err != nil {
		return nil, err
	}

	errmsg, err := jsonparser.GetString(content, "errmsg")
	// err equals to nil on error
	if err == nil {
		return nil, errors.New(errmsg)
	}

	err = json.Unmarshal(content, &result)

	return
}

// UserInfo get user information
func (ff *Fofa) UserInfo() (user *User, err error) {
	user = new(User)
	queryStr := strings.Join([]string{"https://fofa.so/api/v1/info/my?email=", string(ff.email), "&key=", string(ff.key)}, "")

	content, err := ff.Get(queryStr)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(content, user); err != nil {
		return nil, err
	}

	if len(user.Err) != 0 {
		return nil, errors.New(user.Err)
	}

	return user, nil
}

func (u *User) String() string {
	data, err := json.Marshal(u)
	if err != nil {
		log.Fatalf("json marshal failed. err: %s\n", err)
		return ""
	}
	return string(data)
}

func (r *result) String() string {
	data, err := json.Marshal(r)
	if err != nil {
		log.Fatalf("json marshal failed. err: %s\n", err)
		return ""
	}
	return string(data)
}

func (r *Results) String() string {
	data, err := json.Marshal(r)
	if err != nil {
		log.Fatalf("json marshal failed. err: %s\n", err)
		return ""
	}
	return string(data)
}
