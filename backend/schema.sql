CREATE TABLE "car" (
	"id"	INTEGER,
	"make"	TEXT,
	"model"	TEXT,
	"engine"	TEXT,
	"year"	TEXT,
	"odometer"	INTEGER,
	"vin"	TEXT,
	"mot_due"	INTEGER,
	"tax_due"	INTEGER,
	"insured_until"	INTEGER,
	"owner_id"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("owner_id") REFERENCES "users"("id") ON DELETE CASCADE
)

CREATE TABLE "maintenance_records" (
  "id"  INTEGER,
  "car_id" INTEGER,
  "type_id" INTEGER,
  "date_done" INTEGER,
  "mileage" INTEGER,
  "notes" TEXT,
  PRIMARY KEY("id" AUTOINCREMENT),
  FOREIGN KEY("car_id") REFERENCES "car"("id") ON DELETE CASCADE,
  FOREIGN KEY("type_id") REFERENCES "maintenance_types"("id") ON DELETE CASCADE
)


CREATE TABLE "maintenance_types" (
  "id" INTEGER,
  "name" TEXT UNIQUE,
  "interval_miles" INTEGER,
  "interval_months" INTEGER,
  "notes" TEXT,
  PRIMARY KEY("id" AUTOINCREMENT)
)

CREATE TABLE "obd2_readings" (
	"id"	INTEGER,
	"car_id"	INTEGER,
	"timestamp"	INTEGER,
	"odometer"	INTEGER,
	"fuel_level"	REAL,
	"coolant_temp"	REAL,
	"oil_temp"	REAL,
	"engine_rpm"	REAL,
	"horsepower"	REAL,
	"mass_air_flow"	REAL,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("car_id") REFERENCES "car"("id") ON DELETE CASCADE
)

CREATE TABLE "tyre_details" (
	"id"	INTEGER,
	"maintenance_record_id"	INTEGER NOT NULL,
	"position"	TEXT NOT NULL,
	"brand"	TEXT,
	"size"	TEXT,
	"notes"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("maintenance_record_id") REFERENCES "maintenance_records"("id") ON DELETE CASCADE
)


CREATE TABLE "users" (
	"id"	INTEGER,
	"username"	TEXT NOT NULL,
	"email"	TEXT NOT NULL UNIQUE,
	"password_hash"	TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
)

CREATE TABLE reminder (

id INTEGER PRIMARY KEY AUTOINCREMENT,
user_id INTEGER NOT NULL,
car_id INTEGER NOT NULL,
reminder_type TEXT NOT NULL,
title TEXT NOT NULL,
message TEXT NOT NULL,
due_date INTEGER,
maintenance_record_id INTEGER,
maintenance_type_id INTEGER,
threshold_miles INTEGER,
resolved INTEGER DEFAULT 0,
created_at INTEGER NOT NULL,
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
FOREIGN KEY (car_id) REFERENCES car(id) ON DELETE CASCADE,
FOREIGN KEY (maintenance_record_id) REFERENCES maintenance_records(id) ON DELETE CASCADE,
FOREIGN KEY (maintenance_type_id) REFERENCES maintenance_types(id) ON DELETE CASCADE
)
