package domain

func makeSets(parent []int, rank []int, size int) {
	for i := range size {
		parent[i] = i
		rank[i] = 0
	}
}

func findSet(v int, parent []int) int {
	if v == parent[v] {
		return v
	}
	parent[v] = findSet(parent[v], parent)
	return parent[v]
}

func unionSets(v int, u int, parent []int, rank []int) {
	v = findSet(v, parent)
	u = findSet(u, parent)

	if u != v {
		if rank[u] >= rank[v] {
			parent[v] = u
		} else {
			parent[u] = v
		}

		if rank[u] == rank[v] {
			rank[u]++
		}
	}
}
