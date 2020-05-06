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
const POS_MOVE_HASH_SIZE = 1 << PV_HASH_KEY_SIZE_IN_BITS
const POS_MOVE_HASH_MASK = PV_HASH_SIZE - 1

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

const PV_HASH_KEY_SIZE_IN_BITS = 20
const PV_HASH_SIZE = 1 << PV_HASH_KEY_SIZE_IN_BITS
const PV_HASH_MASK = PV_HASH_SIZE - 1

const MAX_PV_MOVES = 4

type PvEntry struct{		
	Depth           int8
	Zobrist         uint64
	Moves           [MAX_PV_MOVES]Move
}

type PvHash struct{
	Entries         [PV_HASH_SIZE]PvEntry
}

func (pvh *PvHash) Get(zobrist uint64) (uint32, PvEntry, bool){
	key := uint32(zobrist & PV_HASH_MASK)

	entry := pvh.Entries[key]

	return key, entry, ( entry.Zobrist == zobrist ) && ( entry.Depth < INFINITE_DEPTH )
}

func (pvh *PvHash) Set(zobrist uint64, pve PvEntry){
	pve.Zobrist = zobrist

	key, oldPve, _ := pvh.Get(zobrist)

	if pve.Depth <= oldPve.Depth{
		pvh.Entries[key] = pve
	}
}
