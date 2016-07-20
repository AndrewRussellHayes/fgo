//
// command/util.go
//
// Copyright (c) 2016 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func generateFromAsset(assetName string, param interface{}) (res []byte, err error) {

	data, err := Asset(assetName)
	if err != nil {
		return
	}
	return generate(data, param)

}

func generateFromFile(path string, param interface{}) (res []byte, err error) {

	fp, err := os.Open(path)
	if err != nil {
		return
	}
	defer fp.Close()

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return
	}

	return generate(data, param)

}

func generate(data []byte, param interface{}) (res []byte, err error) {

	funcMap := template.FuncMap{
		"Title":   strings.Title,
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
	}

	temp, err := template.New("").Funcs(funcMap).Parse(string(data))
	if err != nil {
		return
	}

	buf := bytes.Buffer{}
	if err = temp.Execute(&buf, &param); err != nil {
		return
	}

	res = buf.Bytes()
	return

}

func Sha256(path string) (res string, err error) {

	fp, err := os.Open(path)
	if err != nil {
		return
	}
	defer fp.Close()

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return
	}

	bytes := sha256.Sum256(data)
	res = hex.EncodeToString(bytes[:])
	return

}
