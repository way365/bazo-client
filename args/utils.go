package args

import (
	"crypto/ecdsa"
	"github.com/julwil/bazo-miner/crypto"
	"github.com/julwil/bazo-miner/protocol"
	"strings"
)

// Resolves a public key from a provided string.
// The string can either hold the public key directly
// or a filename pointing to the key file.
func ResolvePublicKey(publicKeyOrFilename string) (publicKey *ecdsa.PublicKey, err error) {
	if len(publicKeyOrFilename) == 0 {
		return nil, err
	}

	// Public key provided indirectly by a filename.
	if len(publicKeyOrFilename) > 0 && strings.Contains(publicKeyOrFilename, ".txt") {
		publicKey, err = crypto.GetOrCreateECDSAPublicKeyFromFile(publicKeyOrFilename)
		if err != nil {
			return publicKey, nil
		}
	}

	// In this case, the key is provided directly.
	// Split the string by \n or \t
	params := strings.Fields(publicKeyOrFilename)
	var size = protocol.ACCOUNT_ADDRESS_SIZE

	// Key is provided as one consecutive string, no \n delimiter.
	if len(params) == 1 && len(publicKeyOrFilename) == (2*size) || len(publicKeyOrFilename) == (3*size) {
		params = make([]string, 2)
		params[0] = publicKeyOrFilename[:size]
		params[1] = publicKeyOrFilename[size : 2*size]
	}

	// Public key provided directly as string.
	// Either including the private key (3 * size) or without.
	if len(params) == 2 || len(params) == 3 {
		publicKey, err = crypto.GetPubKeyFromString(params[0], params[1])
		if err != nil {
			return publicKey, err
		}
	}

	return publicKey, nil
}

// Resolves a private key from a provided string.
// The string can either hold the private key directly
// or a filename pointing to the key file.
func ResolvePrivateKey(privateKeyOrFilename string) (privateKey *ecdsa.PrivateKey, err error) {
	if len(privateKeyOrFilename) == 0 {
		return nil, err
	}

	// Public key provided indirectly by a filename.
	if len(privateKeyOrFilename) > 0 && strings.Contains(privateKeyOrFilename, ".txt") {
		privateKey, err = crypto.ExtractECDSAKeyFromFile(privateKeyOrFilename)
		if err != nil {
			return privateKey, nil
		}
	}

	// In this case, the key is provided directly.
	// Split the string by \n or \t
	params := strings.Fields(privateKeyOrFilename)
	var size = protocol.ACCOUNT_ADDRESS_SIZE

	// Key is provided as one consecutive string, no \n delimiter.
	if len(params) == 1 && len(privateKeyOrFilename) == 3*size {
		params = make([]string, 3)
		params[0] = privateKeyOrFilename[:size]
		params[1] = privateKeyOrFilename[size : 2*size]
		params[2] = privateKeyOrFilename[2*size : 3*size]
	}

	// Private key provided directly. It must contain X, Y, D
	if len(params) == 3 {
		privateKey, err = crypto.GetPrivKeyFromString(params[0], params[1], params[2])
		if err != nil {
			return privateKey, err
		}
	}

	return privateKey, nil
}

// Resolves a set of chameleon hash parameters from a provided string.
// The string can either hold the parameters directly
// or a filename pointing to the ch-parmas file.
func ResolveChParams(chParamsOrFilename string) (chParams *crypto.ChameleonHashParameters, err error) {
	if len(chParamsOrFilename) == 0 {
		return nil, err
	}

	// CH params provided indirectly by a filename.
	if len(chParamsOrFilename) > 0 && strings.Contains(chParamsOrFilename, ".txt") {
		chParams, err = crypto.GetOrCreateChParamsFromFile(chParamsOrFilename)
		if err != nil {
			return nil, err
		}
	}

	// In this case, the parameters are provided directly.
	// Split the string by \n or \t
	params := strings.Fields(chParamsOrFilename)
	var g, p, q, hk, tk string
	var size = crypto.CH_PARAM_SIZE

	// Parameters are provided as one consecutive string, no \n delimiter.
	if len(params) == 1 && len(chParamsOrFilename) == 4*size || len(chParamsOrFilename) == 5*size {
		params = make([]string, 5)
		params[0] = chParamsOrFilename[:size]
		params[1] = chParamsOrFilename[size : 2*size]
		params[2] = chParamsOrFilename[2*size : 3*size]
		params[3] = chParamsOrFilename[3*size : 4*size]

		// If trapdoor is included.
		if len(chParamsOrFilename) == 5 {
			params[4] = chParamsOrFilename[4*size : 5*size]
		}
	}

	// Chameleon hash parameters with and w/o trapdoor key.
	if len(params) == 4 || len(params) == 5 {
		g = params[0]
		p = params[1]
		q = params[2]
		hk = params[3]

		// If trapdoor is included.
		if len(params) == 5 {
			tk = params[4]
		}

		chParams, err = crypto.GetChParamsFromString(g, p, q, hk, tk)
		if err != nil {
			return nil, err
		}
	}

	return chParams, nil
}
