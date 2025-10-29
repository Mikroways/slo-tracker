env "local" {
  src = "file://schema.hcl"
  dev = "docker://postgres/15/slotracker"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "default" { # expected to be used for the helm pre-install hook
  url = getenv("DATABASE_URL")
  migration {
    dir = "file://migrations"
  }
}