package main

import (
	"log"
	"math"

	"github.com/Workiva/go-datastructures/bitarray"
	"github.com/spaolacci/murmur3"
)

// ref: https://hur.st/bloomfilter/
// test whether an element is a member of set or not
// False postives matches are possible but false negative matches are not

// n = ceil(m / (-k / log(1 - exp(log(p) / k)))) Number of items in the filter
// p = pow(1 - exp(-k / (m / n)), k) Probability of false positives, fraction between 0 and 1 or a number indicating 1-in-p
// m = ceil((n * log(p)) / log(1 / pow(2, log(2)))); Number of bits in the filter
// k = round((m / n) * log(2)); Number of hash functions

type BloomFilter struct {
	bitArray     bitarray.BitArray
	bitArraysize uint64
	hashCount    uint64
}

// n: no of items to store in Bloom filter
// probability: it is the error rate of the bloom filter
// depending on the error rate the size of the bit array filter is determined
func NewBloomFilter(n uint64, probability float64) *BloomFilter {
	size := getBitsInFilter(n, probability)
	hashCount := getHashCount(size, n)
	bitArray := bitarray.NewBitArray(size)
	return &BloomFilter{
		bitArray:     bitArray,
		bitArraysize: size,
		hashCount:    hashCount,
	}
}

func (bf *BloomFilter) Add(element string) {
	for i := uint64(0); i < bf.hashCount; i++ {
		hashValue := murmur3.Sum32WithSeed([]byte(element), uint32(i)) % uint32(bf.bitArraysize)
		bf.bitArray.SetBit(uint64(hashValue))
	}
}

func (bf *BloomFilter) Check(element string) bool {
	for i := uint64(0); i < bf.hashCount; i++ {
		hashValue := murmur3.Sum32WithSeed([]byte(element), uint32(i)) % uint32(bf.bitArraysize)
		bit, err := bf.bitArray.GetBit(uint64(hashValue))
		if err != nil {
			log.Fatal(err)
		}
		if !bit {
			return false
		}
	}
	return true
}

func getBitsInFilter(n uint64, p float64) uint64 {
	m := uint64(math.Ceil((float64(n) * math.Log(p)) / math.Log(1/math.Pow(2, math.Log(2)))))
	return m
}

func getHashCount(m uint64, n uint64) uint64 {
	k := int(math.Round((float64(m) / float64(n)) * math.Log(2)))
	return uint64(k)
}
