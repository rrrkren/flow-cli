/*
 * Flow CLI
 *
 * Copyright 2019 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package keys

import (
	"context"
	"fmt"

	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/spf13/cobra"

	"github.com/onflow/flow-cli/flowkit"
	"github.com/onflow/flow-cli/flowkit/output"
	"github.com/onflow/flow-cli/internal/command"
)

type flagsGenerate struct {
	Mnemonic       string `flag:"mnemonic" info:"Mnemonic seed to use"`
	DerivationPath string `default:"m/44'/539'/0'/0/0" flag:"derivationPath" info:"Derivation path"`
	KeySigAlgo     string `default:"ECDSA_P256" flag:"sig-algo" info:"Signature algorithm"`
}

var generateFlags = flagsGenerate{}

var generateCommand = &command.Command{
	Cmd: &cobra.Command{
		Use:     "generate",
		Short:   "Generate a new key-pair",
		Example: "flow keys generate",
	},
	Flags: &generateFlags,
	Run:   generate,
}

func generate(
	_ []string,
	_ command.GlobalFlags,
	_ output.Logger,
	_ flowkit.ReaderWriter,
	flow flowkit.Services,
) (command.Result, error) {
	sigAlgo := crypto.StringToSignatureAlgorithm(generateFlags.KeySigAlgo)
	if sigAlgo == crypto.UnknownSignatureAlgorithm {
		return nil, fmt.Errorf("invalid signature algorithm: %s", generateFlags.KeySigAlgo)
	}

	var err error
	mnemonic := generateFlags.Mnemonic
	if mnemonic == "" {
		_, mnemonic, err = flow.GenerateMnemonicKey(context.Background(), sigAlgo, generateFlags.DerivationPath)
		if err != nil {
			return nil, err
		}
	}

	privateKey, err := flow.DerivePrivateKeyFromMnemonic(
		context.Background(),
		mnemonic,
		sigAlgo,
		generateFlags.DerivationPath,
	)
	if err != nil {
		return nil, err
	}

	return &keyResult{
		privateKey:     privateKey,
		publicKey:      privateKey.PublicKey(),
		sigAlgo:        sigAlgo,
		mnemonic:       mnemonic,
		derivationPath: generateFlags.DerivationPath,
	}, nil
}
