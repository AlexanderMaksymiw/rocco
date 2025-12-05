CREATE TABLE "car" (
	"id"	INTEGER,
	"make"	TEXT,
	"model"	TEXT,
	"engine"	TEXT,
	"year"	TEXT,
	"odometer"	INTEGER,
	"vin"	TEXT,
	"last_service_miles"	INTEGER,
	"last_service_date"	INTEGER,
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
