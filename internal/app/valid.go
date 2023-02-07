package app

import (
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.CustomTypeTagMap.Set("orderID", func(i interface{}, context interface{}) bool {
		v, ok := i.(string)
		if !ok {
			return false
		}
		if v == "" {
			return false
		}
		return goluhn.Validate(v) == nil
	})

}
