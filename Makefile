version := 0.0.1
author  := Ankit Jindal
app     := RabbitMQ Notify
id      := com.invitae.rabbitmq-notify


.PHONY:build

install:
	go get -u ./... && go mod tidy

build-darwin: clean
	GOOS=darwin GOARCH=amd64 go build -v -ldflags "-s -w -X main.version=$(version)" -o bin/darwin/rmq ./

bundle-darwin: build-darwin
	cd bin/darwin && \
		appify \
			-author "$(author)" \
			-id $(id) \
			-version $(version) \
			-name "$(app)" \
			-icon ../../assets/icon.png \
			./rmq
	/usr/libexec/PlistBuddy -c 'Add :LSUIElement bool true' 'bin/darwin/$(app).app/Contents/Info.plist'
	# codesign --force --deep --sign - 'bin/darwin/$(app).app'
	rm 'bin/darwin/$(app).app/Contents/README'
	rm bin/darwin/rmq

bundle-darwin-dmg: bundle-darwin
	cd bin/darwin/ && \
		create-dmg --dmg-title='$(app)' '$(app).app' ./ || true # ignore Error 2
	# Rename .dmg appropriately
	mv 'bin/darwin/$(app) $(version).dmg' bin/darwin/RabbitMQ_Notify_$(version)_darwin_x86_64.dmg
	rm -rf 'bin/darwin/$(app).app'

clean:
	rm -rf bin/

kill:
	pkill "$(app)" || true

start: run
run: bundle-darwin
	pkill "$(app)" || true
	open bin/darwin/"$(app)".app

tag:
	git tag -a v$(version) -m "Version $(version)"
	git push origin --tags

release:
	cd bin/darwin && \
		tar -cvzf '$(app).tar.gz' '$(app).app'
	gh release create v$(version) -n="$(version)" -t="$(version)" 'bin/darwin/$(app).tar.gz'
	# gh release create v$(version) 'bin/darwin/$(app).tar.gz'