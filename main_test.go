package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
)

func TestBloomFilter_AddAndCheck(t *testing.T) {
	bloomFilter := NewBloomFilter(100, 0.01)

	// Test 1: Add 25 strings to the Bloom filter and check for their presence
	stringsToAdd := []string{
		"apple",
		"banana",
		"cherry",
		"date",
		"elderberry",
		"fig",
		"grape",
		"honeydew",
		"imbe",
		"jackfruit",
		"kiwi",
		"lemon",
		"mango",
		"nectarine",
		"orange",
		"papaya",
		"quince",
		"raspberry",
		"strawberry",
		"tangerine",
		"ugli",
		"vanilla",
		"watermelon",
		"xylocarp",
		"yuzu",
	}

	for _, str := range stringsToAdd {
		bloomFilter.Add(str)
		if !bloomFilter.Check(str) {
			t.Errorf("Expected string %s to be present in the Bloom filter, but it is not", str)
		}
	}

	// Test 2: Check for strings that are not present in the Bloom filter
	stringsToCheck := []string{
		"apricot",
		"blueberry",
		"coconut",
		"durian",
		"elderflower",
		"guava",
		"huckleberry",
		"indian gooseberry",
		"jujube",
		"kiwifruit",
		"lime",
		"mulberry",
		"nectar",
		"olive",
		"pineapple",
		"quenepa",
		"rambutan",
		"strudel",
		"tangelo",
		"ugni",
		"vanillin",
		"wax apple",
		"xigua",
		"yellow passionfruit",
		"zinfandel grape",
	}

	for _, str := range stringsToCheck {
		if bloomFilter.Check(str) {
			t.Errorf("Expected string %s to be absent in the Bloom filter, but it is present", str)
		}
	}
}
func TestBloomFilter_FalsePositiveRate(t *testing.T) {
	// Create a Bloom filter with a small size to increase the chance of false positives
	bloomFilter := NewBloomFilter(10, 0.01)

	// Add an element to the Bloom filter
	element := "example"
	bloomFilter.Add(element)

	// Generate 100 non-added elements and check if they are reported as false positives
	falsePositiveCount := 0
	for i := 0; i < 100; i++ {
		nonAddedElement := fmt.Sprintf("non-added-%d", i)
		if bloomFilter.Check(nonAddedElement) {
			falsePositiveCount++
		}
	}

	// Calculate the false positive rate and check if it's within an acceptable range
	falsePositiveRate := float64(falsePositiveCount) / 100.0
	acceptableError := 0.05 // 5% acceptable error
	if math.Abs(falsePositiveRate-0.01) > acceptableError {
		t.Errorf("False positive rate is %f, expected around 0.01", falsePositiveRate)
	}
}

func TestBloomFilterFromFile(t *testing.T) {
	// Open the file
	file, err := os.Open("shakespeare_unique.txt")
	if err != nil {
		t.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	var words []string
	for scanner.Scan() {
		line := scanner.Text()
		words = append(words, strings.Split(line, ",")...)
	}

	// Create a Bloom filter
	bloomFilter := NewBloomFilter(uint64(len(words)/2), 0.0000001)

	// Add half of the words to the Bloom filter
	for i := 0; i < len(words)/2; i++ {
		bloomFilter.Add(words[i])
	}

	// Check if the remaining words are present in the Bloom filter
	for i := len(words) / 2; i < len(words); i++ {
		if bloomFilter.Check(words[i]) {
			t.Errorf("Expected word %s to be absent in the Bloom filter, but it is present", words[i])
		}
	}
}
