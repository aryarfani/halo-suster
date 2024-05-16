CREATE TABLE IF NOT EXISTS "users" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "nip" BIGINT UNIQUE NOT NULL,
    "name" TEXT NOT NULL,
    "role" TEXT NOT NULL,
    "password" TEXT NULL,
    "identity_card_scan_img" TEXT NULL,
    "deleted_at" TIMESTAMP WITH TIME ZONE NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);