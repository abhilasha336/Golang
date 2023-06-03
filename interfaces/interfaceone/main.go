package main

import "fmt"

type Cars interface {
	Mileage() float64
}

type Bmw struct {
	distanceCovered float64
	fuelConsumption float64
	averageDistance float64
}

type Audi struct {
	distanceCovered float64
	fuelConsumption float64
	averageDistance float64
}

func totalMileage(m []Cars) float64 { //slice of interfaces can call mileage() values using this by looping and interface.Mileage()
	t1 := 0.0
	for _, val := range m {
		t1 = t1 + val.Mileage()
	}
	return t1
}

func (b Bmw) Mileage() float64 {
	return b.distanceCovered / b.fuelConsumption
}
func (b Audi) Mileage() float64 {
	return b.distanceCovered / b.fuelConsumption
}

func main() {
	b1 := Bmw{
		distanceCovered: 15.1,
		fuelConsumption: 1,
		averageDistance: 2,
	}
	a1 := Audi{
		distanceCovered: 20.1,
		fuelConsumption: 2,
		averageDistance: 1.5,
	}

	var a Cars
	a = a1
	fmt.Println("mileage:from a is", a.Mileage())
	fmt.Println("mileage:", b1.Mileage())
	fmt.Println("mileage:", a1.Mileage())

	person := []Cars{a1, b1}

	totalMil := totalMileage(person)
	fmt.Println("totalMIleage is", totalMil)

}

func (b Bmw) AverageSpeed() float64 {
	return b.averageDistance
}

func averageSpeed(c Cars) float64 { //slice of interfaces can call mileage() values using this by looping and interface.Mileage()
	as := c.(Bmw) //type assertion  AverageSpeed is not inside interface{Mileage()}
	return as.AverageSpeed()
}
