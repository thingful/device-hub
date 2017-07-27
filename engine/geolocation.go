// Copyright © 2017 thingful
//
// Geolocation minimal js object mimic
// https://dev.w3.org/geo/api/spec-source.html#position

package engine

type coordinates struct {
	Latitude         float64 // specified in decimal degrees
	Longitude        float64 // specified in decimal degrees
	Accuracy         float64 // specified in meters
	Altitude         float64 // specified in meters
	AltitudeAccuracy float64 // specified in meters
	Heading          float64 // specified in degrees
	Speed            float64 // specified in m/s
}

type position struct {
	Coords    coordinates
	Timestamp int64 // represents the time when the position was acquired
}

func (p position) GetLat() float64 {
	return p.Coords.Latitude
}
