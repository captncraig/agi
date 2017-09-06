package binutils

type Reader struct {
	dat    []byte
	idx    int
	labels []int
	Taken  []byte
}

func NewReader(dat []byte) *Reader {
	return &Reader{
		dat:    dat,
		labels: []int{0},
	}
}

func (r *Reader) Take() byte {
	b := r.dat[r.idx]
	r.idx++
	r.Taken = append(r.Taken, b)
	return b
}

func (r *Reader) Peek() byte {
	return r.dat[r.idx]
}

func (r *Reader) Addr() int {
	return r.idx
}

func (r *Reader) ResetTaken() {
	r.Taken = nil
}

func (r *Reader) Take16() uint16 {
	return uint16(r.Take()) | uint16(r.Take())<<8
}

func (r *Reader) Take16AsInt() int {
	return int(r.Take16())
}

func (r *Reader) HasMore() bool {
	return r.idx < len(r.dat)
}

func (r *Reader) MarkCurrent() int {
	r.labels = append(r.labels, r.idx)
	return len(r.labels) - 1
}

func (r *Reader) TakeFrom(lbl int, offset int) byte {
	return r.dat[r.labels[lbl]+offset]
}

func (r *Reader) Take16From(lbl int, offset int) uint16 {
	return uint16(r.TakeFrom(lbl, offset)) | uint16(r.TakeFrom(lbl, offset+1))<<8
}

func (r *Reader) Take16AsIntFrom(lbl int, offset int) int {
	return int(r.Take16From(lbl, offset))
}
