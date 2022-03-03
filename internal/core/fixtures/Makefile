# In order to stay consistent with upstream, we copy over the fixture files from Limes.
#
# Run `make all` and then adjust code and unit tests accordingly.

fixtures = cluster-get-west \
           domain-get-germany domain-list-filtered domain-list \
           project-get-berlin-only-rates project-get-dresden project-list-filtered project-list

all: $(fixtures)

$(fixtures): %: FORCE
	curl -so $@.json "https://raw.githubusercontent.com/sapcc/limes/master/pkg/api/fixtures/$@.json"

.PHONY: FORCE
