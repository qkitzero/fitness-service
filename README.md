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

    namespace tenant {
        class Profile["Profile テナント情報"] {
            TenantID tenantID
            address.Address address
            contact.Phone? phone
            contact.Email? email
            HomepageURL? homepageURL
            Note? note
        }
    }

    namespace organization {
        class Organization["Organization 組織"] {
            OrganizationID id
            tenant.TenantID tenantID
            Name name
        }
    }

    namespace customer {
        class Customer["Customer 顧客"] {
            CustomerID id
            tenant.TenantID tenantID
            Name name
            NameKana nameKana
            Gender gender
            BirthDate birthDate
            contact.Phone? phone
            contact.Email? email
            address.Address address
            EmergencyContactName? emergencyContactName
            EmergencyContactRelationship? emergencyContactRelationship
            contact.Phone? emergencyContactPhone
            organization.OrganizationID? organizationID
            bool active
        }
    }

    namespace measurement {
        class Measurement["Measurement 測定"] {
            MeasurementID id
            customer.CustomerID customerID
            MeasuredOn measuredOn
            staff.StaffID measuredBy
            AgeAtMeasurement ageAtMeasurement
            staff.StaffID updatedBy
            bool isDraft
        }

        class MeasurementEntry["MeasurementEntry 測定エントリ"] {
            measurementitem.MeasurementItemID measurementItemID
            bool unmeasurable
            Note? note
        }

        class MeasurementValue["MeasurementValue 測定値"] {
            TrialIndex trialIndex
            Side side
            Value? value
            Value? valueSecondary
            Choice? valueChoice
        }
    }

    namespace measurementitem {
        class MeasurementItem["MeasurementItem 測定項目"] {
            MeasurementItemID id
            Code code
            Name name
            Category category
            Unit unit
            TrialCount trialCount
            SideMode sideMode
            ValueType valueType
            ScoreDirection? scoreDirection
            SideAggregation sideAggregation
            Element[] elements
        }
    }

    namespace standard {
        class AgeGroupStandard["AgeGroupStandard 年代別基準値"] {
            AgeGroupStandardID id
            measurementitem.MeasurementItemID measurementItemID
            Gender gender
            AgeRange ageRange
            Mean mean
            StandardDeviation standardDeviation
        }

        class RankStandard["RankStandard ランク基準"] {
            Rank rank
            ZScore? zScoreMin
            ZScore? zScoreMax
        }
    }

    namespace judgment {
        class Judgment["Judgment 判定"] {
            measurement.MeasurementID measurementID
            Advice? advice
        }

        class Evaluation["Evaluation 評価"] {
            <<computed, not persisted>>
            MotorAge? motorAge
        }

        class ItemEvaluation["ItemEvaluation 項目別評価"] {
            measurementitem.MeasurementItemID measurementItemID
            measurement.Value value
            standard.Mean mean
            standard.ZScore zScore
            standard.Rank rank
        }

        class ElementEvaluation["ElementEvaluation 体力要素別評価"] {
            measurementitem.Element element
            standard.ZScore zScore
            standard.Rank rank
        }

        class Prescription["Prescription 処方"] {
            <<computed, not persisted>>
        }

        class PrescribedMenu["PrescribedMenu 処方メニュー"] {
            PrescriptionSource source
            measurementitem.Element? element
            training.Part? part
            training.TrainingMenuID trainingMenuID
            training.Name trainingMenuName
            training.Amount amount
            training.Unit unit
            training.Sets sets
        }

        class PrescribedMenuOverride["PrescribedMenuOverride 処方メニューの手動上書き"] {
            PrescribedMenuOverrideID id
            measurement.MeasurementID measurementID
            training.SortOrder sortOrder
            measurementitem.Element? element
            training.Part? part
            training.TrainingMenuID trainingMenuID
            training.Amount amount
            training.Unit unit
            training.Sets sets
        }
    }

    namespace training {
        class TrainingMenu["TrainingMenu トレーニングメニュー"] {
            TrainingMenuID id
            Code code
            Name name
            measurementitem.Element element
            Part part
            Amount amount
            Unit unit
            Sets sets
            Instruction instruction
        }

        class ElementMenu["ElementMenu 体力要素別メニュー"] {
            measurementitem.Element element
            Part part
            Level level
            TrainingMenuID trainingMenuID
        }

        class AgeDecadeMenu["AgeDecadeMenu 年代別メニュー"] {
            Decade decade
            SortOrder sortOrder
            TrainingMenuID trainingMenuID
        }

        class FixedMenu["FixedMenu 固定メニュー"] {
            SortOrder sortOrder
            TrainingMenuID trainingMenuID
        }
    }

    Measurement "1" *-- "0..*" MeasurementEntry
    MeasurementEntry "1" *-- "0..*" MeasurementValue
    Evaluation "1" *-- "0..*" ItemEvaluation
    Evaluation "1" *-- "0..*" ElementEvaluation
    Prescription "1" *-- "0..*" PrescribedMenu

    Customer ..> Organization : organizationID
    Measurement ..> Customer : customerID
    MeasurementEntry ..> MeasurementItem : measurementItemID
    AgeGroupStandard ..> MeasurementItem : measurementItemID
    Judgment ..> Measurement : measurementID
    PrescribedMenuOverride ..> Measurement : measurementID
    PrescribedMenuOverride ..> TrainingMenu : trainingMenuID
    ItemEvaluation ..> MeasurementItem : measurementItemID
    PrescribedMenu ..> TrainingMenu : trainingMenuID
    Evaluation ..> Measurement : input
    Evaluation ..> MeasurementItem : input
    Evaluation ..> AgeGroupStandard : input
    Evaluation ..> RankStandard : input
    Prescription ..> Evaluation : input
    Prescription ..> ElementMenu : input
    Prescription ..> FixedMenu : input
    Prescription ..> AgeDecadeMenu : input
    Prescription ..> TrainingMenu : input
    Prescription ..> PrescribedMenuOverride : input
    ElementMenu ..> TrainingMenu : trainingMenuID
    AgeDecadeMenu ..> TrainingMenu : trainingMenuID
    FixedMenu ..> TrainingMenu : trainingMenuID
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
