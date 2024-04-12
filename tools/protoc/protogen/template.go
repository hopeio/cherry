package main

const ymlTpl = `schema:
  - ./*.graphql

# Where should the generated server code go?
exec:
  filename: ../../{{.SubDir}}/{{.SubDir}}.generated.gql.go
  package: {{.SubDir}}

# Enable Apollo federation support
federation:
  filename: ../../{{.SubDir}}/{{.SubDir}}.federation.gql.go
  package: {{.SubDir}}

model:
  filename: ../../{{.SubDir}}/{{.SubDir}}.model.gql.go
  package: {{.SubDir}}

autobind:
{{range .Packages}}  - {{$.GoMod}}/{{.}}
{{end}}  - github.com/hopeio/cherry/protobuf/request
  - github.com/hopeio/cherry/protobuf/response
  - github.com/hopeio/cherry/protobuf/oauth

struct_fields_always_pointers: false

models:
  ID:
    model:
      - github.com/99designs/gqlgen/graphql.ID
  Id:
    model:
      - github.com/99designs/gqlgen/graphql.ID
  UUID:
    model:
      - github.com/99designs/gqlgen/graphql.UUID
  Int:
    model:
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
      - github.com/hopeio/cherry/utils/net/http/api/graphql.Uint
      - github.com/hopeio/cherry/utils/net/http/api/graphql.Uint64
      - github.com/hopeio/cherry/utils/net/http/api/graphql.Uint32
  Int32:
    model: github.com/99designs/gqlgen/graphql.Int32
  Int64:
    model: github.com/99designs/gqlgen/graphql.Int64
  Uint8:
    model: github.com/hopeio/cherry/utils/net/http/api/graphql.Uint8
  Uint:
    model:
      - github.com/hopeio/cherry/utils/net/http/api/graphql.Uint
      - github.com/hopeio/cherry/utils/net/http/api/graphql.Uint64
      - github.com/hopeio/cherry/utils/net/http/api/graphql.Uint32
  Uint32:
      model: github.com/hopeio/cherry/utils/net/http/api/graphql.Uint32
  Uint64:
      model: github.com/hopeio/cherry/utils/net/http/api/graphql.Uint64
  Float32:
    model: github.com/hopeio/cherry/utils/net/http/api/graphql.Float32
  Float64:
    model: github.com/hopeio/cherry/utils/net/http/api/graphql.Float64
  Float:
    model: github.com/99designs/gqlgen/graphql.Float
  Bytes:
    model: github.com/hopeio/cherry/utils/net/http/api/graphql.Bytes
  Header:
    model: github.com/hopeio/cherry/utils/net/http/api/graphql.Header
`

//经过一番查找，发现yaml语法对格式是非常严格的，不可以有制表符！不可以有制表符！不可以有制表符！
//缩进也有要求
