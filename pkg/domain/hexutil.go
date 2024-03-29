package domain

// @see Reference: https://github.com/ethereum/go-ethereum/blob/master/common/hexutil/hexutil.go

import (
	"math/big"

	"github.com/pkg/errors"
)

const badNibble = ^uint64(0)

var bigWordNibbles int

func init() {
	b, _ := new(big.Int).SetString("FFFFFFFFFF", 16)
	switch len(b.Bits()) {
	case 1:
		bigWordNibbles = 16
	case 2:
		bigWordNibbles = 8
	default:
		panic("weird big.Word size")
	}
}

// Decode decodes a hex string with 0x prefix.
func hexToDecimalString(input string) (int64, error) {
	raw, err := checkNumber(input)
	if err != nil {
		return 0, err
	}
	if len(raw) > 64 {
		return 0, errors.New("hex number > 64 bits")
	}
	words := make([]big.Word, len(raw)/bigWordNibbles+1)
	end := len(raw)
	for i := range words {
		start := end - bigWordNibbles
		if start < 0 {
			start = 0
		}
		for ri := start; ri < end; ri++ {
			nib := decodeNibble(raw[ri])
			if nib == badNibble {
				return 0, errors.Errorf("invalid hex string")
			}
			words[i] *= 16
			words[i] += big.Word(nib)
		}
		end = start
	}
	dec := new(big.Int).SetBits(words)

	return dec.Int64(), nil
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func checkNumber(input string) (raw string, err error) {
	if len(input) == 0 {
		return "", errors.New("empty string")
	}
	if !has0xPrefix(input) {
		return "", errors.New("missing 0x prefix")
	}
	input = input[2:]
	if len(input) == 0 {
		return "", errors.New("hex string 0x")
	}
	if len(input) > 1 && input[0] == '0' {
		return "", errors.New("hex number with leading zero digits")
	}
	return input, nil
}

func decodeNibble(in byte) uint64 {
	switch {
	case in >= '0' && in <= '9':
		return uint64(in - '0')
	case in >= 'A' && in <= 'F':
		return uint64(in - 'A' + 10)
	case in >= 'a' && in <= 'f':
		return uint64(in - 'a' + 10)
	default:
		return badNibble
	}
}
