# Tools

CONVERT   ?= convert
IDENTIFY  ?= identify
BASE64    ?= base64
ENVSUBST  ?= envsubst
SCOUR     ?= scour
INKSCAPE  ?= inkscape

# Variables

STATIC ?= ../../public/static
ICON   ?= lico-icon.svg

# Build

all: app-icon.svg

.PHONY: icons
icons: $(STATIC)/favicon.ico

.PHONY: $(STATIC)/favicon.ico
$(STATIC)/favicon.ico: app-icon-rounded-256x256.png
	$(CONVERT) -background transparent $< -define icon:auto-resize=16,32,48,64,128,256 $@

app-icon.svg: $(ICON)
	cp -vaf $< $@

app-icon-whitebox-256x256.png: app-icon.svg
	$(INKSCAPE) -z -e $@.tmp -w 204.8 -h 204.8 -b white -y 1.0 $<
	$(CONVERT) $@.tmp -background white -gravity center -extent 256x256 $@
	@$(RM) $@.tmp

app-icon-rounded-256x256.png: app-icon-whitebox-256x256.png
	$(CONVERT) -size 256x256 xc:none -draw "roundrectangle 2,2,252,252,126,126" $@.tmp.png
	$(CONVERT) $< -matte $@.tmp.png -compose DstIn -composite $@
	@$(RM) $@.tmp.png

.PHONY: clean
clean:
	$(RM) app-icon-*.png || true
	$(RM) app-icon.svg || true
