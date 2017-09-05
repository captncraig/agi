//Package decrypt implements decryption of agi files from a key embedded in a Sierra loader
package decrypt

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Decryptor interface {
	ReadFile(name string) ([]byte, error)
}

type decrypt struct {
	key [128]byte
}

func New(loaderFile string) (Decryptor, error) {
	dat, err := ioutil.ReadFile(loaderFile)
	if err != nil {
		return nil, err
	}
	const delim = "keyOfs"
	idx := strings.Index(string(dat), delim)
	if idx == -1 {
		return nil, fmt.Errorf("Key not found in loader file")
	}
	start := idx + len(delim) + 2
	key := dat[start : start+128]
	d := &decrypt{}
	copy(d.key[:], key)
	return d, nil
}

func (d *decrypt) ReadFile(name string) ([]byte, error) {
	dat, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	//key := d.key
	for i := range dat {
		keyOff := i % 128
		if keyOff == 0 && i != 0 {
			fmt.Println("ROTATE")
			break
		}
		//dat[i] ^= key[keyOff]
	}
	fmt.Println(string(dat))
	fmt.Println(dat)
	return dat, nil
}
