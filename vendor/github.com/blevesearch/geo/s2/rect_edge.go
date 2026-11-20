// Copyright (c) 2026 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package s2

import (
	"github.com/blevesearch/geo/s1"
)

// This file adds exact relational predicates between a Rect and a geodesic
// edge
//
// Geometry that runs exactly along the boundary of a rectangle cannot otherwise
// be told apart from geometry that just barely crosses it. Therefore the boundary
// is moved by a few epsilons. This margin is roughly 25 nanometers.
const boundaryMargin = 16 * dblEpsilon

// IntersectsEdge reports whether the rectangle, treated as a closed region,
// intersects the geodesic edge AB.
//
// It requires the points to have unit length.
func (r Rect) IntersectsEdge(a, b Point) bool {
	return r.closure().intersectsEdge(a, b)
}

// intersectsEdge is IntersectsEdge with no margin applied.
func (r Rect) intersectsEdge(a, b Point) bool {
	if r.IsEmpty() {
		return false
	}

	// If either endpoint lies in the rectangle we are done.
	if r.ContainsPoint(a) || r.ContainsPoint(b) {
		return true
	}

	// Otherwise the edge is a connected curve with both of its ends
	// outside the rectangle, so it can only reach the rectangle by
	// meeting its boundary.
	return r.BoundaryIntersectsEdge(a, b)
}

// BoundaryIntersectsEdge reports whether the boundary of the rectangle crosses
// the geodesic edge AB. An edge that lies along the boundary does not cross it.
//
// It requires the points to have unit length.
func (r Rect) BoundaryIntersectsEdge(a, b Point) bool {
	if r.IsEmpty() {
		return false
	}

	if intersectsLngEdge(a, b, r.Lat, s1.Angle(r.Lng.Lo)) ||
		intersectsLngEdge(a, b, r.Lat, s1.Angle(r.Lng.Hi)) ||
		intersectsLatEdge(a, b, s1.Angle(r.Lat.Lo), r.Lng) ||
		intersectsLatEdge(a, b, s1.Angle(r.Lat.Hi), r.Lng) {
		return true
	}

	return false
}

// InteriorIntersectsEdge reports whether the interior of the rectangle
// intersects the geodesic edge AB, that is, whether any part of AB gets inside
// the rectangle rather than running along its boundary.
//
// It requires the points to have unit length.
func (r Rect) InteriorIntersectsEdge(a, b Point) bool {
	return r.interior().intersectsEdge(a, b)
}

// ContainsPolyline reports whether the rectangle contains the polyline.
func (r Rect) ContainsPolyline(pl *Polyline) bool {
	if pl == nil || len(*pl) == 0 {
		return false
	}

	return r.closure().Contains(pl.RectBound())
}

// ContainsPolygon reports whether the rectangle contains the polygon, its
// interior included.
func (r Rect) ContainsPolygon(pgn *Polygon) bool {
	if pgn == nil || pgn.IsEmpty() {
		return false
	}

	return r.closure().Contains(pgn.RectBound())
}

// ContainsRect reports whether the rectangle contains the other rectangle.
func (r Rect) ContainsRect(o Rect) bool {
	return r.closure().Contains(o)
}

// interior returns the rectangle pulled in by boundaryMargin on every side.
func (r Rect) interior() Rect {
	if r.IsEmpty() {
		return r
	}

	inset := r.expandedByMargin(-boundaryMargin)
	if inset.IsEmpty() {
		return r
	}

	return inset
}

// closure returns the rectangle pushed out by boundaryMargin on every side.
func (r Rect) closure() Rect {
	if r.IsEmpty() {
		return r
	}

	return r.expandedByMargin(boundaryMargin)
}

func (r Rect) expandedByMargin(margin float64) Rect {
	return r.expanded(LatLng{Lat: s1.Angle(margin), Lng: s1.Angle(margin)})
}
