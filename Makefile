generate:
	sh ./scripts/buf/gen.sh
	sh ./scripts/ent/gen.sh
	sh ./scripts/embed/gen.sh

clean:
	rm -rf generated
	rm -rf web/gen
	rm -f web/tsconfig.app.tsbuildinfo
	rm -f web/tsconfig.node.tsbuildinfo

.PHONY: generate clean
