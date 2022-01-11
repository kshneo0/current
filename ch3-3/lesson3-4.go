package main

import (
	"container/heap"
	"fmt"
	"image"
	"time"
)

const (
	WIDTH = 128
	HEIGHT = 128
	THRESHOLD = 0x7777
)

func outOfBounds(x, y int) bool {
	return x >= WIDTH || y >= HEIGHT || y < 0 || x < 0
}

func retracePath(c *Vertex) [][2]int {
	path := make([][2]int, 0)
	for c.parent != nil {
		path = append(path, c.value)
		c = c.parent
	}
	return path
}

func shortestPath(img image.Image, source [2]int, dest [2]int) [][2]int {
	//create a @D slice of nil Vertex pointers that are the size of the image
	G := make([][]*Vertex, HEIGHT)
	for i := 0; i< HEIGHT; i++ {
		G[i] = make([]*Vertex, WIDTH)
	}

	sourceVertex := &Vertex{source,0,0,nil}

	//initialize priority queue
	pq := make(PriorityQueue,1)
	pq[0] = sourceVertex
	heap.Init(&pq)

	// y x
	G[source[1]][source[0]] = sourceVertex

	// 	n
	// n @ n
    //   n
	neighbors := [4][2]int{{0,1}, {1,0},{-1,0},{0,-1}}

	for pq.Len() > 0 {
		//current Vertex
		c := heap.Pop(&pq).(*Vertex)

		// if we arrived at destination
		if c.value == dest {
			return retracePath(c)
		}

		for _, n := range neighbors {
			func(n [2]int) {
				vx, vy := c.value[0]+n[0], c.value[1]+n[1]
				if outOfBounds(vx, vy) {
					return
				}

				r, g, b, _ := (img).At(vx, vy).RGBA()
				if r + b + g < THRESHOLD {
					return
				}

				edgeWwight := 1

				p := c.priority + edgeWwight

				v := G[vy][vx]

				if v == nil {
					v = &Vertex{[2]int{vx,vy},p,0,c}
					G[vy][vx] = v
				} else if v.priority > p {
					v = &Vertex{[2]int{vx,vy}, p, 0, c}
				} else {
					return
				}

				heap.Push(&pq, v)
			}(n)
		}
	}
	return nil
}

func main() {
	m := loadAndResizeImg("maze.jpg",WIDTH, HEIGHT)
	
	canvas := createCanvasWithImage(m)
	drawPath(canvas,[][2]int{{10,8}})
	drawPath(canvas,[][2]int{{64,64}})

	start := time.Now()
	path := shortestPath(m, [2]int{10,8}, [2]int{64,64})
	duration := time.Since(start)
	fmt.Println(duration)

	drawPath(canvas, path)
//	fmt.Println(path)
	exportCanvasToFile(canvas,"maze-sln.jpeg")
}