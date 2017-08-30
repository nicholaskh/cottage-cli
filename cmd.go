package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/nicholaskh/log4go"
)

type cmd struct {
	command string
	params  []string
}

func newCmd() *cmd {
	this := new(cmd)
	if flag.NArg() == 0 {
		flag.Usage()
	}
	this.command = flag.Arg(0)
	this.params = make([]string, 0)
	for i := 1; i < flag.NArg(); i++ {
		this.params = append(this.params, flag.Arg(i))
	}

	return this
}

func (this *cmd) run() {
	switch this.command {
	case "search":
		if len(this.params) == 0 {
			panic("No search word specified")
		}
		word := this.params[0]
		this.search(word)

	case "install":
		if len(this.params) == 0 {
			panic("No package name specified")
		}
		name := this.params[0]
		this.install(name)
	}
}

func (this *cmd) search(word string) {
	resp, err := http.PostForm("http://127.0.0.1:8844/search", url.Values{"word": {word}})
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	log.Info("============================================================== matched: %s ==============================================================", word)
	for _, name := range this.parseData(body).([]interface{}) {
		log.Info(name)
	}

	return
}

func (this *cmd) install(name string) {
	resp, err := http.PostForm("http://127.0.0.1:8844/max-version", url.Values{"name": {name}})
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	maxVersion := this.parseData(body).(string)

	resp, err = http.PostForm("http://127.0.0.1:8844/download", url.Values{"name": {name}, "version": {maxVersion}})
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	filename := fmt.Sprintf("%s.jar", name)
	err = ioutil.WriteFile(filename, body, 0666)
	if err != nil {
		log.Info("Write file error: %s", err.Error())
	}
	log.Info("%s has installed", filename)
}

func (this *cmd) parseData(body []byte) interface{} {
	var bodyJson map[string]interface{}
	json.Unmarshal(body, &bodyJson)
	return bodyJson["data"]
}
