IMAGE_NAME := memere-admin
IMAGE_TAG  := local
PORT       := 3000

.PHONY: build start docker-build docker-run

## Local dev / CI ----------------------------------------------------------------

build:
	pnpm build

start:
	pnpm start

typecheck:
	pnpm typecheck

lint:
	pnpm lint

## Docker -----------------------------------------------------------------------

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

# Run the container using .env.local for secrets.
# Requires: API_BASE_URL and COOKIE_SECRET to be set in .env.local
docker-run:
	@test -f .env.local || (echo "ERROR: .env.local not found — copy .env.local.example and fill in values" && exit 1)
	docker run --rm \
	  -p $(PORT):3000 \
	  --env-file .env.local \
	  -e NODE_ENV=production \
	  $(IMAGE_NAME):$(IMAGE_TAG)

# Print required runtime env vars
env-check:
	@echo "Required runtime env vars:"
	@echo "  API_BASE_URL    — Go backend base URL (e.g. https://api.memere.app)"
	@echo "  COOKIE_SECRET   — ≥32 random bytes for httpOnly cookie signing"
	@echo "  NODE_ENV        — set to 'production'"
	@echo "  NEXT_PUBLIC_APP_NAME — (optional) display name, defaults to 'Memere Admin'"
