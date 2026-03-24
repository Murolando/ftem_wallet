package entities

type ResultString string

const (
	ErrInvalidMnemonic     = ResultString("Error: Invalid mnemonic phrase")
	ErrPrivateKeyDerivation = ResultString("Error: Failed to derive private key")
)