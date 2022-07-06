package numbers

// Add returns a + b
func Add[N Number](a N, b N) N {
    return a + b
}

// Sub returns a - b
func Sub[N Number](a N, b N) N {
    return a - b
}
