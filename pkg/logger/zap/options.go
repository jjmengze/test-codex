// Copyright (c) 2023 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

//go:build go1.21

package handler

import "log/slog"

// handlerOption configures a slog Handler.
type handlerOption interface {
	apply(*Handler)
}

// HandlerOptionFunc wraps a func so it satisfies the Option interface.
type HandlerOptionFunc func(*Handler)

func (f HandlerOptionFunc) apply(handler *Handler) {
	f(handler)
}

// withCaller configures the Logger to include the filename and line number
// of the caller in log messages--if available.
func withCaller(enabled bool) handlerOption {
	return HandlerOptionFunc(func(handler *Handler) {
		handler.addCaller = enabled
	})
}

// withStacktrace configures the Logger to record a stack trace
// for all messages at or above a given level.
func withStacktrace(lvl slog.Level) handlerOption {
	return HandlerOptionFunc(func(log *Handler) {
		log.addStackAt = lvl
	})
}
