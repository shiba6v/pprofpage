package controller

import (
	"flag"
	"strings"

	"github.com/google/pprof/driver"
	"github.com/labstack/echo/v4"
	"github.com/shiba6v/eu"
)

// https://github.com/kaz/pprotein/blob/master/internal/pprof/flagset.go
type (
	FlagSet struct {
		*flag.FlagSet

		input     []string
		usageMsgs []string
	}
)

func NewFlagSet(input []string) *FlagSet {
	return &FlagSet{
		flag.NewFlagSet("", flag.ContinueOnError),
		input,
		[]string{},
	}
}

func (f *FlagSet) StringList(o, d, c string) *[]*string {
	return &[]*string{f.String(o, d, c)}
}

func (f *FlagSet) ExtraUsage() string {
	return strings.Join(f.usageMsgs, "\n")
}
func (f *FlagSet) AddExtraUsage(eu string) {
	f.usageMsgs = append(f.usageMsgs, eu)
}

func (f *FlagSet) Parse(usage func()) []string {
	f.Usage = usage
	f.FlagSet.Parse(f.input)
	args := f.Args()
	if len(args) == 0 {
		usage()
	}
	return args
}

func (r Controller) GetPProf(c echo.Context) error {
	id := c.Param("id")
	key := c.Param("key")
	path, err := r.storage.GetObjectToTmp(c.Request().Context(), id)
	if err != nil {
		return eu.Wrap(err)
	}
	options := &driver.Options{
		Flagset: NewFlagSet(
			[]string{
				"-no_browser",
				"-http", "0:0",
				path,
			},
		),
		HTTPServer: func(args *driver.HTTPServerArgs) error {
			hs := args.Handlers
			for k, h := range hs {
				if k == "/"+key {
					h.ServeHTTP(c.Response().Writer, c.Request())
					break
				}
			}
			return nil
		},
	}
	if err := driver.PProf(options); err != nil {
		return eu.Wrap(err)
	}
	return nil
}
