package sorter

import "os"

type Sorter struct {
	files []os.FileInfo
}

func New(files []os.FileInfo) *Sorter {
	return &Sorter{files}
}

func (s *Sorter) Len() int {
	return len(s.files)
}

func (s *Sorter) Swap(i, j int) {
	s.files[i], s.files[j] = s.files[j], s.files[i]
}

func (s *Sorter) Less(i, j int) bool {
	return s.files[i].Name() < s.files[j].Name()
}
