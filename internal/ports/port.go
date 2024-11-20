package ports

import "strconv"

type Port int

func (p Port) String() string {
	return strconv.Itoa(int(p))
}
