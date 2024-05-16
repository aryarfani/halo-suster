CREATE TABLE IF NOT EXISTS "patients" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "identity_number" BIGINT UNIQUE NOT NULL,
    "phone_number" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "birth_date" DATE NOT NULL,
    "gender" TEXT NOT NULL,
    "identity_card_scan_img" TEXT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);