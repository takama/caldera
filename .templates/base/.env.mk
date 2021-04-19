{{[ toENV .Name ]}}_DB_ROOT_PASSWORD ?= {{[ randStr ]}}
{{[- if .GKE.Enabled ]}}
GKE_PROJECT_ID_DEV ?= dev-project-id
GKE_PROJECT_ID_PROD ?= prod-project-id
GKE_CLUSTER_NAME_DEV ?= dev
GKE_CLUSTER_NAME_PROD ?= prod
{{[- end ]}}
{{[ toENV .Name ]}}_SSL_CERT_NAME_DEV ?= dev-certs
{{[ toENV .Name ]}}_SSL_CERT_NAME_PROD ?= prod-certs
{{[ toENV .Name ]}}_GRPC_HOST_DEV ?= {{[ .Name ]}}-grpc.dev.host
{{[ toENV .Name ]}}_GRPC_HOST_PROD ?= {{[ .Name ]}}-grpc.prod.host
{{[ toENV .Name ]}}_REST_HOST_DEV ?= {{[ .Name ]}}-rest.dev.host
{{[ toENV .Name ]}}_REST_HOST_PROD ?= {{[ .Name ]}}-rest.prod.host
{{[ toENV .Name ]}}_CORS_ALLOWED_HOSTS ?= localhost,127.0.0.1
