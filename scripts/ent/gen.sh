mkdir ./generated/ent -p

if [ ! -f "./generated/ent/ent.go" ]; then
    touch "./generated/ent/ent.go"
    echo "package ent" > "./generated/ent/ent.go"
fi

go run -mod=mod entgo.io/ent/cmd/ent generate ./pkgs/ent/schema --target ./generated/ent
