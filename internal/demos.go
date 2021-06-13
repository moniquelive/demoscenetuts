package demos

import (
	"image"

	"github.com/moniquelive/demoscenetuts/internal/bifilter"
	"github.com/moniquelive/demoscenetuts/internal/bump"
	"github.com/moniquelive/demoscenetuts/internal/crossfade"
	"github.com/moniquelive/demoscenetuts/internal/cyber1"
	"github.com/moniquelive/demoscenetuts/internal/filters"
	"github.com/moniquelive/demoscenetuts/internal/mandelbrot"
	"github.com/moniquelive/demoscenetuts/internal/plasma"
	"github.com/moniquelive/demoscenetuts/internal/stars"
	"github.com/moniquelive/demoscenetuts/internal/stars3D"
	"github.com/moniquelive/demoscenetuts/internal/textmap"
)

type Demo interface {
	Draw(*image.RGBA)
	Setup() (int, int, int)
}

func FillDemos() map[string]Demo {
	m := make(map[string]Demo)
	m["stars"] = &stars.Stars{}
	m["3d"] = &stars3D.Stars{}
	m["crossfade"] = &crossfade.Cross{}
	m["plasma"] = &plasma.Plasma{}
	m["filter"] = &filters.Filter{}
	m["cyber1"] = &cyber1.Lerp{}
	m["bifilter"] = &bifilter.Bifilter{}
	m["bump"] = &bump.Bump{}
	m["mandelbrot"] = &mandelbrot.Mandelbrot{}
	m["textmap"] = &textmap.Textmap{}
	return m
}

