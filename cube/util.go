package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	image_draw "image/draw"
	_ "image/png"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/gl"
)

// EncodeObject converts float32 vertices into a LittleEndian byte array.
func EncodeObject(vertices ...[]float32) []byte {
	buf := bytes.Buffer{}
	for _, v := range vertices {
		err := binary.Write(&buf, binary.LittleEndian, v)
		if err != nil {
			panic(fmt.Sprintln("binary.Write failed:", err))
		}
	}

	return buf.Bytes()
}

func loadTexture(name string) gl.Texture {
	imgFile, _ := asset.Open(name)
	img, _, _ := image.Decode(imgFile)

	rgba := image.NewRGBA(img.Bounds())
	image_draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, image_draw.Src)

	newTexture := gl.CreateTexture()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, newTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		rgba.Rect.Size().X,
		rgba.Rect.Size().Y,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		rgba.Pix)

	return newTexture
}
