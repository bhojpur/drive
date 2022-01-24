FROM moby/buildkit:v0.9.3
WORKDIR /swarm
COPY swarm README.md /swarm/
ENV PATH=/swarm:$PATH
ENTRYPOINT [ "/bhojpur/swarm" ]