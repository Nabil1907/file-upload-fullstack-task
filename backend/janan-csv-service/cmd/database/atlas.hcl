data "external_schema" "gorm" {
  program = ["env", "ENCORERUNTIME_NOPANIC=1", "go", "run", "./main.go"]
}

env "local" {
  src = data.external_schema.gorm.url
  dev = "postgres://postgres:postgres@localhost:5432/janan?sslmode=disable"
  url = "postgres://postgres:postgres@localhost:5432/janan?sslmode=disable"
  diff {
    skip {
      drop_schema = true
    }
  }
  migration {
    dir = "file://migrations"
    format = golang-migrate
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

