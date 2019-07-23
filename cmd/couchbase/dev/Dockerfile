FROM couchbase:community

ENV MEMORY_QUOTA 512
ENV INDEX_MEMORY_QUOTA 256
ENV FTS_MEMORY_QUOTA 256

ENV USERNAME "admin"
ENV PASSWORD "dev123"
ENV BUCKETS "api_client service_roles"

COPY config.sh /config-entrypoint.sh
RUN ["chmod", "+x", "/config-entrypoint.sh"]
ENTRYPOINT ["/config-entrypoint.sh"]