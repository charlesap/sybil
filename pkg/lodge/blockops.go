package lodge

import (
	"fmt"
//	"bytes"
	"errors"
//	"os"
//	"time"
//	"crypto/rand"
//	"crypto/ed25519"
//	"crypto/sha256"
//        "path/filepath"
	"encoding/binary"

//	"github.com/nofeaturesonlybugs/z85"
)

func Hash2block( h *Hash, i int, l uint64) (uint64,error) { //i may range from 0 to 19 on bounce to next hash location

	if l < 1 {return 0,errors.New("limit must be greater than zero")}
	return (uint64(h[i])+
	      (uint64(h[i+1])<<8)+
	      (uint64(h[i+2])<<16)+
	      (uint64(h[i+3])<<24)+
	      (uint64(h[i+4])<<32)+
	      (uint64(h[i+5])<<40)+
	      (uint64(h[i+6])<<48)+
	      (uint64(h[i+7])<<56))%l,nil
}


func (b Base) isfree( i uint64, s int) (avail bool) {

   avail = false

   if (s > 0) && (s < 3) {

//   fmt.Println("checking to see if block ", i, " with count ",s," is available where the limit is ",b.Limit)

        k1 := make([]byte, 1)
	k1[0] = 255
	k2 := make([]byte, 1)
	k2[0] = 255


	if i > b.Limit { return false }

	_, _ = b.Store.Seek(int64(i*256), 0)
	_, _ = b.Store.Read(k1)

	if (s == 2) && (k1[0]==0) {
		if i+1 > b.Limit { return false }
		_, _ = b.Store.Seek(int64((i+1)*256), 0)
		_, _ = b.Store.Read(k2)
		if k2[0] == 0 {avail = true}
	}else{
		if k1[0] == 0 {avail = true}
	}
   }

   return avail
}

func (b Base) ReadKnodBlock (i uint64 ) (*Knod, error) {

	if i > b.Limit { return nil, errors.New("attempt to read beyond end of store") }

	_, err := b.Store.Seek(int64(i*256), 0)
	if err != nil {
		fmt.Println(b.Store)
		return nil, err
	}

	k := Knod{}

	e := binary.Read(b.Store, binary.LittleEndian, &k)

	return &k, e
}

func (b Base) WriteKnodBlock (k * Knod, i uint64 ) error {

	if i > b.Limit { errors.New("attempt to write beyond end of store") }

	_, err := b.Store.Seek(int64(i*256), 0)
	if err != nil {
	   return err
	}

	e := binary.Write(b.Store, binary.LittleEndian, k)

	return e
}

func (b Base) ReadBodyBlock (i uint64 ) (*Body, error) {

	if i > b.Limit { return nil, errors.New("attempt to read beyond end of store") }

	_, err := b.Store.Seek(int64(i*256), 0)
	if err != nil {
	   return nil, err
	}

	kb := Body{}

	e := binary.Read(b.Store, binary.LittleEndian, &kb)

	return &kb, e
}

func (b Base) WriteBodyBlock (kb * Body, i uint64 ) error {

	if i > b.Limit { errors.New("attempt to write beyond end of store") }

	_, err := b.Store.Seek(int64(i*256), 0)
	if err != nil {
	   return err
	}

	e := binary.Write(b.Store, binary.LittleEndian, kb)

	return e
}

