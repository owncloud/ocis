package ocis.authz

default allow = true

allow := false {
  input.method == "PUT"
  endswith(input.path, ".txt")
}
