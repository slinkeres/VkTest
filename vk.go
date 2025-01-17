package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct{
	x, y int
}

type Node struct{
	point Point
	priority int
}

type PriorityQueue []Node

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool{
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int){
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}){
	*pq = append(*pq, x.(Node))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

var directions = []Point{
	{0, 1}, 
	{1, 0},
	{0,-1}, 
	{-1, 0}, 
}

func parseInput() ([][]int, Point, Point, error){
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan(){
		return nil, Point{}, Point{}, errors.New("Не указаны размеры лабиринта")
	}
	size := strings.Fields(scanner.Text())
	if len(size) != 2{
		return nil, Point{}, Point{}, errors.New("Неправильный формат размера лабиринта")
	}
	rows, err1 := strconv.Atoi(size[0])
	cols, err2 := strconv.Atoi(size[1])
	if err1 != nil || err2 != nil || rows <= 0 || cols <= 0{
		return nil, Point{}, Point{}, errors.New("Неверный размер лабиринта")
	}
	maze := make([][]int, rows)
	for i := 0;i < rows; i++{
		if !scanner.Scan(){
			return nil, Point{}, Point{}, errors.New("Недостаточно строк в лабиринте")
		}
		line := strings.Fields(scanner.Text())
		if len(line) != cols{
			return nil, Point{}, Point{}, errors.New("Неверное количество столбцов в строке лабиринта")
		}
		maze[i] = make([]int, cols)
		for j, cell := range line {
			maze[i][j], err1 = strconv.Atoi(cell)
			if err1 != nil || maze[i][j] < 0 || maze[i][j] > 9 {
				return nil, Point{}, Point{}, errors.New("Неверное значение клетки лабиринта")
			}
		}
	}
	if !scanner.Scan(){
		return nil, Point{}, Point{}, errors.New("Не указаны стартовая и конечная точки")
	}
	points := strings.Fields(scanner.Text())
	if len(points) != 4{
		return nil, Point{}, Point{}, errors.New("Неправильный формат стартовой и конечной точек")
	}
	startX, err1 := strconv.Atoi(points[0])
	startY, err2 := strconv.Atoi(points[1])
	endX, err3 := strconv.Atoi(points[2])
	endY, err4 := strconv.Atoi(points[3])
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil{
		return nil, Point{}, Point{}, errors.New("Неверные координаты точек")
	}
	start := Point{startX, startY}
	end := Point{endX, endY}
	if !isValidPoint(start, rows, cols) || !isValidPoint(end, rows, cols){
		return nil, Point{}, Point{}, errors.New("Стартовая или конечная точка вне границ лабиринта")
	}
	if maze[startX][startY] == 0{
		return nil, Point{}, Point{}, errors.New("Стартовая точка не может быть стеной")
	}
	if maze[endX][endY] == 0{
		return nil, Point{}, Point{}, errors.New("Конечная точка не может быть стеной")
	}
	return maze, start, end, nil
}

func isValidPoint(p Point, rows, cols int) bool{
	return p.x >= 0 && p.x < rows && p.y >= 0 && p.y < cols
}

func findShortestPath(maze [][]int, start, end Point) ([]Point, error){
	rows := len(maze)
	cols := len(maze[0])
	const inf = 1 << 30
	cost := make([][]int, rows)
	for i := range cost{
		cost[i] = make([]int, cols)
		for j := range cost[i]{
			cost[i][j] = inf
		}
	}
	cost[start.x][start.y] = 0
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, Node{start, 0})
	prev := make(map[Point]*Point)
	for pq.Len() > 0 {
		curr := heap.Pop(pq).(Node)
		if curr.point == end{
			break
		}
		for _, d := range directions{
			neighbor := Point{curr.point.x + d.x, curr.point.y + d.y}
			if isValidPoint(neighbor, rows, cols) && maze[neighbor.x][neighbor.y] != 0{
				newCost := cost[curr.point.x][curr.point.y] + maze[neighbor.x][neighbor.y]
				if newCost < cost[neighbor.x][neighbor.y]{
					cost[neighbor.x][neighbor.y] = newCost
					heap.Push(pq, Node{neighbor, newCost})
					prev[neighbor] = &curr.point
				}
			}
		}
	}
	if cost[end.x][end.y] == inf{
		return nil, errors.New("Путь не найден")
	}
	path := []Point{}
	for p := &end; p != nil; p = prev[*p]{
		path = append([]Point{*p}, path...)
	}
	return path, nil
}

func main(){
	maze, start, end , err := parseInput()
	if err != nil{
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	path, err := findShortestPath(maze, start, end)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, p := range path{
		fmt.Printf("%d %d\n", p.x, p.y)
	}
	fmt.Println(".")
}