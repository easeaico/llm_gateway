
build:
	docker build . -t llm-mesh

run:
	docker run -d -p 5984:5984 llm-mesh

package:
	docker build . -t llm-mesh:amd64 --platform linux/amd64
	docker save llm-mesh:amd64 -o llm-mesh.tar
