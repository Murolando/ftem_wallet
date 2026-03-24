package algorithm

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"slices"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/text/unicode/norm"
)

func BIP39Mnemonic() [12]string {
	var resultWords [12]string
	// создать словарь
	d := NewDictionary()

	// случайно сгенерировать 128 бит из 0 и 1
	var b [16]byte // 128 bits
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}

	out := make([]byte, 128)
	k := 0
	for i := 0; i < 16; i++ {
		// от старшего бита к младшему
		for bit := 7; bit >= 0; bit-- {
			out[k] = (b[i] >> uint(bit)) & 1
			k++
		}
	}
	fmt.Println("out", out, len(out))

	// посчитать 4 бита чек суммы от sha256
	cs := sha256.Sum256(b[:])
	fmt.Println("checksum", cs)

	// взять первые 4 бита
	csbytes := fmt.Sprintf("%08b", cs[0])
	fmt.Println("csbytes", csbytes)

	outcs := append(out, csbytes[0]-'0', csbytes[1]-'0', csbytes[2]-'0', csbytes[3]-'0')

	fmt.Println("outcs", outcs, len(outcs)) // 132 бита

	intWords := make([]int, 0)

	// разбить на групки по 11 бит и перевести их в десятичные числа
	for i := 0; i < 132; i += 11 {
		group := outcs[i : i+11]
		koef := 1
		num := 0
		for j := 10; j >= 0; j-- {
			num += int(group[j]) * koef
			koef *= 2
		}
		intWords = append(intWords, num)
	}
	fmt.Println("intWords", intWords)

	// перевести десятичные числа в слова
	for i := 0; i < 12; i++ {
		resultWords[i], _ = d.GetWord(intWords[i])
	}
	fmt.Println("resultWords", resultWords)
	// вернуть список слов
	return resultWords
}

func BIP39IsValidMnemomic(words [12]string) bool {
	d := NewDictionary()

	// перевести все слова в десятичные числа
	indexWords := make([]int, 0)
	for _, word := range words {
		index, ok := d.GetIndex(word)
		if !ok {
			return false // слово не найдено в словаре
		}
		indexWords = append(indexWords, index)
	}

	// перевести каждое число в 11 бит и собрать в общую последовательность
	allBits := make([]byte, 132) // 12 слов × 11 бит = 132 бита
	bitIndex := 0

	for _, index := range indexWords {
		// преобразуем индекс в 11 бит (от старшего к младшему)
		for bit := 10; bit >= 0; bit-- {
			allBits[bitIndex] = byte(index >> uint(bit) & 1)
			bitIndex++
		}
	}

	// взять первые 128 бит (энтропия) и перевести их в байты
	entropyBytes := make([]byte, 16) // 128 бит = 16 байт
	for i := 0; i < 16; i++ {
		var b byte = 0
		for bit := 0; bit < 8; bit++ {
			if allBits[i*8+bit] == 1 {
				b |= 1 << uint(7-bit)
			}
		}
		entropyBytes[i] = b
	}

	// прогнать байты через sha256
	hash := sha256.Sum256(entropyBytes)

	// взять первые 4 бита
	csbytes := fmt.Sprintf("%08b", hash[0])

	var expectedChecksum []byte
	expectedChecksum = append(expectedChecksum, csbytes[0]-'0', csbytes[1]-'0', csbytes[2]-'0', csbytes[3]-'0')

	return slices.Equal(allBits[128:132], expectedChecksum)
}

func BIP39SeedFromMnemomic(words [12]string, password string) []byte {
	// Объединяем слова в строку через пробелы
	mnemonic := strings.Join(words[:], " ")

	// Нормализуем по стандарту NFKD
	normalizedMnemonic := norm.NFKD.String(mnemonic)
	normalizedPassword := norm.NFKD.String(password)

	// Создаем соль: "mnemonic" + нормализованный пароль
	salt := "mnemonic" + normalizedPassword

	// Применяем PBKDF2 с SHA-512
	seed := pbkdf2.Key([]byte(normalizedMnemonic), []byte(salt), 2048, 64, sha512.New)

	return seed
}
