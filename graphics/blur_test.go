// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphics

import (
	"image"
	"testing"

	_ "image/png"
)

var blurOneColorTests = []transformOneColorTest{
	{
		"1x1-blank", 1, 1, 1, 1,
		&BlurOptions{0.83, 1},
		[]uint8{0xff},
		[]uint8{0xff},
	},
	{
		"1x1-spreadblank", 1, 1, 1, 1,
		&BlurOptions{0.83, 2},
		[]uint8{0xff},
		[]uint8{0xff},
	},
	{
		"3x3-blank", 3, 3, 3, 3,
		&BlurOptions{0.83, 2},
		[]uint8{
			0xff, 0xff, 0xff,
			0xff, 0xff, 0xff,
			0xff, 0xff, 0xff,
		},
		[]uint8{
			0xff, 0xff, 0xff,
			0xff, 0xff, 0xff,
			0xff, 0xff, 0xff,
		},
	},
	{
		"3x3-dot", 3, 3, 3, 3,
		&BlurOptions{0.34, 1},
		[]uint8{
			0x00, 0x00, 0x00,
			0x00, 0xff, 0x00,
			0x00, 0x00, 0x00,
		},
		[]uint8{
			0x00, 0x03, 0x00,
			0x03, 0xf2, 0x03,
			0x00, 0x03, 0x00,
		},
	},
	{
		"5x5-dot", 5, 5, 5, 5,
		&BlurOptions{0.34, 1},
		[]uint8{
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0xff, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
		},
		[]uint8{
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x03, 0x00, 0x00,
			0x00, 0x03, 0xf2, 0x03, 0x00,
			0x00, 0x00, 0x03, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
		},
	},
	{
		"5x5-dot-spread", 5, 5, 5, 5,
		&BlurOptions{0.85, 1},
		[]uint8{
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0xff, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
		},
		[]uint8{
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x10, 0x20, 0x10, 0x00,
			0x00, 0x20, 0x40, 0x20, 0x00,
			0x00, 0x10, 0x20, 0x10, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
		},
	},
	{
		"4x4-box", 4, 4, 4, 4,
		&BlurOptions{0.34, 1},
		[]uint8{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0xff, 0xff, 0x00,
			0x00, 0xff, 0xff, 0x00,
			0x00, 0x00, 0x00, 0x00,
		},
		[]uint8{
			0x00, 0x03, 0x03, 0x00,
			0x03, 0xf8, 0xf8, 0x03,
			0x03, 0xf8, 0xf8, 0x03,
			0x00, 0x03, 0x03, 0x00,
		},
	},
	{
		"5x5-twodots", 5, 5, 5, 5,
		&BlurOptions{0.34, 1},
		[]uint8{
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x96, 0x00, 0x96, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
		},
		[]uint8{
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x02, 0x00, 0x02, 0x00,
			0x02, 0x8e, 0x04, 0x8e, 0x02,
			0x00, 0x02, 0x00, 0x02, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
		},
	},
}

func TestBlurOneColor(t *testing.T) {
	for _, oc := range blurOneColorTests {
		dst := oc.newDst()
		src := oc.newSrc()
		opt := oc.opt.(*BlurOptions)
		Blur(dst, src, opt)

		if !checkTransformTest(t, &oc, dst, src) {
			continue
		}
	}
}

func TestBlurGopher(t *testing.T) {
	src, err := loadImage("../testdata/gopher.png")
	if err != nil {
		t.Error(err)
		return
	}

	dst := image.NewRGBA(src.Bounds())
	Blur(dst, src, &BlurOptions{StdDev: 1.1})

	cmp, err := loadImage("../testdata/gopher-blur.png")
	if err != nil {
		t.Error(err)
		return
	}
	err = imageWithinTolerance(dst, cmp, 0x101)
	if err != nil {
		t.Error(err)
		return
	}
}

func benchBlur(b *testing.B, bounds image.Rectangle) {
	b.StopTimer()

	// Construct a fuzzy image.
	src := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			src.SetRGBA(x, y, image.RGBAColor{
				uint8(5 * x % 0x100),
				uint8(7 * y % 0x100),
				uint8((7*x + 5*y) % 0x100),
				0xff,
			})
		}
	}
	dst := image.NewRGBA(bounds)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Blur(dst, src, &BlurOptions{0.84, 3})
	}
}

func BenchmarkBlur400x400x3(b *testing.B) {
	benchBlur(b, image.Rect(0, 0, 400, 400))
}

// Exactly twice the pixel count of 400x400.
func BenchmarkBlur400x800x3(b *testing.B) {
	benchBlur(b, image.Rect(0, 0, 400, 800))
}

// Exactly twice the pixel count of 400x800
func BenchmarkBlur400x1600x3(b *testing.B) {
	benchBlur(b, image.Rect(0, 0, 400, 1600))
}
