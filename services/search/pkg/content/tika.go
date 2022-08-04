package content

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bbalet/stopwords"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/google/go-tika/tika"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
	"io"
	"strings"
)

// Tika is used to extract content from a resource,
// it uses apache tika to retrieve all the data.
type Tika struct {
	*Basic
	Retriever
	tika *tika.Client
}

// NewTikaExtractor creates a new Tika instance.
func NewTikaExtractor(gw gateway.GatewayAPIClient, logger log.Logger, cfg *config.Config) (*Tika, error) {
	basic, err := NewBasicExtractor(logger)
	if err != nil {
		return nil, err
	}

	tk := tika.NewClient(nil, cfg.Extractor.Tika.TikaURL)
	tkv, err := tk.Version(context.Background())
	if err != nil {
		return nil, err
	}
	logger.Info().Msg(fmt.Sprintf("Tika version: %s", tkv))

	return &Tika{
		Basic:     basic,
		Retriever: newCS3Retriever(gw, logger, cfg.MachineAuthAPIKey, cfg.Extractor.CS3AllowInsecure),
		tika:      tika.NewClient(nil, cfg.Extractor.Tika.TikaURL),
	}, nil
}

// Extract loads a resource from its underlying storage, passes it to tika and processes the result into a Document.
func (t Tika) Extract(ctx context.Context, ref *provider.Reference, ri *provider.ResourceInfo) (Document, error) {
	doc, err := t.Basic.Extract(ctx, ref, ri)
	if err != nil {
		return doc, err
	}

	data, err := t.Retrieve(ctx, ref, &user.User{Id: ri.Owner})
	if err != nil {
		return doc, err
	}
	defer data.Close()

	var d1, d2 bytes.Buffer
	if _, err := io.Copy(io.MultiWriter(&d1, &d2), data); err != nil {
		return doc, err
	}

	lang, err := t.tika.Language(ctx, bytes.NewReader(d1.Bytes()))
	if err != nil {
		return doc, nil
	}

	metas, err := t.tika.MetaRecursive(ctx, bytes.NewReader(d2.Bytes()))
	if err != nil {
		return doc, err
	}

	for _, meta := range metas {
		if title, err := getFirstValue(meta, "title"); err == nil {
			doc.Title = strings.TrimSpace(fmt.Sprintf("%s %s", doc.Title, title))
		}

		if content, err := getFirstValue(meta, "X-TIKA:content"); err == nil {
			if lang != "" {
				content = stopwords.CleanString(content, lang, true)
			}

			doc.Content = strings.TrimSpace(fmt.Sprintf("%s %s", doc.Content, content))
		}
	}

	return doc, nil
}
