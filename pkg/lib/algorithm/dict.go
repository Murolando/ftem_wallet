package algorithm

import (
	_ "embed"
	"strings"
)

// BIP39EnglishWords содержит встроенный список из 2048 английских слов BIP-39
//
//go:embed bip39_english.txt
var BIP39EnglishWords string

// Dictionary представляет собой словарь BIP-39
type Dictionary struct {
	words   []string
	wordMap map[string]int // слово -> индекс
}

// NewDictionary создает новый словарь BIP-39
func NewDictionary() *Dictionary {
	words := strings.Split(strings.TrimSpace(BIP39EnglishWords), "\n")
	wordMap := make(map[string]int, len(words))

	for i, word := range words {
		wordMap[strings.TrimSpace(word)] = i
	}

	return &Dictionary{
		words:   words,
		wordMap: wordMap,
	}
}

// GetWord возвращает слово по индексу (0-2047)
func (d *Dictionary) GetWord(index int) (string, bool) {
	if index < 0 || index >= len(d.words) {
		return "", false
	}
	return d.words[index], true
}

// GetIndex возвращает индекс слова в словаре
func (d *Dictionary) GetIndex(word string) (int, bool) {
	index, exists := d.wordMap[word]
	return index, exists
}

// GetAllWords возвращает все слова словаря
func (d *Dictionary) GetAllWords() []string {
	result := make([]string, len(d.words))
	copy(result, d.words)
	return result
}

// Size возвращает размер словаря (должен быть 2048 для BIP-39)
func (d *Dictionary) Size() int {
	return len(d.words)
}

// IsValidWord проверяет, есть ли слово в словаре
func (d *Dictionary) IsValidWord(word string) bool {
	_, exists := d.wordMap[word]
	return exists
}
