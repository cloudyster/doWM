PREFIX ?= /usr/local
BINDIR := $(PREFIX)/bin
XSESSIONS := /usr/share/xsessions
USER_CONFIG := $(HOME)/.config/doWM

build:
	go build -o doWM
	@echo "Built successfully!"

install:
	# Install binary locally
	mkdir -p $(BINDIR)
	sudo install -m755 doWM $(BINDIR)/doWM

	# Install .desktop session file
	mkdir -p $(XSESSIONS)
	sudo install -m644 doWM.desktop $(XSESSIONS)/doWM.desktop


	@echo "Installed successfully!"

config:
	@echo "Setting up doWM user config..."
	mkdir -p $(USER_CONFIG)
	@if [ ! -f $(USER_CONFIG)/autostart.sh ]; then \
		cp exampleConfig/autostart.sh $(USER_CONFIG)/autostart.sh && \
		chmod +x $(USER_CONFIG)/autostart.sh && \
		echo "Installed example autostart.sh"; \
	else \
		echo "autostart.sh already exists, skipping..."; \
	fi
	@if [ ! -f $(USER_CONFIG)/doWM.yml ]; then \
		cp exampleConfig/doWM.yml $(USER_CONFIG)/doWM.yml && \
		echo "Installed example doWM.yml"; \
	else \
		echo "doWM.yml already exists, skipping..."; \
	fi
	@echo "Config setup complete!"

.PHONY: all install uninstall
