
mds = $(wildcard *.md)
htmls = $(mds:.md=.html)

default: $(htmls) pfs.1
	@true

%.html: %.md
	markdown $< > $@

pfs.1: pfs_man_v0.1.md
	echo "build manual page -- FIXME"
	touch $@

