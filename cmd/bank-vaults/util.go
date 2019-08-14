// Copyright © 2018 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"plugin"
	"strings"

	"github.com/banzaicloud/bank-vaults/pkg/kv"
	"github.com/banzaicloud/bank-vaults/pkg/kv/alibabakms"
	"github.com/banzaicloud/bank-vaults/pkg/kv/alibabaoss"
	"github.com/banzaicloud/bank-vaults/pkg/kv/awskms"
	"github.com/banzaicloud/bank-vaults/pkg/kv/azurekv"
	"github.com/banzaicloud/bank-vaults/pkg/kv/dev"
	"github.com/banzaicloud/bank-vaults/pkg/kv/file"
	"github.com/banzaicloud/bank-vaults/pkg/kv/gckms"
	"github.com/banzaicloud/bank-vaults/pkg/kv/gcs"
	"github.com/banzaicloud/bank-vaults/pkg/kv/k8s"
	"github.com/banzaicloud/bank-vaults/pkg/kv/s3"
	"github.com/banzaicloud/bank-vaults/pkg/sdk/vault"
	"github.com/spf13/viper"
)

func vaultConfigForConfig(cfg *viper.Viper) (vault.Config, error) {

	return vault.Config{
		SecretShares:    appConfig.GetInt(cfgSecretShares),
		SecretThreshold: appConfig.GetInt(cfgSecretThreshold),

		InitRootToken:  appConfig.GetString(cfgInitRootToken),
		StoreRootToken: appConfig.GetBool(cfgStoreRootToken),

		PreFlightChecks: appConfig.GetBool(cfgPreFlightChecks),
	}, nil
}

func kvStoreForConfig(cfg *viper.Viper) (kv.Service, error) {

	switch mode := cfg.GetString(cfgMode); mode {

	case cfgModeValueGoogleCloudKMSGCS:
		gcs, err := gcs.New(
			cfg.GetString(cfgGoogleCloudStorageBucket),
			cfg.GetString(cfgGoogleCloudStoragePrefix),
		)

		if err != nil {
			return nil, fmt.Errorf("error creating google cloud storage kv store: %s", err.Error())
		}

		kms, err := gckms.New(gcs,
			cfg.GetString(cfgGoogleCloudKMSProject),
			cfg.GetString(cfgGoogleCloudKMSLocation),
			cfg.GetString(cfgGoogleCloudKMSKeyRing),
			cfg.GetString(cfgGoogleCloudKMSCryptoKey),
		)

		if err != nil {
			return nil, fmt.Errorf("error creating google cloud kms kv store: %s", err.Error())
		}

		return kms, nil

	case cfgModeValueAWSKMS3:
		s3, err := s3.New(
			cfg.GetString(cfgAWSS3Region),
			cfg.GetString(cfgAWSS3Bucket),
			cfg.GetString(cfgAWSS3Prefix),
		)

		if err != nil {
			return nil, fmt.Errorf("error creating AWS S3 kv store: %s", err.Error())
		}

		kms, err := awskms.New(s3, cfg.GetString(cfgAWSKMSRegion), cfg.GetString(cfgAWSKMSKeyID))

		if err != nil {
			return nil, fmt.Errorf("error creating AWS KMS kv store: %s", err.Error())
		}

		return kms, nil

	case cfgModeValueAzureKeyVault:
		akv, err := azurekv.New(cfg.GetString(cfgAzureKeyVaultName))
		if err != nil {
			return nil, fmt.Errorf("error creating Azure Key Vault kv store: %s", err.Error())
		}

		return akv, nil

	case cfgModeValueAlibabaKMSOSS:
		accessKeyID := cfg.GetString(cfgAlibabaAccessKeyID)
		accessKeySecret := cfg.GetString(cfgAlibabaAccessKeySecret)

		if accessKeyID == "" || accessKeySecret == "" {
			return nil, fmt.Errorf("Alibaba accessKeyID or accessKeySecret can't be empty")
		}

		bucket := cfg.GetString(cfgAlibabaOSSBucket)

		if bucket == "" {
			return nil, fmt.Errorf("Alibaba OSS bucket should be specified")
		}

		oss, err := alibabaoss.New(
			cfg.GetString(cfgAlibabaOSSEndpoint),
			accessKeyID,
			accessKeySecret,
			bucket,
			cfg.GetString(cfgAlibabaOSSPrefix),
		)
		if err != nil {
			return nil, fmt.Errorf("error creating Alibaba OSS kv store: %s", err.Error())
		}

		kms, err := alibabakms.New(
			cfg.GetString(cfgAlibabaKMSRegion),
			accessKeyID,
			accessKeySecret,
			cfg.GetString(cfgAlibabaKMSKeyID),
			oss)
		if err != nil {
			return nil, fmt.Errorf("error creating Alibaba KMS kv store: %s", err.Error())
		}

		return kms, nil

	case cfgModeValueK8S:
		k8s, err := k8s.New(
			cfg.GetString(cfgK8SNamespace),
			cfg.GetString(cfgK8SSecret),
		)

		if err != nil {
			return nil, fmt.Errorf("error creating K8S Secret kv store: %s", err.Error())
		}

		return k8s, nil

	case cfgModeValueDev:
		dev, err := dev.New()
		if err != nil {
			return nil, fmt.Errorf("error creating Dev Secret kv store: %s", err.Error())
		}

		return dev, nil

	case cfgModeValueFile:
		file, err := file.New(cfg.GetString(cfgFilePath))
		if err != nil {
			return nil, fmt.Errorf("error creating File kv store: %s", err.Error())
		}

		return file, nil

	case cfgModeValuePlugin:

		// go build -buildmode=plugin -o mykv.so ./hack/mykv.go
		// bank-vaults configure --mode plugin --plugin-config path=mykv.so

		configRaw := cfg.GetStringSlice(cfgPluginConfig)

		// Hack until https://github.com/spf13/viper/issues/608 gets fixed
		config := map[string]string{}
		for _, c := range configRaw {
			s := strings.Split(c, "=")
			config[s[0]] = s[1]
		}

		path, ok := config["path"]
		if !ok {
			return nil, fmt.Errorf("plugin path is missing")
		}

		plug, err := plugin.Open(path)
		if err != nil {
			return nil, fmt.Errorf("unsupported plugin: '%s'", cfg.GetString(cfgMode))
		}

		New, err := plug.Lookup("New")
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		//call the NewService function
		return New.(func(map[string]string) (kv.Service, error))(config)

	default:
		return nil, fmt.Errorf("unsupported backend mode: '%s'", cfg.GetString(cfgMode))
	}
}
