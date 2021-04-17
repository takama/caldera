FROM scratch

# Common configuration
{{[- if .API.Enabled ]}}
ENV {{[ toENV .Name ]}}_SERVER_PORT {{[ .API.Config.Port ]}}
{{[- if .API.Gateway ]}}
ENV {{[ toENV .Name ]}}_SERVER_GATEWAY_PORT {{[ .API.Config.Gateway.Port ]}}
{{[- end ]}}
{{[- end ]}}
ENV {{[ toENV .Name ]}}_INFO_PORT 8080
ENV {{[ toENV .Name ]}}_LOGGER_LEVEL 0

# Exposing ports
{{[- if .API.Enabled ]}}
EXPOSE ${{[ toENV .Name ]}}_SERVER_PORT
{{[- if .API.Gateway ]}}
EXPOSE ${{[ toENV .Name ]}}_SERVER_GATEWAY_PORT
{{[- end ]}}
{{[- end ]}}
EXPOSE ${{[ toENV .Name ]}}_INFO_PORT

# Copy dependecies
COPY certs /etc/ssl/certs/
{{[- if .Storage.Enabled ]}}
{{[- if .Example ]}}
COPY fixtures /fixtures/
{{[- end ]}}
COPY migrations /migrations/
{{[- end ]}}
COPY bin/linux-amd64/{{[ .Bin ]}} /

CMD ["/{{[ .Bin ]}}", "serve"]
