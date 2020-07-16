package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-store/pkg/config"
	"github.com/owncloud/ocis-store/pkg/proto/v0"
	"google.golang.org/protobuf/encoding/protojson"
)

// BleveDocument wraps the generated Record.Metadata and adds a property that is used to distinguish documents in the index.
type BleveDocument struct {
	Metadata map[string]*proto.Field
	Database string `json:"database"`
	Table    string `json:"table"`
}

// New returns a new instance of Service
func New(opts ...Option) (s *Service, err error) {
	options := newOptions(opts...)
	logger := options.Logger
	cfg := options.Config

	recordsDir := filepath.Join(cfg.Datapath, "databases")
	{
		var fi os.FileInfo
		if fi, err = os.Stat(recordsDir); err != nil {
			if os.IsNotExist(err) {
				// create store directory
				if err = os.MkdirAll(recordsDir, 0700); err != nil {
					return nil, err
				}
			}
		} else if !fi.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", recordsDir)
		}
	}

	indexMapping := bleve.NewIndexMapping()
	// keep all symbols in terms to allow exact matching, eg. emails
	indexMapping.DefaultAnalyzer = keyword.Name

	s = &Service{
		id:     cfg.GRPC.Namespace + "." + "store",
		log:    logger,
		Config: cfg,
	}

	indexDir := filepath.Join(cfg.Datapath, "index.bleve")
	// for now recreate index on every start
	if err = os.RemoveAll(indexDir); err != nil {
		return nil, err
	}
	if s.index, err = bleve.New(indexDir, indexMapping); err != nil {
		return
	}
	// if err = s.indexRecords(recordsDir); err != nil {
	// 	return nil, err
	// }
	return
}

// Service implements the AccountsServiceHandler interface
type Service struct {
	id     string
	log    log.Logger
	Config *config.Config
	index  bleve.Index
}

// Read implements the StoreHandler interface.
func (s *Service) Read(c context.Context, rreq *proto.ReadRequest, rres *proto.ReadResponse) error {
	if len(rreq.Key) != 0 {
		file := filepath.Join(s.Config.Datapath, "databases", rreq.Options.Database, rreq.Options.Table, rreq.Key)

		var data []byte
		rec := &proto.Record{}
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return merrors.NotFound(s.id, "could not read record")
		}

		if err = protojson.Unmarshal(data, rec); err != nil {
			return merrors.InternalServerError(s.id, "could not unmarshal record")
		}

		rres.Records = append(rres.Records, rec)
		return nil
	}

	s.log.Info().Interface("requeest", rreq).Msg("read request")
	if rreq.Options.Where != nil {
		// build bleve query
		// execute search
		// fetch the actual record if there's a hit
		dtq := bleve.NewTermQuery(rreq.Options.Database)
		ttq := bleve.NewTermQuery(rreq.Options.Table)
		dtq.SetField("database")
		ttq.SetField("table")

		query := bleve.NewConjunctionQuery(dtq, ttq)
		for k, v := range rreq.Options.Where {
			ntq := bleve.NewTermQuery(v.Value)
			ntq.SetField("metadata." + k + ".value")
			query.AddQuery(ntq)
		}

		searchRequest := bleve.NewSearchRequest(query)
		var searchResult *bleve.SearchResult
		searchResult, err := s.index.Search(searchRequest)
		if err != nil {
			s.log.Error().Err(err).Msg("could not execute bleve search")
			return merrors.InternalServerError(s.id, "could not execute bleve search: %v", err.Error())
		}

		for _, hit := range searchResult.Hits {
			rec := &proto.Record{}

			dest := filepath.Join(s.Config.Datapath, "databases", hit.ID)

			var data []byte
			data, err := ioutil.ReadFile(dest)
			s.log.Info().Str("path", dest).Interface("hit", hit).Msgf("hit info")
			if err != nil {
				s.log.Info().Str("path", dest).Interface("hit", hit).Msgf("file not found")
				return merrors.NotFound(s.id, "could not read record")
			}

			if err = protojson.Unmarshal(data, rec); err != nil {
				return merrors.InternalServerError(s.id, "could not unmarshal record")
			}

			rres.Records = append(rres.Records, rec)
		}
		return nil
	}

	return merrors.InternalServerError(s.id, "neither id nor metadata present")
}

// Write implements the StoreHandler interface.
func (s *Service) Write(c context.Context, wreq *proto.WriteRequest, wres *proto.WriteResponse) error {
	// TODO sanitize key. As it may contain invalid characters, such as slashes.
	// file: /var/tmp/ocis-store/databases/{database}/{table}/{record.key}.
	file := filepath.Join(s.Config.Datapath, "databases", wreq.Options.Database, wreq.Options.Table, wreq.Record.Key)

	var bytes []byte
	bytes, err := protojson.Marshal(wreq.Record)
	if err != nil {
		return merrors.InternalServerError(s.id, "could not marshal record")
	}

	err = os.MkdirAll(filepath.Dir(file), 0700)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, bytes, 0600)
	if err != nil {
		return merrors.InternalServerError(s.id, "could not write record")
	}

	doc := BleveDocument{
		Metadata: wreq.Record.Metadata,
		Database: wreq.Options.Database,
		Table:    wreq.Options.Table,
	}
	// TODO sanitize input.
	if err := s.index.Index(strings.Join([]string{wreq.Options.Database, wreq.Options.Table, wreq.Record.Key}, "/"), doc); err != nil {
		s.log.Error().Err(err).Interface("document", doc).Msg("could not index record metadata")
		return err
	}

	return nil
}

// Delete implements the StoreHandler interface.
func (s *Service) Delete(c context.Context, dreq *proto.DeleteRequest, dres *proto.DeleteResponse) error {
	file := filepath.Join(s.Config.Datapath, "databases", dreq.Options.Database, dreq.Options.Table, dreq.Key)
	if err := os.Remove(file); err != nil {
		if os.IsNotExist(err) {
			return merrors.NotFound(s.id, "could not find record")
		}

		return merrors.InternalServerError(s.id, "could not delete record")
	}
	return nil
}

// List implements the StoreHandler interface.
func (s *Service) List(context.Context, *proto.ListRequest, proto.Store_ListStream) error {
	return nil
}

// Databases implements the StoreHandler interface.
func (s *Service) Databases(c context.Context, dbreq *proto.DatabasesRequest, dbres *proto.DatabasesResponse) error {
	file := filepath.Join(s.Config.Datapath, "databases")
	f, err := os.Open(file)
	if err != nil {
		return merrors.InternalServerError(s.id, "could not open database directory")
	}
	defer f.Close()

	dnames, err := f.Readdirnames(0)
	if err != nil {
		return merrors.InternalServerError(s.id, "could not read database directory")
	}

	dbres.Databases = dnames
	return nil
}

// Tables implements the StoreHandler interface.
func (s *Service) Tables(ctx context.Context, in *proto.TablesRequest, out *proto.TablesResponse) error {
	file := filepath.Join(s.Config.Datapath, "databases", in.Database)
	f, err := os.Open(file)
	if err != nil {
		return merrors.InternalServerError(s.id, "could not open tables directory")
	}
	defer f.Close()

	tnames, err := f.Readdirnames(0)
	if err != nil {
		return merrors.InternalServerError(s.id, "could not read tables directory")
	}

	out.Tables = tnames
	return nil
}
