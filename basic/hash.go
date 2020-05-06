package basic

import (
	"github.com/spaolacci/murmur3"
)

//https://stackoverflow.com/questions/35371385/how-can-i-convert-an-int64-into-a-byte-array-in-go
func Uint64ToByteSlice(v uint64) []byte{
	return []byte{
		byte(0xff & v),
        byte(0xff & (v >> 8)),
        byte(0xff & (v >> 16)),
        byte(0xff & (v >> 24)),
        byte(0xff & (v >> 32)),
        byte(0xff & (v >> 40)),
        byte(0xff & (v >> 48)),
        byte(0xff & (v >> 56)),
	}
}

func Uint32ToByteSlice(v uint32) []byte{
	return []byte{
		byte(0xff & v),
        byte(0xff & (v >> 8)),
        byte(0xff & (v >> 16)),
        byte(0xff & (v >> 24)),        
	}
}

//https://github.com/spaolacci/murmur3/blob/539464a789e9b9f01bc857458ffe2c5c1a2ed382/murmur32.go#L106
func Mix32Uint64AndUint32(v64 uint64, v32 uint32) uint32{
	return murmur3.Sum32(append(Uint64ToByteSlice(v64), Uint32ToByteSlice(v32)...))
}

const POS_MOVE_HASH_KEY_SIZE_IN_BITS = 25
const POS_MOVE_HASH_SIZE = 1 << POS_MOVE_HASH_KEY_SIZE_IN_BITS
const POS_MOVE_HASH_MASK = POS_MOVE_HASH_SIZE - 1

type PosMoveEntry struct{
	Used            bool
	Depth           int8
	SubTree         int
}

type PosMoveHash struct{
	Entries         [POS_MOVE_HASH_SIZE]PosMoveEntry
}

func (pmh *PosMoveHash) Get(zobrist uint64, move Move) (uint32, PosMoveEntry){
	key := Mix32Uint64AndUint32(zobrist, uint32(move)) & POS_MOVE_HASH_MASK

	return key, pmh.Entries[key]
}

func (pmh *PosMoveHash) Set(zobrist uint64, move Move, pme PosMoveEntry){
	pme.Used = true

	key, oldPme := pmh.Get(zobrist, move)

	if !oldPme.Used{
		pmh.Entries[key] = pme
		return
	}

	if pme.Depth > oldPme.Depth{
		return
	}

	pmh.Entries[key] = pme
}
