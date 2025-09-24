package config

import (
	"database/sql"
	"log"
)

func Migrate(db *sql.DB) {
	ModelQuery := []string{
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`,

		`CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);`,

		`CREATE TABLE IF NOT EXISTS products (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT NOT NULL,
		price DOUBLE PRECISION NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		deleted_at TIMESTAMPTZ 
	);`,

		`CREATE TABLE IF NOT EXISTS carts (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);`,

		`CREATE TABLE IF NOT EXISTS cart_products (
		id BIGSERIAL PRIMARY KEY ,
		cart_id UUID NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
		product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
		quantity INTEGER NOT NULL
	);`,

		`CREATE TABLE IF NOT EXISTS purchase_histories (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);`,

		`CREATE TABLE IF NOT EXISTS purchase_items (
		id BIGSERIAL PRIMARY KEY ,
		purchase_id UUID NOT NULL REFERENCES purchase_histories(id) ON DELETE CASCADE,
		product_id UUID NOT NULL REFERENCES products(id),
		quantity INTEGER NOT NULL
	);`,
	}

	for _, q := range ModelQuery {
		if _, err := db.Exec(q); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	}
	log.Println("Migration เสร็จสิ้น")

}
