// Handles the user's authentication
package account

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

const ECDSA uint8 = 1

type SmartcontractDeveloper struct {
	Address         string
	AccountType     uint8             // The cryptographic algorithm key
	EcdsaPublicKey  *ecdsa.PublicKey  // If the account type is ECDSA, then this one will keep the pub key
	EcdsaPrivateKey *ecdsa.PrivateKey //
}

// Creates a new SmartcontractDeveloper with a public key but without private key
func NewEcdsaPublicKey(pub_key *ecdsa.PublicKey) *SmartcontractDeveloper {
	return &SmartcontractDeveloper{
		Address:         crypto.PubkeyToAddress(*pub_key).Hex(),
		AccountType:     ECDSA,
		EcdsaPublicKey:  pub_key,
		EcdsaPrivateKey: nil,
	}
}

// Creates a new SmartcontractDeveloper with a private key
func NewEcdsaPrivateKey(private_key *ecdsa.PrivateKey) *SmartcontractDeveloper {
	return &SmartcontractDeveloper{
		Address:         crypto.PubkeyToAddress(private_key.PublicKey).Hex(),
		AccountType:     ECDSA,
		EcdsaPublicKey:  &private_key.PublicKey,
		EcdsaPrivateKey: private_key,
	}
}

// Encrypts the given data with a public key
// The result could be decrypted by the private key
//
// If the account has a private key, then the public key derived from it would be used
func (account *SmartcontractDeveloper) Encrypt(plain_text []byte) ([]byte, error) {
	if account.AccountType != ECDSA {
		return []byte{}, errors.New("only ECDSA protocol supported")
	}
	if account.EcdsaPrivateKey != nil {
		account.EcdsaPublicKey = &account.EcdsaPrivateKey.PublicKey
	}
	if account.EcdsaPublicKey == nil {
		return []byte{}, errors.New("the account has no public key")
	}
	// We get the public key for from the signature.
	// We also get the account address from public key using crypto.PubkeyToAddress(*ethereum_pb).Hex()
	// this account should be checked against the whitelisted accounts in the database.
	//
	// Encrypt the message with the public key
	curve_pb := ecies.ImportECDSAPublic(account.EcdsaPublicKey)
	cipher_text, err := ecies.Encrypt(rand.Reader, curve_pb, plain_text, nil, nil)

	return cipher_text, err
}

func (account *SmartcontractDeveloper) Decrypt(cipher_text []byte) ([]byte, error) {
	if account.AccountType != ECDSA {
		return []byte{}, errors.New("only ECDSA is supported")
	}

	if account.EcdsaPrivateKey == nil {
		return []byte{}, errors.New("the account has no private key")
	}

	curve_secret_key := ecies.ImportECDSA(account.EcdsaPrivateKey)
	plain_text, err := curve_secret_key.Decrypt(cipher_text, nil, nil)
	return plain_text, err
}
