package ball

import "math"

type ball struct {
	pos    vec2
	vel    vec2
	acc    vec2
	radius float64
	mass   float64
}

func (b *ball) update() {
	b.vel.x += b.acc.x
	b.vel.y += b.acc.y

	b.pos.x += b.vel.x
	b.pos.y += b.vel.y

	mag := (b.vel.x * b.vel.x) + (b.vel.y * b.vel.y)
	if mag > 12*12 {
		b.vel.x *= 0.95
		b.vel.y *= 0.95
	}

	if b.pos.x-b.radius < 0 {
		b.vel.x = b.vel.x * -1.0
		b.pos.x = b.radius
	}
	if b.pos.x+b.radius > screenWidth {
		b.vel.x = b.vel.x * -1.0
		b.pos.x = screenWidth - b.radius
	}
	if b.pos.y+b.radius > screenHeight {
		b.vel.y = b.vel.y * -1.0
		b.pos.y = screenHeight - b.radius
	}
	if b.pos.y-b.radius < 0 {
		b.vel.y = b.vel.y * -1.0
		b.pos.y = b.radius
	}
}

func collide(b1 *ball, b2 *ball) {

	// 1. find unit normal and unit tangent vectors
	n := vec2{x: b2.pos.x - b1.pos.x, y: b2.pos.y - b1.pos.y}
	magnitude := math.Sqrt((n.x * n.x) + n.y*n.y)
	un := vec2{x: n.x / magnitude, y: n.y / magnitude}
	ut := vec2{x: -un.y, y: un.x}

	//fmt.Println(n, magnitude, un, ut)

	// 2. Create the initial(before the collision velocity vectors
	// already done

	// 3. Project velocity vectors onto unit normal and unit tangent vectors
	v1n := dot(un, b1.vel)
	v1t := dot(ut, b1.vel)
	v2n := dot(un, b2.vel)
	v2t := dot(ut, b2.vel)

	// 4. Find new tangential velocities after the collision
	// same as original because no force between circles in the tangential direction
	// v'1t = v1t    v'2t = v2t
	v1tPrime := v1t
	v2tPrime := v2t

	// 5. find new normal velocities
	v1nPrime := (v1n*(b1.mass-b2.mass) + (2 * b2.mass * v2n)) / (b1.mass + b2.mass)
	v2nPrime := (v2n*(b2.mass-b1.mass) + (2 * b1.mass * v1n)) / (b1.mass + b2.mass)

	// 6. convert scalar normal and tangential velocities into vectors
	v1nPrimeVec := vec2{x: v1nPrime * un.x, y: v1nPrime * un.y}
	v1tPrimeVec := vec2{x: v1tPrime * ut.x, y: v1tPrime * ut.y}
	v2nPrimeVec := vec2{x: v2nPrime * un.x, y: v2nPrime * un.y}
	v2tPrimeVec := vec2{x: v2tPrime * ut.x, y: v2tPrime * ut.y}

	// 7. Find final velocity vectors by adding normal and tangential components for each
	v1Prime := vec2{x: v1nPrimeVec.x + v1tPrimeVec.x, y: v1nPrimeVec.y + v1tPrimeVec.y}
	v2Prime := vec2{x: v2nPrimeVec.x + v2tPrimeVec.x, y: v2nPrimeVec.y + v2tPrimeVec.y}

	b1.vel = v1Prime
	b2.vel = v2Prime
}
