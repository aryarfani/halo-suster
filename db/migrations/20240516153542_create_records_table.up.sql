CREATE TABLE IF NOT EXISTS "records" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "patient_identity_number" BIGINT NOT NULL,
    "symptoms" TEXT NOT NULL,
    "medications" TEXT NOT NULL,
    "user_id" UUID NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "records" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "records" ADD FOREIGN KEY ("patient_identity_number") REFERENCES "patients" ("identity_number");