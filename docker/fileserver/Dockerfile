# Build context should be project root
FROM scratch
COPY docker/fileserver/data build/fileserver /
ENTRYPOINT ["/fileserver"]
CMD ["/data"]