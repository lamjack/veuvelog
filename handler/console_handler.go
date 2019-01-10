// Copyright 2019 The WIZ Technology Co. Ltd. Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package handler

import (
	"fmt"
	"github.com/lamjack/veuvelog"
	"runtime"
)

var (
	yellow = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	red    = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue   = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	reset  = string([]byte{27, 91, 48, 109})
)

type ConsoleHandler struct {
	level        log.Level
	disableColor bool
}

func NewConsoleHandler(lvl log.Level) *ConsoleHandler {
	h := new(ConsoleHandler)
	h.level = lvl
	if runtime.GOOS == "windows" {
		h.disableColor = true
	} else {
		h.disableColor = false
	}
	return h
}

func (h *ConsoleHandler) GetLevel() log.Level {
	return h.level
}

func (h *ConsoleHandler) Handle(record log.Record) {
	lvlText := h.getLevelColor(record.Level)
	fmt.Printf(lvlText + record.Message() + "\n")
}

func (ConsoleHandler) Close() error {
	return nil
}

func (h ConsoleHandler) getLevelColor(lvl log.Level) string {
	var r string
	if !h.disableColor {
		if lvl == log.DEBUG || lvl == log.INFO || lvl == log.NOTICE {
			r = "[" + fmt.Sprintf(blue+lvl.String()+reset) + "]"
		} else if lvl == log.WARNING {
			r = "[" + fmt.Sprintf(yellow+lvl.String()+reset) + "]"
		} else if lvl == log.ERROR || lvl == log.CRITICAL {
			r = "[" + fmt.Sprintf(red+lvl.String()+reset) + "]"
		}
	}
	if r == "" {
		r = fmt.Sprintf("[%s]", lvl.String())
	}
	return r
}
