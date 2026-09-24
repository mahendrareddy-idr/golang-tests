package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	for {
		showMenu()

		fmt.Print("Choose an option: ")
		var choice string
		fmt.Scanln(&choice)

		switch strings.ToLower(choice) {
		case "1", "a":
			performAdd()
		case "2", "s":
			performSubtract()
		case "3", "m":
			performMultiply()
		case "4", "d":
			performDivide()
		case "5", "p":
			performPower()
		case "6", "r":
			performSquareRoot()
		case "7", "i":
			performSin()
		case "8", "c":
			performCos()
		case "9", "t":
			performTan()
		case "10", "l":
			performLog()
		case "11", "n":
			performNaturalLog()
		case "0", "e":
			fmt.Println("Exiting scientific calculator...")
			os.Exit(0)
		default:
			fmt.Println("Invalid option. Please try again.")
		}

		fmt.Println()
	}
}

func showMenu() {
	fmt.Println("====================================")
	fmt.Println("Scientific Calculator")
	fmt.Println("====================================")
	fmt.Println("1. Add (A)")
	fmt.Println("2. Subtract (S)")
	fmt.Println("3. Multiply (M)")
	fmt.Println("4. Divide (D)")
	fmt.Println("5. Power (P)")
	fmt.Println("6. Square Root (R)")
	fmt.Println("7. Sine (I)")
	fmt.Println("8. Cosine (C)")
	fmt.Println("9. Tangent (T)")
	fmt.Println("10. Log Base 10 (L)")
	fmt.Println("11. Natural Log (N)")
	fmt.Println("0. Exit (E)")
}

func readNumber(prompt string) float64 {
	for {
		fmt.Print(prompt)
		var input string
		fmt.Scanln(&input)

		value, err := strconv.ParseFloat(input, 64)
		if err == nil {
			return value
		}

		fmt.Println("Invalid number. Please enter a valid numeric value.")
	}
}

func performAdd() {
	a := readNumber("Enter first number: ")
	b := readNumber("Enter second number: ")
	fmt.Printf("Result: %.2f\n", add(a, b))
}

func add(a, b float64) float64 {
	return a + b
}

func performSubtract() {
	a := readNumber("Enter first number: ")
	b := readNumber("Enter second number: ")
	fmt.Printf("Result: %.2f\n", subtract(a, b))
}

func subtract(a, b float64) float64 {
	return a - b
}

func performMultiply() {
	a := readNumber("Enter first number: ")
	b := readNumber("Enter second number: ")
	fmt.Printf("Result: %.2f\n", multiply(a, b))
}

func multiply(a, b float64) float64 {
	return a * b
}

func performDivide() {
	a := readNumber("Enter dividend: ")
	b := readNumber("Enter divisor: ")
	if b == 0 {
		fmt.Println("Error: division by zero.")
		return
	}
	fmt.Printf("Result: %.2f\n", divide(a, b))
}

func divide(a, b float64) float64 {
	return a / b
}

func performPower() {
	base := readNumber("Enter base: ")
	exponent := readNumber("Enter exponent: ")
	fmt.Printf("Result: %.2f\n", power(base, exponent))
}

func power(base, exponent float64) float64 {
	return math.Pow(base, exponent)
}

func performSquareRoot() {
	value := readNumber("Enter value: ")
	if value < 0 {
		fmt.Println("Error: square root of a negative number is not allowed.")
		return
	}
	fmt.Printf("Result: %.2f\n", squareRoot(value))
}

func squareRoot(value float64) float64 {
	return math.Sqrt(value)
}

func performSin() {
	angle := readNumber("Enter angle in degrees: ")
	fmt.Printf("Result: %.2f\n", sinDegrees(angle))
}

func sinDegrees(angleInDegrees float64) float64 {
	return math.Sin(angleInDegrees * math.Pi / 180)
}

func performCos() {
	angle := readNumber("Enter angle in degrees: ")
	fmt.Printf("Result: %.2f\n", cosDegrees(angle))
}

func cosDegrees(angleInDegrees float64) float64 {
	return math.Cos(angleInDegrees * math.Pi / 180)
}

func performTan() {
	angle := readNumber("Enter angle in degrees: ")
	fmt.Printf("Result: %.2f\n", tanDegrees(angle))
}

func tanDegrees(angleInDegrees float64) float64 {
	return math.Tan(angleInDegrees * math.Pi / 180)
}

func performLog() {
	value := readNumber("Enter value: ")
	if value <= 0 {
		fmt.Println("Error: log is only defined for positive numbers.")
		return
	}
	fmt.Printf("Result: %.2f\n", logBase10(value))
}

func logBase10(value float64) float64 {
	return math.Log10(value)
}

func performNaturalLog() {
	value := readNumber("Enter value: ")
	if value <= 0 {
		fmt.Println("Error: natural log is only defined for positive numbers.")
		return
	}
	fmt.Printf("Result: %.2f\n", naturalLog(value))
}

func naturalLog(value float64) float64 {
	return math.Log(value)
}
