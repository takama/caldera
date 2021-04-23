# Used defaut namespace for the environment
# Namespace: public, commerce, billing,
NAMESPACE ?= default

# Cluster: dev, prod, ...
CLUSTER ?= dev

# Some database engines require it
{{[ toENV .Name ]}}_DB_ROOT_PASSWORD ?= {{[ randStr ]}}

{{[- if .GKE.Enabled ]}}

# GKE environments initialisation
GKE_CLUSTER_NAME = $(CLUSTER)
GKE_PROJECT_ID = $(GKE_CLUSTER_NAME)-project-id
GKE_PROJECT_REGION ?= {{[ .GKE.Region ]}}
{{[- end ]}}

# SSL/hosts environments initialisation
{{[ toENV .Name ]}}_SSL_CERT_NAME ?= $(CLUSTER)-certs
{{[ toENV .Name ]}}_GRPC_HOST ?= {{[ .Name ]}}-grpc.$(CLUSTER).host
{{[ toENV .Name ]}}_REST_HOST ?= {{[ .Name ]}}-rest.$(CLUSTER).host
{{[ toENV .Name ]}}_CORS_ALLOWED_HOSTS ?= $({{[ toENV .Name ]}}_REST_HOST),localhost,127.0.0.1
