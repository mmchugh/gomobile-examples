package main

import (
	"log"
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

var (
	program      gl.Program
	vertCoord    gl.Attrib
	vertTexCoord gl.Attrib
	projection   gl.Uniform
	view         gl.Uniform
	model        gl.Uniform
	buf          gl.Buffer
	texture      gl.Texture
	touchLoc     geom.Point
	started      time.Time
)

func main() {
	app.Run(app.Callbacks{
		Start:  start,
		Stop:   stop,
		Draw:   draw,
		Touch:  touch,
		Config: config,
	})
}

func start() {
	var err error
	program, err = glutil.CreateProgram(vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	buf = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, buf)
	gl.BufferData(gl.ARRAY_BUFFER, EncodeObject(cubeData), gl.STATIC_DRAW)

	vertCoord = gl.GetAttribLocation(program, "vertCoord")
	vertTexCoord = gl.GetAttribLocation(program, "vertTexCoord")

	projection = gl.GetUniformLocation(program, "projection")
	view = gl.GetUniformLocation(program, "view")
	model = gl.GetUniformLocation(program, "model")

	texture = loadTexture("gopher.png")

	started = time.Now()
}

func stop() {
	gl.DeleteProgram(program)
	gl.DeleteBuffer(buf)
}

func config(new, old event.Config) {
	touchLoc = geom.Point{new.Width / 2, new.Height / 2}
}

func touch(t event.Touch, c event.Config) {
	touchLoc = t.Loc
}

func draw(c event.Config) {
	since := time.Now().Sub(started)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Clear(gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(program)

	m := mgl.Perspective(0.785, float32(c.Width/c.Height), 0.1, 10.0)
	gl.UniformMatrix4fv(projection, m[:])

	eye := mgl.Vec3{3, 3, 3}
	center := mgl.Vec3{0, 0, 0}
	up := mgl.Vec3{0, 1, 0}

	m = mgl.LookAtV(eye, center, up)
	gl.UniformMatrix4fv(view, m[:])

	m = mgl.HomogRotate3D(float32(since.Seconds()), mgl.Vec3{0, 1, 0})
	gl.UniformMatrix4fv(model, m[:])

	gl.BindBuffer(gl.ARRAY_BUFFER, buf)

	gl.EnableVertexAttribArray(vertCoord)
	gl.VertexAttribPointer(vertCoord, coordsPerVertex, gl.FLOAT, false, 20, 0) // 4 bytes in float, 5 values per vertex

	gl.EnableVertexAttribArray(vertTexCoord)
	gl.VertexAttribPointer(vertTexCoord, texCoordsPerVertex, gl.FLOAT, false, 20, 12)

	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)

	gl.DisableVertexAttribArray(vertCoord)

	debug.DrawFPS(c)
}

var cubeData = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

	// Right
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
}

var (
	coordsPerVertex    = 3
	texCoordsPerVertex = 2
	vertexCount        = len(cubeData) / (coordsPerVertex + texCoordsPerVertex)
)

const vertexShader = `#version 100
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

attribute vec3 vertCoord;
attribute vec2 vertTexCoord;

varying vec2 fragTexCoord;

void main() {
	fragTexCoord = vertTexCoord;
    gl_Position = projection * view * model * vec4(vertCoord, 1);
}`

const fragmentShader = `#version 100
precision mediump float;

uniform sampler2D tex;

varying vec2 fragTexCoord;

void main() {
    gl_FragColor = texture2D(tex, fragTexCoord);
}`
