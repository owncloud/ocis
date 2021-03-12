package command

import (
	"path/filepath"

	"github.com/owncloud/ocis/storage/pkg/config"
)

func drivers(cfg *config.Config) map[string]interface{} {

	// UploadInfoDir must always be absolute for REVA
	uploadInfoDir, _ := filepath.Abs(cfg.Reva.Storages.OwnCloud.UploadInfoDir)

	return map[string]interface{}{
		"eos": map[string]interface{}{
			"namespace":              cfg.Reva.Storages.EOS.Root,
			"shadow_namespace":       cfg.Reva.Storages.EOS.ShadowNamespace,
			"uploads_namespace":      cfg.Reva.Storages.EOS.UploadsNamespace,
			"share_folder":           cfg.Reva.Storages.EOS.ShareFolder,
			"eos_binary":             cfg.Reva.Storages.EOS.EosBinary,
			"xrdcopy_binary":         cfg.Reva.Storages.EOS.XrdcopyBinary,
			"master_url":             cfg.Reva.Storages.EOS.MasterURL,
			"slave_url":              cfg.Reva.Storages.EOS.SlaveURL,
			"cache_directory":        cfg.Reva.Storages.EOS.CacheDirectory,
			"sec_protocol":           cfg.Reva.Storages.EOS.SecProtocol,
			"keytab":                 cfg.Reva.Storages.EOS.Keytab,
			"single_username":        cfg.Reva.Storages.EOS.SingleUsername,
			"enable_logging":         cfg.Reva.Storages.EOS.EnableLogging,
			"show_hidden_sys_files":  cfg.Reva.Storages.EOS.ShowHiddenSysFiles,
			"force_single_user_mode": cfg.Reva.Storages.EOS.ForceSingleUserMode,
			"use_keytab":             cfg.Reva.Storages.EOS.UseKeytab,
			"gatewaysvc":             cfg.Reva.Storages.EOS.GatewaySVC,
		},
		"eoshome": map[string]interface{}{
			"namespace":              cfg.Reva.Storages.EOS.Root,
			"shadow_namespace":       cfg.Reva.Storages.EOS.ShadowNamespace,
			"uploads_namespace":      cfg.Reva.Storages.EOS.UploadsNamespace,
			"share_folder":           cfg.Reva.Storages.EOS.ShareFolder,
			"eos_binary":             cfg.Reva.Storages.EOS.EosBinary,
			"xrdcopy_binary":         cfg.Reva.Storages.EOS.XrdcopyBinary,
			"master_url":             cfg.Reva.Storages.EOS.MasterURL,
			"slave_url":              cfg.Reva.Storages.EOS.SlaveURL,
			"cache_directory":        cfg.Reva.Storages.EOS.CacheDirectory,
			"sec_protocol":           cfg.Reva.Storages.EOS.SecProtocol,
			"keytab":                 cfg.Reva.Storages.EOS.Keytab,
			"single_username":        cfg.Reva.Storages.EOS.SingleUsername,
			"user_layout":            cfg.Reva.Storages.EOS.UserLayout,
			"enable_logging":         cfg.Reva.Storages.EOS.EnableLogging,
			"show_hidden_sys_files":  cfg.Reva.Storages.EOS.ShowHiddenSysFiles,
			"force_single_user_mode": cfg.Reva.Storages.EOS.ForceSingleUserMode,
			"use_keytab":             cfg.Reva.Storages.EOS.UseKeytab,
			"gatewaysvc":             cfg.Reva.Storages.EOS.GatewaySVC,
		},
		"eosgrpc": map[string]interface{}{
			"namespace":              cfg.Reva.Storages.EOS.Root,
			"shadow_namespace":       cfg.Reva.Storages.EOS.ShadowNamespace,
			"share_folder":           cfg.Reva.Storages.EOS.ShareFolder,
			"eos_binary":             cfg.Reva.Storages.EOS.EosBinary,
			"xrdcopy_binary":         cfg.Reva.Storages.EOS.XrdcopyBinary,
			"master_url":             cfg.Reva.Storages.EOS.MasterURL,
			"master_grpc_uri":        cfg.Reva.Storages.EOS.GrpcURI,
			"slave_url":              cfg.Reva.Storages.EOS.SlaveURL,
			"cache_directory":        cfg.Reva.Storages.EOS.CacheDirectory,
			"sec_protocol":           cfg.Reva.Storages.EOS.SecProtocol,
			"keytab":                 cfg.Reva.Storages.EOS.Keytab,
			"single_username":        cfg.Reva.Storages.EOS.SingleUsername,
			"user_layout":            cfg.Reva.Storages.EOS.UserLayout,
			"enable_logging":         cfg.Reva.Storages.EOS.EnableLogging,
			"show_hidden_sys_files":  cfg.Reva.Storages.EOS.ShowHiddenSysFiles,
			"force_single_user_mode": cfg.Reva.Storages.EOS.ForceSingleUserMode,
			"use_keytab":             cfg.Reva.Storages.EOS.UseKeytab,
			"enable_home":            cfg.Reva.Storages.EOS.EnableHome,
			"gatewaysvc":             cfg.Reva.Storages.EOS.GatewaySVC,
		},
		"local": map[string]interface{}{
			"root":         cfg.Reva.Storages.Local.Root,
			"share_folder": cfg.Reva.Storages.Local.ShareFolder,
		},
		"localhome": map[string]interface{}{
			"root":         cfg.Reva.Storages.Local.Root,
			"share_folder": cfg.Reva.Storages.Local.ShareFolder,
			"user_layout":  cfg.Reva.Storages.Local.UserLayout,
		},
		"owncloud": map[string]interface{}{
			"datadirectory":   cfg.Reva.Storages.OwnCloud.Root,
			"upload_info_dir": uploadInfoDir,
			"sharedirectory":  cfg.Reva.Storages.OwnCloud.ShareFolder,
			"user_layout":     cfg.Reva.Storages.OwnCloud.UserLayout,
			"redis":           cfg.Reva.Storages.OwnCloud.Redis,
			"enable_home":     cfg.Reva.Storages.OwnCloud.EnableHome,
			"scan":            cfg.Reva.Storages.OwnCloud.Scan,
			"userprovidersvc": cfg.Reva.Users.Endpoint,
		},
		"ocis": map[string]interface{}{
			"root":                cfg.Reva.Storages.Common.Root,
			"enable_home":         cfg.Reva.Storages.Common.EnableHome,
			"user_layout":         cfg.Reva.Storages.Common.UserLayout,
			"treetime_accounting": true,
			"treesize_accounting": true,
			"owner":               "95cb8724-03b2-11eb-a0a6-c33ef8ef53ad", // the accounts service system account uuid
		},
		"s3": map[string]interface{}{
			"region":     cfg.Reva.Storages.S3.Region,
			"access_key": cfg.Reva.Storages.S3.AccessKey,
			"secret_key": cfg.Reva.Storages.S3.SecretKey,
			"endpoint":   cfg.Reva.Storages.S3.Endpoint,
			"bucket":     cfg.Reva.Storages.S3.Bucket,
			"prefix":     cfg.Reva.Storages.S3.Root,
		},
	}
}
