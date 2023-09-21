WORDIR = $(shell pwd)

create_package: 
	@read -p "🏷️ Package name: " name; \
	cd $(WORDIR); \
	mkdir src/$$name; \
	mkdir src/$$name/domain src/$$name/application src/$$name/infrastructure;  \
	mkdir tests/$$name;

delete_pakage: 
	@read -p "🏷️ Package name: " name; \
	cd $(WORDIR); \
	rm -rf src/$$name; \
	rm -rf tests/$$name;