generate:
	sh ./scripts/buf/gen.sh
	sh ./scripts/ent/gen.sh
	sh ./scripts/embed/gen.sh

.PHONY: generate
