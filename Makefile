VERSION = v0.0.1

release:
	git tag -a ${VERSION} -m ${VERSION}
	git push origin ${VERSION}
	goreleaser release --rm-dist