package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const SCREEN_WIDTH = 800
const SCREEN_HEIGHT = 600
const FRAME_RATE = 60 // frames per second

const STAR_RADIUS = 4.0
const UNIVERSE_DEPTH = 2000.0

const VANISHING_X = SCREEN_WIDTH / 2
const VANISHING_Y = SCREEN_HEIGHT / 2
const VANISHING_Z = UNIVERSE_DEPTH

type ThreeDPoint struct {
	x float64
	y float64
	z float64
}

type Vector struct {
	angle     float64
	magnitude float64
}

// define a struct for a three dimensional direction and speed
type ThreeDSpeed struct {
	x_speed float64
	y_speed float64
	z_speed float64
}

type ThreeDLine struct {
	start ThreeDPoint
	end   ThreeDPoint
	color color.RGBA
	speed ThreeDSpeed
}

type Scene struct {
	stars []ThreeDPoint
	nodes []ThreeDPoint
	lines []ThreeDLine
}

func rand_speed() float64 {
	max_speed := 50.0
	return (rand.Float64() * max_speed) - (max_speed / 2)
}

func (s *Scene) GeneratePoints() {
	// Generate a bunch of points
	s.nodes = make([]ThreeDPoint, 0)
	s.lines = make([]ThreeDLine, 0)

	distant_point := ThreeDPoint{SCREEN_WIDTH / 2, SCREEN_HEIGHT / 2, UNIVERSE_DEPTH}

	s.nodes = append(s.nodes, distant_point)

	for i := 0; i < 20; i++ {
		color := color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255}
		starting_point := ThreeDPoint{rand.Float64() * SCREEN_WIDTH, rand.Float64() * SCREEN_HEIGHT, 0}
		speed_per_second := ThreeDSpeed{rand_speed(), rand_speed(), rand_speed()}
		s.lines = append(s.lines, ThreeDLine{starting_point, distant_point, color, speed_per_second})
	}

	// generate a thousand stars
	s.stars = make([]ThreeDPoint, 0)
	for i := 0; i < 1000; i++ {
		depth := rand.Float64() * UNIVERSE_DEPTH
		star := ThreeDPoint{rand.Float64() * SCREEN_WIDTH, rand.Float64() * SCREEN_HEIGHT, depth}
		s.stars = append(s.stars, star)
	}
}

func (s *Scene) Update() error {
	// move each line's start point by the speed
	for i := 0; i < len(s.lines); i++ {
		line := &s.lines[i]
		line.start.x += line.speed.x_speed / FRAME_RATE
		line.start.y += line.speed.y_speed / FRAME_RATE
		line.start.z += line.speed.z_speed / FRAME_RATE
	}

	// all stars move towards the camera at the same speed
	for i := 0; i < len(s.stars); i++ {
		star := &s.stars[i]
		star.z -= 1
		if star.z < 0-UNIVERSE_DEPTH {
			star.z = UNIVERSE_DEPTH
		}
	}

	return nil
}

func (s *Scene) Draw(screen *ebiten.Image) {
	for _, line := range s.lines {
		ebitenutil.DrawLine(screen, line.start.x, line.start.y, line.end.x, line.end.y, line.color)
	}

	for _, star := range s.stars {
		star_color := color.RGBA{255, 255, 255, 128}

		star_x := star.x + ((VANISHING_X-star.x)/VANISHING_Z)*star.z
		star_y := star.y + ((VANISHING_Y-star.y)/VANISHING_Z)*star.z

		scale_factor := (UNIVERSE_DEPTH / (UNIVERSE_DEPTH + star.z))

		// adjust star width and height based on distance
		star_radius := STAR_RADIUS * scale_factor

		// fmt.Printf("Star %d - depth: %f, orig coords: %f, %f coords: %f, %f, radius: %f\n", i, star.z, star.x, star.y, star_x, star_y, star_radius)
		// star_color = color.RGBA{255, 0, 0, 255}

		ebitenutil.DrawCircle(screen, star_x, star_y, star_radius, star_color)
	}
}

func (s *Scene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func (s *Scene) Describe() {
	fmt.Println("Points")
	for _, point := range s.nodes {
		fmt.Println(point)
	}

	fmt.Println("Lines")
	for _, line := range s.lines {
		fmt.Println(line)
	}
}

func main() {
	scene := &Scene{}
	scene.GeneratePoints()

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("3D Lines")

	scene.Describe()

	if err := ebiten.RunGame(scene); err != nil {
		log.Fatal(err)
	}

}
