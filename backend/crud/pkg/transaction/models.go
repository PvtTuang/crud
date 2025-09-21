package transaction

import "gorm.io/gorm"

type Transaction struct {
	tx *gorm.DB
}
