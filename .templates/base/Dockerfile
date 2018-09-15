FROM scratch

# Common configuration
{{[- if .API.Enabled ]}}
ENV {{[ toENV .Name ]}}_SERVER_PORT {{[ .API.Config.Port ]}}
{{[- end ]}}
ENV {{[ toENV .Name ]}}_INFO_PORT 8080
ENV {{[ toENV .Name ]}}_LOG_LEVEL 0

# Exposing ports
{{[- if .API.Enabled ]}}
EXPOSE ${{[ toENV .Name ]}}_SERVER_PORT
{{[- end ]}}
EXPOSE ${{[ toENV .Name ]}}_INFO_PORT

# Copy dependecies
COPY certs /etc/ssl/certs/
{{[- if .Storage.Enabled ]}}
COPY migrations /migrations/
{{[- end ]}}
COPY bin/linux-amd64/service /

CMD ["/service", "serve"]
