package domain

type edge struct {
	u int
	v int
}

func distAndNextPosToTarget(start *Object, target *Object, level *Level) *Vector {
	if start == target {
		return nil
	}

	q := createQueue()
	q.enqueue(start)

	dx := []int{0, 0, -1, 1}
	dy := []int{-1, 1, 0, 0}
	dirs := []Directions{Forward, Back, Left, Right}

	visit := make([][]int, RoomsInHeight*RegionHeight)
	dist := make([][]int, RoomsInHeight*RegionHeight)
	dirParent := make([][]Directions, RoomsInHeight*RegionHeight)
	parent := make([][]Object, RoomsInHeight*RegionHeight)

	for i := range visit {
		visit[i] = make([]int, RoomsInWidth*RegionWidth)
		dist[i] = make([]int, RoomsInWidth*RegionWidth)
		dirParent[i] = make([]Directions, RoomsInWidth*RegionWidth)
		parent[i] = make([]Object, RoomsInWidth*RegionWidth)
	}

	visit[start.Coordinates[Y]][start.Coordinates[X]] = 1

	for !q.isEmpty() {
		current := q.dequeue()
		currentDist := dist[current.Coordinates[Y]][current.Coordinates[X]]

		for i := range len(dirs) {
			coords := current
			coords.Coordinates[X] += dx[i]
			coords.Coordinates[Y] += dy[i]

			y, x := coords.Coordinates[Y], coords.Coordinates[X]

			if visit[y][x] == 0 && !checkOutsideBorder(&coords, level) && checkUnoccupiedLevel(&coords, level) {
				q.enqueue(&coords)
				dist[y][x] = currentDist + 1
				visit[y][x] = 1
				parent[y][x] = current
				dirParent[y][x] = dirs[i]
			}
		}
	}

	var path *Vector
	if visit[target.Coordinates[Y]][target.Coordinates[X]] == 1 {
		path = createVector()
		coords := *target

		for !CheckEqualCoords(parent[coords.Coordinates[Y]][coords.Coordinates[X]].Coordinates, start.Coordinates) {
			path.pushBack(dirParent[coords.Coordinates[Y]][coords.Coordinates[X]])
			coords = parent[coords.Coordinates[Y]][coords.Coordinates[X]]
		}
		path.pushBack(dirParent[coords.Coordinates[Y]][coords.Coordinates[X]])
		path.reverseVector()
	}

	return path
}
