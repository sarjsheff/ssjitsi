
all: docs ssjitsi

docs:
	swag init -g cmd/server/main.go -o docs

ssjitsi: build-ui
	go build -o ssjitsi cmd/server/main.go

build-ui:
	cd web/app && npm run build
	mkdir -p internal/pkg/ssjitsi/web
	cp -r web/app/dist/* internal/pkg/ssjitsi/web/

clean:
	rm -rf ssjitsi
	rm -rf docs
	rm -rf internal/pkg/ssjitsi/web