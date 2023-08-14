package goheif

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"

	"github.com/mholt/goheif/heif"
	"github.com/mholt/goheif/libde265"
)

// SafeEncoding uses more memory but seems to make
// the library safer to use in containers.
var SafeEncoding bool

type gridBox struct {
	columns, rows int
	width, height int
}

func newGridBox(data []byte) (*gridBox, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("invalid data")
	}
	// version := data[0]
	flags := data[1]
	rows := int(data[2]) + 1
	columns := int(data[3]) + 1

	var width, height int
	if (flags & 1) != 0 {
		if len(data) < 12 {
			return nil, fmt.Errorf("invalid data")
		}

		width = int(data[4])<<24 | int(data[5])<<16 | int(data[6])<<8 | int(data[7])
		height = int(data[8])<<24 | int(data[9])<<16 | int(data[10])<<8 | int(data[11])
	} else {
		width = int(data[4])<<8 | int(data[5])
		height = int(data[6])<<8 | int(data[7])
	}

	return &gridBox{columns: columns, rows: rows, width: width, height: height}, nil
}

func decodeHevcItem(dec *libde265.Decoder, hf *heif.File, item *heif.Item) (*image.YCbCr, error) {
	if item.Info.ItemType != "hvc1" {
		return nil, fmt.Errorf("Unsupported item type: %s", item.Info.ItemType)
	}

	hvcc, ok := item.HevcConfig()
	if !ok {
		return nil, fmt.Errorf("No hvcC")
	}

	hdr := hvcc.AsHeader()
	data, err := hf.GetItemData(item)
	if err != nil {
		return nil, err
	}

	dec.Reset()
	dec.Push(hdr)
	tile, err := dec.DecodeImage(data)
	if err != nil {
		return nil, err
	}

	ycc, ok := tile.(*image.YCbCr)
	if !ok {
		return nil, fmt.Errorf("Tile is not YCbCr")
	}

	return ycc, nil
}

func ExtractExif(ra io.ReaderAt) ([]byte, error) {
	hf := heif.Open(ra)
	return hf.EXIF()
}

func Decode(r io.Reader) (image.Image, error) {
	ra, err := asReaderAt(r)
	if err != nil {
		return nil, err
	}

	hf := heif.Open(ra)

	it, err := hf.PrimaryItem()
	if err != nil {
		return nil, err
	}

	width, height, ok := it.SpatialExtents()
	if !ok {
		return nil, fmt.Errorf("No dimension")
	}

	if it.Info == nil {
		return nil, fmt.Errorf("No item info")
	}

	dec, err := libde265.NewDecoder(libde265.WithSafeEncoding(SafeEncoding))
	if err != nil {
		return nil, err
	}
	defer dec.Free()
	if it.Info.ItemType == "hvc1" {
		return decodeHevcItem(dec, hf, it)
	}

	if it.Info.ItemType != "grid" {
		return nil, fmt.Errorf("No grid")
	}

	data, err := hf.GetItemData(it)
	if err != nil {
		return nil, err
	}

	grid, err := newGridBox(data)
	if err != nil {
		return nil, err
	}

	dimg := it.Reference("dimg")
	if dimg == nil {
		return nil, fmt.Errorf("No dimg")
	}

	if len(dimg.ToItemIDs) != grid.columns*grid.rows {
		return nil, fmt.Errorf("Tiles number not matched")
	}

	var out *image.YCbCr
	var tileWidth, tileHeight int
	for i, y := 0, 0; y < grid.rows; y += 1 {
		for x := 0; x < grid.columns; x += 1 {
			id := dimg.ToItemIDs[i]
			item, err := hf.ItemByID(id)
			if err != nil {
				return nil, err
			}

			ycc, err := decodeHevcItem(dec, hf, item)
			if err != nil {
				return nil, err
			}

			rect := ycc.Bounds()
			if tileWidth == 0 {
				tileWidth, tileHeight = rect.Dx(), rect.Dy()
				width, height := tileWidth*grid.columns, tileHeight*grid.rows
				out = image.NewYCbCr(image.Rectangle{image.Pt(0, 0), image.Pt(width, height)}, ycc.SubsampleRatio)
			}

			if tileWidth != rect.Dx() || tileHeight != rect.Dy() {
				return nil, fmt.Errorf("Inconsistent tile dimensions")
			}

			// copy y stride data
			for i := 0; i < rect.Dy(); i += 1 {
				copy(out.Y[(y*tileHeight+i)*out.YStride+x*ycc.YStride:], ycc.Y[i*ycc.YStride:(i+1)*ycc.YStride])
			}

			// height of c strides
			cHeight := len(ycc.Cb) / ycc.CStride

			// copy c stride data
			for i := 0; i < cHeight; i += 1 {
				copy(out.Cb[(y*cHeight+i)*out.CStride+x*ycc.CStride:], ycc.Cb[i*ycc.CStride:(i+1)*ycc.CStride])
				copy(out.Cr[(y*cHeight+i)*out.CStride+x*ycc.CStride:], ycc.Cr[i*ycc.CStride:(i+1)*ycc.CStride])
			}

			i += 1
		}
	}

	//crop to actual size when applicable
	out.Rect = image.Rectangle{image.Pt(0, 0), image.Pt(width, height)}
	return out, nil
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	var config image.Config

	ra, err := asReaderAt(r)
	if err != nil {
		return config, err
	}

	hf := heif.Open(ra)

	it, err := hf.PrimaryItem()
	if err != nil {
		return config, err
	}

	width, height, ok := it.SpatialExtents()
	if !ok {
		return config, fmt.Errorf("No dimension")
	}

	config = image.Config{
		ColorModel: color.YCbCrModel,
		Width:      width,
		Height:     height,
	}
	return config, nil
}

func asReaderAt(r io.Reader) (io.ReaderAt, error) {
	if ra, ok := r.(io.ReaderAt); ok {
		return ra, nil
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func init() {
	libde265.Init()
	// they check for "ftyp" at the 5th bytes, let's do the same...
	// https://github.com/strukturag/libheif/blob/master/libheif/heif.cc#L94
	image.RegisterFormat("heic", "????ftyp", Decode, DecodeConfig)
}
