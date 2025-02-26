
# Variables
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=janan
LOCAL_DATABASE_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable


## hash: Hash (re-)creates an integrity hash file for the migration directory.
.PHONY: hash
hash:
	cd cmd/database/ && atlas migrate hash --env local

##diff: Calculate and print the diff between two schemas.
##usage migrate diff msg=yourmessage
.PHONY: diff
diff:
ifndef msg
	@echo "Error: Please provide a value for the 'msg' parameter."
	@exit 1
endif
	cd cmd/database/ && atlas migrate diff --env local ${msg}

##migrate-up : execute the SQL statements in the migration files
.PHONY: migrate-up
migrate-up:
	cd cmd/database/ && migrate -path ./migrations/ -database ${LOCAL_DATABASE_URL} -verbose up
##migrate-down : to revert back the migration
.PHONY: migrate-down
migrate-down:
	cd cmd/database/ && migrate -path ./migrations/ -database ${LOCAL_DATABASE_URL} -verbose down

##schema : shows the dn schema url
.PHONY: schema
schema:
	cd cmd/database/ && atlas schema inspect --web --env local

##force: force migrate the database to a specific version
.PHONY: force
force:
ifndef version
	@echo "Error: Please provide a value for the 'version' parameter."
	@exit 1
endif
	cd cmd/database/ && migrate -path migrations -database ${LOCAL_DATABASE_URL} force ${version}
# ==================================================================================== #
# Swagger
# ==================================================================================== #
swaggo:
	echo "Starting swagger level 1 generating"
	swag init   -g ./**/*.go 