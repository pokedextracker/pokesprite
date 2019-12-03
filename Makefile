PNGCRUSH := $(shell command -v pngcrush 2> /dev/null)

# main targets

.PHONY: rename
rename:
	go run cmd/rename/main.go

.PHONY: spritesheet
spritesheet: output
	go run cmd/spritesheet/main.go
ifdef PNGCRUSH
	pngcrush -l 9 output/pokesprite.png output/pokesprite-optimized.png
	mv output/pokesprite-optimized.png output/pokesprite.png
else
	@echo "---> skipping png optimization because pngcrush is not installed"
endif

.PHONY: scss
scss: output
	go run cmd/scss/main.go

# helper targets

.PHONY: clean
clean:
	rm -rf output

output:
	mkdir output

.PHONY: setup
setup:
	brew install pngcrush
