// Package videotest provides dummy video driver for testing.
package videotest

import (
	"context"
	"image"
	"io"
	"time"

	"github.com/pion/mediadevices/pkg/driver"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/io/video"
	"github.com/pion/mediadevices/pkg/prop"
)

func init() {
	driver.GetManager().Register(
		newVideoTest(),
		driver.Info{Label: "VideoTest", DeviceType: driver.Camera},
	)
}

type dummy struct {
	closed <-chan struct{}
	cancel func()
	tick   *time.Ticker
}

func newVideoTest() *dummy {
	return &dummy{}
}

func (d *dummy) Open() error {
	ctx, cancel := context.WithCancel(context.Background())
	d.closed = ctx.Done()
	d.cancel = cancel
	return nil
}

func (d *dummy) Close() error {
	d.cancel()
	if d.tick != nil {
		d.tick.Stop()
	}
	return nil
}

func (d *dummy) VideoRecord(p prop.Media) (video.Reader, error) {
	colors := [][3]byte{
		{235, 128, 128},
		{210, 16, 146},
		{170, 166, 16},
		{145, 54, 34},
		{107, 202, 222},
		{82, 90, 240},
		{41, 240, 110},
	}

	// numbers := [][12]bool{
	// 	{false, false, false, true, false, false, false, true, false, false, false, true},
	// 	{false, true, true, false, false, false, true, false, false, false, true, true},
	// }

	yi := p.Width * p.Height
	ci := yi / 2
	yy := make([]byte, yi)
	cb := make([]byte, ci)
	cr := make([]byte, ci)
	yyBase := make([]byte, yi)
	cbBase := make([]byte, ci)
	crBase := make([]byte, ci)
	hColorBarEnd := p.Height * 3 / 4
	wGradationEnd := p.Width
	for y := 0; y < hColorBarEnd; y++ {
		yi := p.Width * y
		ci := p.Width * y / 2
		// Color bar
		for x := 0; x < p.Width; x++ {
			c := x * 7 / p.Width
			yyBase[yi+x] = uint8(uint16(colors[c][0]) * 75 / 100)
			cbBase[ci+x/2] = colors[c][1]
			crBase[ci+x/2] = colors[c][2]
		}
	}
	for y := hColorBarEnd; y < p.Height; y++ {
		yi := p.Width * y
		ci := p.Width * y / 2
		for x := 0; x < wGradationEnd; x++ {
			// Gray gradation
			yyBase[yi+x] = uint8(x * 255 / wGradationEnd)
			cbBase[ci+x/2] = 128
			crBase[ci+x/2] = 128
		}
	}

	tick := time.NewTicker(time.Duration(float32(time.Second) / p.FrameRate))
	d.tick = tick
	closed := d.closed

	r := video.ReaderFunc(func() (image.Image, func(), error) {
		select {
		case <-closed:
			return nil, func() {}, io.EOF
		default:
		}

		<-tick.C

		copy(yy, yyBase)
		copy(cb, cbBase)
		copy(cr, crBase)

		// write the time
		now := time.Now()
		// format: sssmm (s - second, m - millisecond)
		timestamp := int(now.UnixMilli() / 10 % 100000)
		numbers := [5]int{
			timestamp / 10000,
			timestamp / 1000 % 10,
			timestamp / 100 % 10,
			timestamp / 10 % 10,
			timestamp % 10,
		}

		for y := 0; y < 5; y++ {
			yi := p.Width * (y * 5 + 16)
			for x := 0; x < numbers[y]; x++ {
				xi := x * 5 + 16

				yy[yi+xi] = byte(x % 3 * 60)
				yy[yi+xi+1] = byte(x % 3 * 60)
				yy[yi+xi+2] = byte(x % 3 * 60)
				yy[yi+xi+p.Width] = byte(x % 3 * 60)
				yy[yi+xi+p.Width+1] = byte(x % 3 * 60)
				yy[yi+xi+p.Width+2] = byte(x % 3 * 60)
				yy[yi+xi+p.Width+p.Width] = byte(x % 3 * 60)
				yy[yi+xi+p.Width+p.Width+1] = byte(x % 3 * 60)
				yy[yi+xi+p.Width+p.Width+2] = byte(x % 3 * 60)
			}
		}

		return &image.YCbCr{
			Y:              yy,
			YStride:        p.Width,
			Cb:             cb,
			Cr:             cr,
			CStride:        p.Width / 2,
			SubsampleRatio: image.YCbCrSubsampleRatio422,
			Rect:           image.Rect(0, 0, p.Width, p.Height),
		}, func() {}, nil
	})

	return r, nil
}

func (d dummy) Properties() []prop.Media {
	return []prop.Media{
		{
			Video: prop.Video{
				Width:       640,
				Height:      480,
				FrameFormat: frame.FormatYUYV,
				FrameRate:   30,
			},
		},
		{
			Video: prop.Video{
				Width:       1280,
				Height:      720,
				FrameFormat: frame.FormatYUYV,
				FrameRate:   30,
			},
		},
		{
			Video: prop.Video{
				Width:       1280,
				Height:      720,
				FrameFormat: frame.FormatYUYV,
				FrameRate:   60,
			},
		},
		{
			Video: prop.Video{
				Width:       1920,
				Height:      1080,
				FrameFormat: frame.FormatYUYV,
				FrameRate:   30,
			},
		},
		{
			Video: prop.Video{
				Width:       1920,
				Height:      1080,
				FrameFormat: frame.FormatYUYV,
				FrameRate:   60,
			},
		},
		{
			Video: prop.Video{
				Width:       3840,
				Height:      2160,
				FrameFormat: frame.FormatYUYV,
				FrameRate:   30,
			},
		},
		{
			Video: prop.Video{
				Width:       3840,
				Height:      2160,
				FrameFormat: frame.FormatYUYV,
				FrameRate:   60,
			},
		},
	}
}