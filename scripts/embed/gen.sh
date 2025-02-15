sh scripts/embed/setup.sh

cd web && pnpm run build --outDir "../generated/embed/frontend/" --emptyOutDir
