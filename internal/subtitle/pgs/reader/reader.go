package reader

import "fmt"

type Reader struct {
	Data           []byte
	Position       int
	LimitInUse     bool
	LimitRemaining int
}

func NewReader(data []byte) *Reader {
	return &Reader{
		Data:           data,
		Position:       0,
		LimitInUse:     false,
		LimitRemaining: 0,
	}
}

func (r *Reader) SetLimit(limit int) {
	r.LimitInUse = true
	r.LimitRemaining = limit
}

func (r *Reader) RemoveLimit() {
	r.LimitInUse = false
	r.LimitRemaining = 0
}

func (r *Reader) RemainingLimit() int {
	return r.LimitRemaining
}

func (r *Reader) SkipPastLimit() {
	r.Position += r.LimitRemaining
	r.RemoveLimit()
}

func (r *Reader) RemainingBytes() int {
	return len(r.Data) - r.Position - 1
}

func (r *Reader) ReachedEnd() bool {
	return r.Position >= len(r.Data)
}

func (r *Reader) Read(count int) ([]byte, error) {
	if count == 0 {
		return make([]byte, 0), nil
	}

	if r.LimitInUse && count > r.LimitRemaining {
		buf, err := r.Read(r.LimitRemaining)
		if err != nil {
			return make([]byte, 0), err
		}

		return append(buf, make([]byte, count-len(buf))...), nil
	}

	r.Position += count
	r.LimitRemaining -= count

	if r.ReachedEnd() {
		return make([]byte, 0), fmt.Errorf("Read position out of bound")
	}

	return r.Data[r.Position-count : r.Position], nil
}

func (r *Reader) ReadByte() (byte, error) {
	if r.LimitInUse && r.LimitRemaining <= 0 {
		return 0x00, nil
	}

	if r.ReachedEnd() {
		return 0x00, fmt.Errorf("Read position out of bound")
	}

	r.Position++
	r.LimitRemaining--

	return r.Data[r.Position-1], nil
}
