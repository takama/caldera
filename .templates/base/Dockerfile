FROM scratch

ENV {{[ toENV .Name ]}}_SERVER_PORT {{[ .API.Config.Port ]}}
ENV {{[ toENV .Name ]}}_INFO_PORT 8080
ENV {{[ toENV .Name ]}}_LOG_LEVEL 0

EXPOSE ${{[ toENV .Name ]}}_SERVER_PORT
EXPOSE ${{[ toENV .Name ]}}_INFO_PORT

COPY certs /etc/ssl/certs/
COPY migrations /migrations/
COPY bin/linux-amd64/service /

CMD ["/service", "serve"]
