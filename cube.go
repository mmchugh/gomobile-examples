package main

import (
	"encoding/binary"
	"log"
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/app/debug"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
	"golang.org/x/mobile/gl/glutil"
)

func Mat2Float(m *f32.Mat4) []float32 {
	return []float32{
		m[0][0], m[0][1], m[0][2], m[0][3],
		m[1][0], m[1][1], m[1][2], m[1][3],
		m[2][0], m[2][1], m[2][2], m[2][3],
		m[3][0], m[3][1], m[3][2], m[3][3],
	}
}

var (
	program   gl.Program
	vertCoord gl.Attrib
	//	vertTexCoord gl.Attrib
	projection gl.Uniform
	view       gl.Uniform
	model      gl.Uniform
	buf        gl.Buffer

	touchLoc geom.Point
	started  time.Time
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
	gl.BufferData(gl.ARRAY_BUFFER, cubeData, gl.STATIC_DRAW)

	vertCoord = gl.GetAttribLocation(program, "vertCoord")

	projection = gl.GetUniformLocation(program, "projection")
	view = gl.GetUniformLocation(program, "view")
	model = gl.GetUniformLocation(program, "model")

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

	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(program)

	identity := &f32.Mat4{}
	identity.Identity()

	m := &f32.Mat4{}
	m.Perspective(f32.Radian(0.785), float32(c.Width/c.Height), 0.1, 10.0)
	gl.UniformMatrix4fv(projection, Mat2Float(m))

	eye := f32.Vec3{3, 3, 3}
	center := f32.Vec3{0, 0, 0}
	up := f32.Vec3{0, 1, 0}

	m.LookAt(&eye, &center, &up)
	gl.UniformMatrix4fv(view, Mat2Float(m))

	m.Identity()
	m.Rotate(m, f32.Radian(since.Seconds()), &f32.Vec3{0, 1, 0})
	gl.UniformMatrix4fv(model, Mat2Float(m))

	gl.BindBuffer(gl.ARRAY_BUFFER, buf)

	gl.EnableVertexAttribArray(vertCoord)
	gl.VertexAttribPointer(vertCoord, coordsPerVertex, gl.FLOAT, false, 0, 0)
	//	gl.EnableVertexAttribArray(texture)
	//	gl.VertexAttribPointer(vertCoord, texCoordsPerVertex, gl.FLOAT, false, 5, 3)

	gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)

	gl.DisableVertexAttribArray(vertCoord)

	debug.DrawFPS(c)
}

var cubeData = f32.Bytes(binary.LittleEndian,
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, 1.0,
	-1.0, -1.0, 1.0,

	// Top
	-1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0,
	1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,

	// Left
	-1.0, -1.0, 1.0,
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,

	// Right
	1.0, -1.0, 1.0,
	1.0, -1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, 1.0,
)

var (
	coordsPerVertex    = 3
	texCoordsPerVertex = 2
	vertexCount        = len(cubeData) / coordsPerVertex
)

const vertexShader = `#version 100
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

attribute vec3 vertCoord;

void main() {
    gl_Position = projection * view * model * vec4(vertCoord, 1);
}`

const fragmentShader = `#version 100
void main() {
    gl_FragColor = vec4(1.0, 0.0, 0.0, 1.0);
}`
