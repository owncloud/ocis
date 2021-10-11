package storagedrivers

import (
	"github.com/owncloud/ocis/storage/pkg/config"
)

func MetadataDrivers(cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"eos": map[string]interface{}{
			"namespace":              cfg.Reva.MetadataStorage.EOS.Root,
			"shadow_namespace":       cfg.Reva.MetadataStorage.EOS.ShadowNamespace,
			"uploads_namespace":      cfg.Reva.MetadataStorage.EOS.UploadsNamespace,
			"eos_binary":             cfg.Reva.MetadataStorage.EOS.EosBinary,
			"xrdcopy_binary":         cfg.Reva.MetadataStorage.EOS.XrdcopyBinary,
			"master_url":             cfg.Reva.MetadataStorage.EOS.MasterURL,
			"slave_url":              cfg.Reva.MetadataStorage.EOS.SlaveURL,
			"cache_directory":        cfg.Reva.MetadataStorage.EOS.CacheDirectory,
			"sec_protocol":           cfg.Reva.MetadataStorage.EOS.SecProtocol,
			"keytab":                 cfg.Reva.MetadataStorage.EOS.Keytab,
			"single_username":        cfg.Reva.MetadataStorage.EOS.SingleUsername,
			"enable_logging":         cfg.Reva.MetadataStorage.EOS.EnableLogging,
			"show_hidden_sys_files":  cfg.Reva.MetadataStorage.EOS.ShowHiddenSysFiles,
			"force_single_user_mode": cfg.Reva.MetadataStorage.EOS.ForceSingleUserMode,
			"use_keytab":             cfg.Reva.MetadataStorage.EOS.UseKeytab,
			"gatewaysvc":             cfg.Reva.MetadataStorage.EOS.GatewaySVC,
			"enable_home":            false,
		},
		"eosgrpc": map[string]interface{}{
			"namespace":              cfg.Reva.MetadataStorage.EOS.Root,
			"shadow_namespace":       cfg.Reva.MetadataStorage.EOS.ShadowNamespace,
			"eos_binary":             cfg.Reva.MetadataStorage.EOS.EosBinary,
			"xrdcopy_binary":         cfg.Reva.MetadataStorage.EOS.XrdcopyBinary,
			"master_url":             cfg.Reva.MetadataStorage.EOS.MasterURL,
			"master_grpc_uri":        cfg.Reva.MetadataStorage.EOS.GrpcURI,
			"slave_url":              cfg.Reva.MetadataStorage.EOS.SlaveURL,
			"cache_directory":        cfg.Reva.MetadataStorage.EOS.CacheDirectory,
			"sec_protocol":           cfg.Reva.MetadataStorage.EOS.SecProtocol,
			"keytab":                 cfg.Reva.MetadataStorage.EOS.Keytab,
			"single_username":        cfg.Reva.MetadataStorage.EOS.SingleUsername,
			"user_layout":            cfg.Reva.MetadataStorage.EOS.UserLayout,
			"enable_logging":         cfg.Reva.MetadataStorage.EOS.EnableLogging,
			"show_hidden_sys_files":  cfg.Reva.MetadataStorage.EOS.ShowHiddenSysFiles,
			"force_single_user_mode": cfg.Reva.MetadataStorage.EOS.ForceSingleUserMode,
			"use_keytab":             cfg.Reva.MetadataStorage.EOS.UseKeytab,
			"enable_home":            false,
			"gatewaysvc":             cfg.Reva.MetadataStorage.EOS.GatewaySVC,
		},
		"local": map[string]interface{}{
			"root": cfg.Reva.MetadataStorage.Local.Root,
		},
		"ocis": map[string]interface{}{
			"root":                cfg.Reva.MetadataStorage.OCIS.Root,
			"enable_home":         false,
			"user_layout":         cfg.Reva.MetadataStorage.OCIS.UserLayout,
			"treetime_accounting": false,
			"treesize_accounting": false,
			"owner":               cfg.Reva.MetadataStorage.OCIS.ServiceUserUUID, // the accounts service system account uuid
		},
		"s3": map[string]interface{}{
			"region":     cfg.Reva.MetadataStorage.S3.Region,
			"access_key": cfg.Reva.MetadataStorage.S3.AccessKey,
			"secret_key": cfg.Reva.MetadataStorage.S3.SecretKey,
			"endpoint":   cfg.Reva.MetadataStorage.S3.Endpoint,
			"bucket":     cfg.Reva.MetadataStorage.S3.Bucket,
		},
		"s3ng": map[string]interface{}{
			"root":          cfg.Reva.MetadataStorage.S3NG.Root,
			"enable_home":   false,
			"user_layout":   cfg.Reva.MetadataStorage.S3NG.UserLayout,
			"s3.region":     cfg.Reva.MetadataStorage.S3NG.Region,
			"s3.access_key": cfg.Reva.MetadataStorage.S3NG.AccessKey,
			"s3.secret_key": cfg.Reva.MetadataStorage.S3NG.SecretKey,
			"s3.endpoint":   cfg.Reva.MetadataStorage.S3NG.Endpoint,
			"s3.bucket":     cfg.Reva.MetadataStorage.S3NG.Bucket,
		},
	}
}
