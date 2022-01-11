package main

import "fmt"

type void struct {}

var member void

type Set struct {
	members map[int]void
}

func newSet() Set{
	a := make(map[int]void)
	return Set{a}
}

func(a Set) Add(n int) {
	a.members[n] = member
}

func (a Set) Union(b Set) Set {
	c := newSet()
	for k := range a.members {
		c.members[k] = member
	}
	for k := range a.members {
		c.members[k] = member
	}
	return c
}

func (a Set) Contains(k int) bool {
	_, exists := a.members[k]
	return exists
}

func (a Set) Intersection(b Set) Set {
	c := newSet()
	for k := range a.members {
		if b.Contains(k) {
			c.members[k] = member
		}		
	}
	return c
}

func (a Set) Subtract(b Set) Set {
	c := newSet()
	for k := range a.members {
		if !b.Contains(k) {
			c.members[k] = member
		}
	}
}

func (a Set) Remove(n int) {
	delete(a.members,n)
}

func (a Set) Empty() bool {
	return len(a.members) == 0
}

func main() {
	
	setA := newSet()
	setB := newSet()

	setA.Add(1)
	setA.Add(2)

	setB.Add(2)
	setB.Add(3)
	setB.Add(3)

	fmt.Println("setA:", setA.members)
	fmt.Println("setB:", setB.members)

	fmt.Println("Union of setA and setB:", setA.Union(setB).members)

	setC := newSet()
	fmt.Println("setC is empty?: ", setC.Empty())
	fmt.Println("setA is empty?: ", setA.Empty())

	fmt.Println("Intersection of setA and setB", setA.Intersection(setB))

	setD := newSet()
	setD.Add(4)
	fmt.Println("setD", setD.members)
	setD.Remove(4)
	fmt.Println("setD", setD.members)

	fmt.Println("setA - setB:", setA.Subtract(setB))
	fmt.Println("setB - setA:", setB.Subtract(setA))
}