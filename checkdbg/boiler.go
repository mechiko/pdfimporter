package checkdbg

import (
	"fmt"

	"github.com/mechiko/utility"
)

func (c *Checks) TestUtilityParseCis() error {
	znak := []string{
		`0104630277410026215!&CS1r93VoB/`,
		`0104630277410361215keTalmIz+ZiL\u001d935njT`,
	}
	for _, code := range znak {
		if cis, err := utility.ParseCisInfo(code); err != nil {
			c.Logger().Errorf("parse znak %v", err)
			return fmt.Errorf("%w", err)
		} else {
			c.Logger().Debugf("тест км результат: серийник %s gtin %s", cis.Serial, cis.Gtin)
		}
	}
	return nil
}
