package mathutils

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func NCk(n int, k int) float64 {
	// http://blog.plover.com/math/choose.html
	if k > n {
		return 0
	}
	if k == 0 {
		return 1
	}
	r := float64(1)
	for d := 1; d <= k; d++ {
		r *= float64(n)
		r /= float64(d)
		n--
	}
	return r
}

func Factorial(n int) float64 {
	// unoptimized, called only on small n
	if n < 2 {
		return 1
	}
	f := float64(1)
	for i := 2; i <= n; i++ {
		f *= float64(i)
	}
	return f
}
