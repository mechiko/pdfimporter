package pdfproc

import (
	"fmt"
	"slices"
	"strings"

	"github.com/mechiko/maroto/v2/pkg/components/code"
	"github.com/mechiko/maroto/v2/pkg/components/col"
	"github.com/mechiko/maroto/v2/pkg/components/image"
	"github.com/mechiko/maroto/v2/pkg/components/page"
	"github.com/mechiko/maroto/v2/pkg/components/row"
	"github.com/mechiko/maroto/v2/pkg/components/text"
	"github.com/mechiko/maroto/v2/pkg/core"
)

func (p *pdfProc) Page(t *MarkTemplate, kod string, ser string, idx string) (core.Page, error) {
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
						code.NewMatrix(kod, rowSingle.PropsRect()),
					)
					if rowSingle.ImageDebug {
						cols[i].WithStyle(colStyle)
					}
				}
				if rowSingle.Bar != "" {
					cols[i].Add(
						code.NewBar(kod, rowSingle.PropsBar()),
					)
					if rowSingle.ImageDebug {
						cols[i].WithStyle(colStyle)
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
							cols[i].WithStyle(colStyle)
						}
					} else {
						return nil, fmt.Errorf("page image assets not available (assets is nil) for %q", rowSingle.Image)
					}
				}
				if rowSingle.Value != "" {
					value := strings.ReplaceAll(rowSingle.Value, "@ser", ser)
					value = strings.ReplaceAll(value, "@idx", idx)
					if len(kod) == 20 {
						kod1 := fmt.Sprintf("(%s)%s", kod[:2], kod[2:])
						value = strings.ReplaceAll(value, "@kod", kod1)
					} else {
						value = strings.ReplaceAll(value, "@kod", kod)
					}
					cols[i].Add(text.New(value, rowSingle.PropsText()))
				} else {
					if len(rowSingle.Values) > 0 {
						comps := make([]core.Component, 0)
						for _, val := range rowSingle.Values {
							value := strings.ReplaceAll(val.Value, "@ser", ser)
							value = strings.ReplaceAll(value, "@idx", idx)
							if len(kod) == 20 {
								kod1 := fmt.Sprintf("(%s)%s", kod[:2], kod[2:])
								value = strings.ReplaceAll(value, "@kod", kod1)
							} else {
								value = strings.ReplaceAll(value, "@kod", kod)
							}
							comps = append(comps, text.New(value, val.PropsText()))
						}
						cols[i].Add(comps...)
						// cols[i] = text.NewCol(rowSingle.ColWidth, value, rowSingle.PropsText())
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
