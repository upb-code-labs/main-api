WORDIR = $(shell pwd)

create_package: 
	@read -p "🏷️ Package name: " name; \
	cd $(WORDIR); \
	mkdir src/$$name; \
	mkdir src/$$name/domain src/$$name/application src/$$name/infrastructure; 

delete_pakage: 
	@read -p "🏷️ Package name: " name; \
	cd $(WORDIR); \
	rm -rf src/$$name; 
	
coverage:
	go clean -testcache; \
	go test -coverpkg ./... -coverprofile coverage.txt __tests__/integration/*.go; \
	go tool cover -html=coverage.txt -o coverage.html; 