package df2014

type WorldDat struct {
	Version     uint32
	Compression uint32
}

func (r *Reader) WorldDat() (w WorldDat, err error) {
	w.Version, w.Compression, err = r.Header()
	if err != nil {
		return
	}

	return
}
