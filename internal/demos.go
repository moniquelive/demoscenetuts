package demos

import (
	"image"

	"github.com/moniquelive/demoscenetuts/internal/bifilter"
	"github.com/moniquelive/demoscenetuts/internal/bump"
	"github.com/moniquelive/demoscenetuts/internal/crossfade"
	"github.com/moniquelive/demoscenetuts/internal/cyber1"
	"github.com/moniquelive/demoscenetuts/internal/filters"
	"github.com/moniquelive/demoscenetuts/internal/mandelbrot"
	"github.com/moniquelive/demoscenetuts/internal/particles"
	"github.com/moniquelive/demoscenetuts/internal/plane"
	"github.com/moniquelive/demoscenetuts/internal/plasma"
	"github.com/moniquelive/demoscenetuts/internal/polygon"
	"github.com/moniquelive/demoscenetuts/internal/rotozoom"
	"github.com/moniquelive/demoscenetuts/internal/span"
	"github.com/moniquelive/demoscenetuts/internal/stars"
	"github.com/moniquelive/demoscenetuts/internal/stars3D"
	"github.com/moniquelive/demoscenetuts/internal/textmap"
)

type Demo interface {
	Draw(*image.RGBA)
	Setup() (int, int, int)
}

func FillDemos() map[string]Demo {
	return map[string]Demo{
		"stars":      &stars.Stars{},
		"3d":         &stars3D.Stars{},
		"crossfade":  &crossfade.Cross{},
		"plasma":     &plasma.Plasma{},
		"filter":     &filters.Filter{},
		"cyber1":     &cyber1.Lerp{},
		"bifilter":   &bifilter.Bifilter{},
		"bump":       &bump.Bump{},
		"mandelbrot": &mandelbrot.Mandelbrot{},
		"textmap":    &textmap.Textmap{},
		"rotozoom":   &rotozoom.Rotozoom{},
		"particles":  &particles.Particles{},
		"span":       &span.Span{},
		"polygon":    &polygon.Polygon{},
		"plane":      &plane.Plane{},
	}
}
