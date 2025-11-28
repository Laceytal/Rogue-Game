package domain

type Vector struct {
	Data []Directions
}

func createVector() *Vector {
	return &Vector{
		Data: make([]Directions, 0),
	}
}

func (v *Vector) pushBack(dir Directions) {
	v.Data = append(v.Data, dir)
}

func (v *Vector) reverseVector() {
	for i := range len(v.Data) / 2 {
		j := len(v.Data) - i - 1
		v.Data[i], v.Data[j] = v.Data[j], v.Data[i]
	}
}
