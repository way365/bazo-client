package args

import (
	"crypto/ecdsa"
	"github.com/way365/bazo-miner/crypto"
	"github.com/way365/bazo-miner/protocol"
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
func ResolveParameters(parametersOrFilename string) (parameters *crypto.ChameleonHashParameters, err error) {
	if len(parametersOrFilename) == 0 {
		return nil, err
	}

	// CH params provided indirectly by a filename.
	if len(parametersOrFilename) > 0 && strings.Contains(parametersOrFilename, ".txt") {
		parameters, err = crypto.GetOrCreateParametersFromFile(parametersOrFilename)
		if err != nil {
			return nil, err
		}
	}

	// In this case, the parameters are provided directly.
	// Split the string by \n or \t
	lines := strings.Fields(parametersOrFilename)
	var g, p, q, hk, tk string
	var lenOfOneParameter = crypto.CH_PARAM_SIZE

	// Parameters are provided as one consecutive string, no \n delimiter.
	if len(lines) == 1 && len(parametersOrFilename) == 4*lenOfOneParameter || len(parametersOrFilename) == 5*lenOfOneParameter {
		lines = make([]string, 5)
		lines[0] = parametersOrFilename[:lenOfOneParameter]
		lines[1] = parametersOrFilename[lenOfOneParameter : 2*lenOfOneParameter]
		lines[2] = parametersOrFilename[2*lenOfOneParameter : 3*lenOfOneParameter]
		lines[3] = parametersOrFilename[3*lenOfOneParameter : 4*lenOfOneParameter]

		// If trapdoor is included.
		if len(parametersOrFilename) == 5 {
			lines[4] = parametersOrFilename[4*lenOfOneParameter : 5*lenOfOneParameter]
		}
	}

	// Chameleon hash parameters with and w/o trapdoor key.
	if len(lines) == 4 || len(lines) == 5 {
		g = lines[0]
		p = lines[1]
		q = lines[2]
		hk = lines[3]

		// If trapdoor is included.
		if len(lines) == 5 {
			tk = lines[4]
		}

		parameters, err = crypto.GetParametersFromString(g, p, q, hk, tk)
		if err != nil {
			return nil, err
		}
	}

	return parameters, nil
}
