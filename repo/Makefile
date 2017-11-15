# The exact list is `ls packages/local` + $PKGS
PKGS = alcotest depext

.PHONY: commit update-osx list-licenses

all: commit
	@

update-osx:
	./update-upstream.sh $(PKGS)

commit:
	$(MAKE) update-osx
	git add darwin/packages/upstream
	git commit . -s -m "AUTO: Update upstream packages"

list-licenses:
	mkdir -p licenses
	cd licenses && ../opam-licenses.sh hyperkit
	./list-licenses.sh licenses > OCAML-LICENSES
