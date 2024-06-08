package gzip

import (
	"net/http"
	"regexp"
	"strings"
)

var (
	DefaultExcludedExtentions = NewExcludedExtensions([]string{
		".png", ".gif", ".jpeg", ".jpg",
	})
	DefaultOptions = &Options{
		ExcludedExtensions: DefaultExcludedExtentions,
	}
)

type Options struct {
	ExcludedExtensions ExcludedExtensions
	ExcludedPaths      ExcludedPaths
	ExcludedPathsRegex ExcludedPathsRegex
	Handler            http.HandlerFunc
}

type Option func(*Options)

func WithExcludedExtensions(args []string) Option {
	return func(o *Options) {
		o.ExcludedExtensions = NewExcludedExtensions(args)
	}
}

func WithExcludedPaths(args []string) Option {
	return func(o *Options) {
		o.ExcludedPaths = NewExcludedPaths(args)
	}
}

func WithExcludedPathsRegexes(args []string) Option {
	return func(o *Options) {
		o.ExcludedPathsRegex = NewExcludedPathsRegex(args)
	}
}

func WithHandler(decompressFn http.HandlerFunc) Option {
	return func(o *Options) {
		o.Handler = decompressFn
	}
}

// Using map for better lookup performance
type ExcludedExtensions map[string]bool

func NewExcludedExtensions(extensions []string) ExcludedExtensions {
	res := make(ExcludedExtensions)
	for _, e := range extensions {
		res[e] = true
	}
	return res
}

func (e ExcludedExtensions) Contains(target string) bool {
	_, ok := e[target]
	return ok
}

type ExcludedPaths []string

func NewExcludedPaths(paths []string) ExcludedPaths {
	return ExcludedPaths(paths)
}

func (e ExcludedPaths) Contains(requestURI string) bool {
	for _, path := range e {
		if strings.HasPrefix(requestURI, path) {
			return true
		}
	}
	return false
}

type ExcludedPathsRegex []*regexp.Regexp

func NewExcludedPathsRegex(regexes []string) ExcludedPathsRegex {
	result := make([]*regexp.Regexp, len(regexes), len(regexes))
	for i, reg := range regexes {
		result[i] = regexp.MustCompile(reg)
	}
	return result
}

func (e ExcludedPathsRegex) Contains(requestURI string) bool {
	for _, reg := range e {
		if reg.MatchString(requestURI) {
			return true
		}
	}
	return false
}
