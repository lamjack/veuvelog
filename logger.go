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
package log

import (
	"bytes"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

var (
	sequenceNo uint64
)

type Level uint8

func (l Level) String() string {
	names := [...]string{
		"CRITICAL",
		"ERROR",
		"WARNING",
		"NOTICE",
		"INFO",
		"DEBUG",
	}
	if l < 0 || l > 5 {
		return "UNKNOWN"
	}
	return names[l]
}

// Log levels
const (
	CRITICAL Level = iota
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

type Record struct {
	ID        uint64
	Time      time.Time
	Module    string
	Level     Level
	LevelName string
	Args      []interface{}

	message []byte
	fmt     *string
}

func (r *Record) Message() string {
	if r.message == nil {
		var buf bytes.Buffer
		if r.fmt != nil {
			_, _ = fmt.Fprintf(&buf, *r.fmt, r.Args...)
		} else {
			// use Fprintln to make sure we always get space between arguments
			_, _ = fmt.Fprintln(&buf, r.Args...)
			buf.Truncate(buf.Len() - 1) // strip newline
		}
		r.message = buf.Bytes()
	}
	return string(r.message)
}

type Logger struct {
	Name string

	level    Level
	handlers []Handler // mock a stack with slice

	closed bool
}

func NewLogger(name string) *Logger {
	return &Logger{Name: name}
}

func (l *Logger) PushHandler(h Handler) {
	l.handlers = append([]Handler{h}, l.handlers...)
}

func (l *Logger) PopHandler() Handler {
	if len(l.handlers) == 0 {
		panic("You tried to pop from an empty handler stack.")
	}
	var h Handler
	h, l.handlers = l.handlers[len(l.handlers)-1], l.handlers[:len(l.handlers)-1]
	return h
}

func (l *Logger) SetHandlers(handlers []Handler) {
	l.handlers = handlers
}

func (l *Logger) GetHandlers() []Handler {
	return l.handlers
}

func (l *Logger) Close() {
	if l.closed {
		return
	}
	l.closed = true
	for _, h := range l.handlers {
		_ = h.Close()
	}
}
func (l *Logger) log(lvl Level, format *string, args ...interface{}) {
	record := &Record{
		ID:        atomic.AddUint64(&sequenceNo, 1),
		Time:      time.Now(),
		Module:    l.Name,
		Level:     lvl,
		LevelName: lvl.String(),
		Args:      args,
		fmt:       format,
	}
	for _, h := range l.handlers {
		if h.GetLevel() >= lvl {
			// @todo 下一个handlers不是formatter handler的情况下,message要重置
			if fh, ok := h.(FormattableHandler); ok {
				record.message = fh.GetFormatter().format(*record)
			}
			h.Handle(*record)
		}
	}
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log(CRITICAL, nil, args...)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(CRITICAL, &format, args...)
	os.Exit(1)
}

func (l *Logger) Panic(args ...interface{}) {
	l.log(CRITICAL, nil, args...)
	panic(fmt.Sprint(args...))
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.log(CRITICAL, &format, args...)
	panic(fmt.Sprintf(format, args...))
}

func (l *Logger) Critical(args ...interface{}) {
	l.log(CRITICAL, nil, args...)
}

func (l *Logger) Criticalf(format string, args ...interface{}) {
	l.log(CRITICAL, &format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log(ERROR, nil, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(ERROR, &format, args...)
}

func (l *Logger) Warning(args ...interface{}) {
	l.log(WARNING, nil, args...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.log(WARNING, &format, args...)
}

func (l *Logger) Notice(args ...interface{}) {
	l.log(NOTICE, nil, args...)
}

func (l *Logger) Noticef(format string, args ...interface{}) {
	l.log(NOTICE, &format, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log(INFO, nil, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(INFO, &format, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.log(DEBUG, nil, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DEBUG, &format, args...)
}
