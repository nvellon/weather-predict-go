package solarsystem

import (
	"math"
)

type (
	point struct {
		X float64
		Y float64
	}

	planet struct {
		AngularSpeed    float64
		Angle           float64
		DistanceFromSun float64
		Location        *point
	}

	SolarSystem struct {
		days      int
		Sun       point
		Ferengi   *planet
		Betasoide *planet
		Vulcano   *planet
	}
)

const (
	InitialAngle      = 90
	FerengiSpeed      = 1
	BetasoideSpeed    = 3
	VulcanoSpeed      = -5
	FerengiDistance   = 500
	BetasoideDistance = 2000
	VulcanoDistance   = 1000
)

func NewSolarSystem() *SolarSystem {
	ss := SolarSystem{}

	ss.days = 0

	ss.Sun = point{0, 0}

	ss.Ferengi = newPlanet(deg2rad(FerengiSpeed), FerengiDistance)
	ss.Betasoide = newPlanet(deg2rad(BetasoideSpeed), BetasoideDistance)
	ss.Vulcano = newPlanet(deg2rad(VulcanoSpeed), VulcanoDistance)

	return &ss
}

func newPlanet(angularSpeed float64, distanceFromSun float64) *planet {
	p := planet{
		AngularSpeed:    angularSpeed,
		Angle:           deg2rad(InitialAngle),
		DistanceFromSun: distanceFromSun,
		Location:        &point{0, distanceFromSun},
	}
	return &p
}

func (ss *SolarSystem) NextDay() {
	ss.days += 1

	ss.Ferengi.nextDay()
	ss.Betasoide.nextDay()
	ss.Vulcano.nextDay()
}

func (ss *SolarSystem) GetDay() int {
	return ss.days
}

func (p *planet) nextDay() {
	p.Angle = p.Angle + p.AngularSpeed

	p.Location.X = p.DistanceFromSun * math.Cos(p.Angle)
	p.Location.Y = p.DistanceFromSun * math.Sin(p.Angle)
}

func (ss *SolarSystem) IsDrought() bool {
	a := arePointsAligned(ss.Ferengi.Location, ss.Betasoide.Location, ss.Vulcano.Location)
	b := arePointsAligned(ss.Ferengi.Location, ss.Betasoide.Location, &ss.Sun)
	return a && b
}

// OptimumTemperaturePressure returns true when
// the planets are aligned beetween them
func (ss *SolarSystem) IsOptimumTemperaturePressure() bool {
	a := arePointsAligned(ss.Ferengi.Location, ss.Betasoide.Location, ss.Vulcano.Location)
	b := arePointsAligned(ss.Ferengi.Location, ss.Betasoide.Location, &ss.Sun)
	return a && !b
}

// pointsAligned return true when three points are aligned
// Following the approach defined by this post:
// https://www.urbanpro.com/gre/how-to-determine-if-points-are-collinear
func arePointsAligned(a *point, b *point, c *point) bool {
	return ((a.X-b.X)*(b.Y-c.Y))-((b.X-c.X)*(a.Y-b.Y))/2 == 0
}

// RainSeason returns true when
// the planets are forming a triangle
// and the sun is inside this triangle's area
func (ss *SolarSystem) IsRainSeason() bool {
	// if planets don't form a triangle
	// no need to do more checks
	if arePointsAligned(ss.Ferengi.Location, ss.Betasoide.Location, ss.Vulcano.Location) {
		return false
	}

	return pointInTriangle(&ss.Sun, ss.Ferengi.Location, ss.Betasoide.Location, ss.Vulcano.Location)
}

func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180
}

// pointInTriangle returns true when
// the point pt is inside de area of the triangle
// formed by v1, v2 and v3
// Following the approach from this post:
// https://stackoverflow.com/questions/2049582/how-to-determine-if-a-point-is-in-a-2d-triangle
func pointInTriangle(pt *point, v1 *point, v2 *point, v3 *point) bool {
	var d1, d2, d3 float64
	var has_neg, has_pos bool

	d1 = sign(pt, v1, v2)
	d2 = sign(pt, v2, v3)
	d3 = sign(pt, v3, v1)

	has_neg = (d1 < 0) || (d2 < 0) || (d3 < 0)
	has_pos = (d1 > 0) || (d2 > 0) || (d3 > 0)

	return !(has_neg && has_pos)
}

func sign(p1 *point, p2 *point, p3 *point) float64 {
	return (p1.X-p3.X)*(p2.Y-p3.Y) - (p2.X-p3.X)*(p1.Y-p3.Y)
}
