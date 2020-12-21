PNGCRUSH := $(shell command -v pngcrush 2> /dev/null)
UNAME := $(shell uname -s)

# default

default: scale trim spritesheet scss copy

# main targets

.PHONY: rename
rename:
	go run cmd/rename/main.go

.PHONY: scale
scale:
	go run cmd/scale/main.go

.PHONY: trim
trim:
	go run cmd/trim/main.go

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

.PHONY: copy
copy:
	cp output/pokesprite.png ../pokedextracker.com/public
	cp output/pokesprite.scss ../pokedextracker.com/app/styles

# helper targets

.PHONY: clean
clean:
	rm -rf output

output:
	mkdir output

.PHONY: setup
setup:
ifeq ($(UNAME), Linux)
	sudo apt-get install pngcrush
endif
ifeq ($(UNAME), Darwin)
	brew install pngcrush
endif
