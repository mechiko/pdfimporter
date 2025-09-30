package pdfproc

import (
	"bytes"
	"fmt"
	imgstd "image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"pdfimporter/domain"
	"slices"
	"strings"

	svgbar "github.com/juliankoehn/barcode"
	"github.com/mechiko/barcode"
	"github.com/mechiko/barcode/ean"
	"github.com/mechiko/maroto/v2/pkg/components/code"
	"github.com/mechiko/maroto/v2/pkg/components/col"
	"github.com/mechiko/maroto/v2/pkg/components/image"
	"github.com/mechiko/maroto/v2/pkg/components/page"
	"github.com/mechiko/maroto/v2/pkg/components/row"
	"github.com/mechiko/maroto/v2/pkg/components/text"
	"github.com/mechiko/maroto/v2/pkg/consts/extension"
	"github.com/mechiko/maroto/v2/pkg/core"
	"github.com/mechiko/utility"
)

func (p *pdfProc) Page(t *domain.MarkTemplate, cis *utility.CisInfo, party string, idx string) (core.Page, error) {
	pg := page.New()
	rowKeys := make([]string, 0, len(t.Rows))
	for k := range t.Rows {
		rowKeys = append(rowKeys, k)
	}
	slices.Sort(rowKeys)
	for _, rowKey := range rowKeys {
		rowTempl := t.Rows[rowKey]
		switch {
		case len(rowTempl) == 0:
		case len(rowTempl) == 1:
			// одна строка автороу
			row1 := rowTempl[0]
			if row1.Value == "" {
				// пустая строка с высотой
				pg.Add(
					row.New(row1.RowHeight).Add(),
				)
			} else {
				if row1.RowHeight == 0 {
					pg.Add(
						text.NewAutoRow(row1.Value, row1.PropsText()),
					)
				} else {
					pg.Add(
						row.New(row1.RowHeight).Add(
							text.NewCol(12, row1.Value, row1.PropsText()),
						),
					)
				}
			}
		case len(rowTempl) > 1:
			cols := make([]core.Col, len(rowTempl))
			// строки с колонками
			for i, rowSingle := range rowTempl {
				cols[i] = col.New(rowSingle.ColWidth)
				if rowSingle.DataMatrix != "" {
					cols[i].Add(
						code.NewMatrix(cis.FNC1(), rowSingle.PropsRect()),
					)
					if rowSingle.ImageDebug {
						cols[i].WithStyle(domain.ColStyle)
					}
				}
				if rowSingle.Bar != "" {
					switch rowSingle.Bar {
					case "ean13":
						ean13 := strings.Trim(cis.Gtin, "0")
						cols[i].Add(
							code.NewBar(ean13, rowSingle.PropsBar()),
						)
					}
					if rowSingle.ImageDebug {
						cols[i].WithStyle(domain.ColStyle)
					}
				}
				if rowSingle.Image != "" {
					if p.assets != nil {
						img, err := p.assets.Jpg(rowSingle.Image)
						if err != nil {
							return nil, fmt.Errorf("page image assets error %w", err)
						}
						if len(img) == 0 {
							return nil, fmt.Errorf("page image assets empty for %q", rowSingle.Image)
						}
						cols[i].Add(
							image.NewFromBytes(img, rowSingle.ImageExt, rowSingle.PropsRect()),
						)
						if rowSingle.ImageDebug {
							cols[i].WithStyle(domain.ColStyle)
						}
					} else {
						return nil, fmt.Errorf("page image assets not available (assets is nil) for %q", rowSingle.Image)
					}
				}
				if rowSingle.Value != "" {
					value := strings.ReplaceAll(rowSingle.Value, "@party", party)
					value = strings.ReplaceAll(value, "@idx", idx)
					ean13 := strings.Trim(cis.Gtin, "0")
					ean13 = fmt.Sprintf("%s  %s  %s", ean13[:1], ean13[1:7], ean13[7:])
					value = strings.ReplaceAll(value, "@ean", ean13)
					cols[i].Add(text.New(value, rowSingle.PropsText()))
				} else {
					if len(rowSingle.Values) > 0 {
						comps := make([]core.Component, 0)
						for _, val := range rowSingle.Values {
							value := ""
							if val.Value != "" {
								value = strings.ReplaceAll(val.Value, "@party", party)
								value = strings.ReplaceAll(value, "@idx", idx)
								ean13 := strings.Trim(cis.Gtin, "0")
								ean13 = fmt.Sprintf("%s  %s  %s", ean13[:1], ean13[1:7], ean13[7:])
								value = strings.ReplaceAll(value, "@ean", ean13)
								comps = append(comps, text.New(value, val.PropsText()))
							}
							if val.DataMatrix != "" {
								comps = append(comps, code.NewMatrix(cis.FNC1(), val.PropsRect()))
							}
							if val.Bar != "" {
								switch val.Bar {
								case "ean13":
									ean13 := strings.Trim(cis.Gtin, "0")
									comps = append(comps, code.NewBar(ean13, rowSingle.PropsBar()))
								case "ean13b":
									ean13 := strings.Trim(cis.Gtin, "0")
									img, err := barImg(ean13)
									if err != nil {
										return pg, fmt.Errorf("ean13 bar error %w", err)
									}
									comps = append(comps, image.NewFromBytes(img, extension.Jpg, val.PropsRect()))
								case "ean13svg":
									ean13 := strings.Trim(cis.Gtin, "0")
									img, err := svgImg(ean13)
									if err != nil {
										return pg, fmt.Errorf("ean13 svgImg error %w", err)
									}
									comps = append(comps, image.NewFromBytes(img, extension.Jpg, val.PropsRect()))
								case "ean13j":
									img, err := p.assets.Jpg("gtin")
									if err != nil {
										return nil, fmt.Errorf("page image assets error %w", err)
									}
									if len(img) == 0 {
										return nil, fmt.Errorf("page image assets empty for %q", rowSingle.Image)
									}
									comps = append(comps, image.NewFromBytes(img, extension.Jpg, val.PropsRect()))
								}
							}
						}
						cols[i].Add(comps...)
					}
				}
			}
			pg.Add(
				row.New(rowTempl[0].RowHeight).Add(cols...),
			)
		default:
		}
	}
	return pg, nil
}

func barImg(code string) ([]byte, error) {
	bcImg, err := ean.Encode(code)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	scaled, err := barcode.Scale(bcImg, 100, 25)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	var b bytes.Buffer
	jpeg.Encode(&b, scaled, nil)
	png.Encode(&b, scaled)
	return b.Bytes(), nil
}

func svgImg(code string) ([]byte, error) {
	var b bytes.Buffer
	width := 2
	height := 35
	color := "black"
	showCode := false
	which := "EAN13"
	extension := ".jpg"
	if extension == ".png" || extension == ".jpg" || extension == ".jpeg" || extension == ".gif" {
		var img *imgstd.RGBA
		if extension == ".png" {
			_, img = svgbar.GetBarcodeFile(code, which, width, height, color, showCode, false, true)
		} else {
			_, img = svgbar.GetBarcodeFile(code, which, width, height, color, showCode, false, false)
		}
		switch extension {
		case ".png":
			png.Encode(&b, img)
		case ".jpg", ".jpeg":
			jpeg.Encode(&b, img, &jpeg.Options{Quality: 100})
		case ".gif":
			gif.Encode(&b, img, nil)
		}
	}
	return b.Bytes(), nil
}
