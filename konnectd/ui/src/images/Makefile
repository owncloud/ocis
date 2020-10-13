# Tools

CONVERT   ?= convert
IDENTIFY  ?= identify
BASE64    ?= base64
ENVSUBST  ?= envsubst
SCOUR     ?= scour
INKSCAPE  ?= inkscape

# Variables

LOGINSCREEN_BG_SRC ?= kopano-augmenting-teamwork-bg-7680.jpg
LOGINSCREEN_BG_OVERLAY_SRC ?= kopano-augmenting-teamwork-overlay.svg

LOGINSCREEN_BG_WIDTH            = $(shell $(IDENTIFY) -format '%w' loginscreen-bg.jpg)
LOGINSCREEN_BG_HEIGHT           = $(shell $(IDENTIFY) -format '%h' loginscreen-bg.jpg)
LOGINSCREEN_BG_THUMB_BASE64     = $(shell $(BASE64) -w0 loginscreen-bg-thumb.jpg)
LOGINSCREEN_BG_THUMB_BASE64_SVG = $(shell $(BASE64) -w0 loginscreen-bg-thumb.svg)

STATIC ?= ../../public/static
ICON   ?= kopano-konnect-icon.svg

# Build

all: loginscreen-bg.css

loginscreen-bg.jpg: $(LOGINSCREEN_BG_SRC)
	$(CONVERT) -geometry x1080 -strip -interlace Plane -gaussian-blur 0.05 -define jpeg:dct-method=float -quality 75% $< $@

loginscreen-bg-thumb.jpg: loginscreen-bg.jpg
	$(CONVERT) -geometry x40 -strip -define jpeg:dct-method=float -quality 50% $< $@

loginscreen-bg-thumb.svg: loginscreen-bg-thumb.svg.in | loginscreen-bg-thumb.jpg loginscreen-bg.jpg
	WIDTH=$(LOGINSCREEN_BG_WIDTH) HEIGHT=$(LOGINSCREEN_BG_HEIGHT) \
		IMAGE_DATA=$(LOGINSCREEN_BG_THUMB_BASE64) \
		$(ENVSUBST) < $< > $@

loginscreen-bg-overlay.svg: $(LOGINSCREEN_BG_OVERLAY_SRC)
	$(SCOUR) --enable-viewboxing --create-groups --shorten-ids --enable-id-stripping \
		--enable-comment-stripping --disable-embed-rasters --remove-metadata --strip-xml-prolog -p 9 \
		-i $< -o $@

loginscreen-bg.css: loginscreen-bg.css.in | loginscreen-bg.jpg loginscreen-bg-thumb.svg loginscreen-bg-overlay.svg
	IMAGE_DATA=$(LOGINSCREEN_BG_THUMB_BASE64_SVG) \
		$(ENVSUBST) < $< > $@

.PHONY: icons
icons: $(STATIC)/favicon.ico

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
	$(RM) loginscreen-bg.jpg
	$(RM) loginscreen-bg-thumb.jpg
	$(RM) loginscreen-bg-thumb.svg
	$(RM) loginscreen-bg.css
	$(RM) loginscreen-bg-overlay.svg
	$(RM) app-icon-*.png || true
	$(RM) app-icon.svg || true
