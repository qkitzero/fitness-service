# Fitness Service

[![release](https://img.shields.io/github/v/release/qkitzero/fitness-service?logo=github)](https://github.com/qkitzero/fitness-service/releases)
[![test](https://github.com/qkitzero/fitness-service/actions/workflows/test.yml/badge.svg)](https://github.com/qkitzero/fitness-service/actions/workflows/test.yml)
[![Lint](https://github.com/qkitzero/fitness-service/actions/workflows/lint.yml/badge.svg)](https://github.com/qkitzero/fitness-service/actions/workflows/lint.yml)
[![codecov](https://codecov.io/gh/qkitzero/fitness-service/graph/badge.svg)](https://codecov.io/gh/qkitzero/fitness-service)
[![Buf CI](https://github.com/qkitzero/fitness-service/actions/workflows/buf-ci.yaml/badge.svg)](https://github.com/qkitzero/fitness-service/actions/workflows/buf-ci.yaml)

- Microservices Architecture
- gRPC
- gRPC Gateway
- Buf ([buf.build/qkitzero-org/fitness-service](https://buf.build/qkitzero-org/fitness-service))
- Clean Architecture
- Docker
- Test
- Codecov
- Cloud Build
- Cloud Run

```mermaid
classDiagram
    direction LR

    class Customer {
        id
        name
        createdAt
        updatedAt
    }
```

```mermaid
flowchart TD
    subgraph gcp[GCP]
        secret_manager[Secret Manager]

        subgraph cloud_build[Cloud Build]
            build_fitness_service(Build fitness-service)
            push_fitness_service(Push fitness-service)
            deploy_fitness_service(Deploy fitness-service)

            build_fitness_service_gateway(Build fitness-service-gateway)
            push_fitness_service_gateway(Push fitness-service-gateway)
            deploy_fitness_service_gateway(Deploy fitness-service-gateway)
        end


        subgraph artifact_registry[Artifact Registry]
            fitness_service_image[(fitness-service image)]
            fitness_service_gateway_image[(fitness-service-gateway image)]
        end

        subgraph cloud_run[Cloud Run]
            fitness_service(Fitness Service)
            fitness_service_gateway(Fitness Service Gateway)
        end
    end

    subgraph external[External]
        fitness_db[(Fitness DB)]
    end

    build_fitness_service --> push_fitness_service --> fitness_service_image
    build_fitness_service_gateway --> push_fitness_service_gateway --> fitness_service_gateway_image

    fitness_service_image --> deploy_fitness_service --> fitness_service
    fitness_service_gateway_image --> deploy_fitness_service_gateway --> fitness_service_gateway

    secret_manager --> deploy_fitness_service
    secret_manager --> deploy_fitness_service_gateway

    fitness_service_gateway --> fitness_service
    fitness_service --> fitness_db
```
