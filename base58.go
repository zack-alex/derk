package derk

import (
	"math/big"
)

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func base58Encode(data []byte) string {
	zeros := 0
	for zeros = 0; zeros < len(data); zeros++ {
		if data[zeros] != 0 {
			break
		}
	}

	var digits []int64

	var number big.Int
	number.SetBytes(data)
	var base big.Int
	base.SetInt64(58)
	for number.Cmp(big.NewInt(0)) != 0 {
		var mod big.Int
		number.DivMod(&number, &base, &mod)
		digits = append(digits, mod.Int64())
	}

	for i := 0; i < zeros; i++ {
		digits = append(digits, 0)
	}

	n := len(digits)
	res := make([]byte, n)
	for i, digit := range digits {
		res[n-1-i] = base58Alphabet[digit]
	}

	return string(res)
}
