# Video Uploader API - Description

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Description](#description)
  - [Microservices - Orders](#microservice---orders)
- [Architecture](#architecture)
  - [DDD](#ddd)
  - [Clean Archtecture](#clean-archtecture)
- [Design Patterns](#design-patterns)
- [Unit Testing](#unit-testing)
  - [BDD](#bdd)
- [Docker build and run](#docker-build-and-run)
- [How to use](#how-to-use)
  - [Check app status](#check-app-status)
- [AWS](#aws)
- [Kubernetes](#kubernetes)
- [API endpoints](#api-endpoints)
- [Upload videos](#upload-video)


## Description

The Hackathon Video Uploader 1 aims to do a solution for a video processing and zipping print images. With this software, the user can upload ad many videos as she/he wants and get, after a while, a ZIP URL.

All the Endpoints can be called by accessing `http://localhost:3210/api` API url.

To build and run this project. Follow the Docker section

### Microservice - Orders

The Project is divided in `2 microservices`. Only one has Database. These 2 MSs are stored in `EKS Cluster`. The microservices are:

- API Video Uploader
- Video Processing

## Architecture

This project uses a mixin with two architectures to make it scalable and secure by protecting the **domain** layer. Those are `DDD (Domain Driven Design)` and `Clean Archtecture`

### DDD ###

To design this application was chosen `DDD (Domain Drive Design)` architecture to follow the principle of **protecting the model**.

![ddd_image](https://github.com/thiagoluis88git/tech1-orders/assets/166969350/2016bfff-3c19-4172-837f-8d5d428525f7)

### Clean Archtecture ###

The other one is `Clean Archtecture`. With it, we add some extra layers to organize even more the project.

![CleanArchitecture](https://github.com/user-attachments/assets/a49c2aab-562c-4b6c-82f2-7ffe9e4aec74)

The folder project was created to follow this main principle:

- **data**: Here we have all the implementations, such as Repositories, Remotes and Locals
- **domain**: All the `application business logic`, such as UseCases
- **handler**: Also known as `presentation layer` resides all the `Controllers` handled by the `Web` **interface** given by the `/cmd/api Framework & Driver` 

- **integrations**: This is not part of `DDD` or `Clean Arch`, but is important separate some external packages or integrations, such `Mercado Livre API`

## Design Patterns

To improve and make a good standard project pattern, some `Design Patterns` were used in this application.

- Strategy: All the business logic must be protected by the external implementations. To do it, we use a combo with **Interfaces** and **Dependency Inversion solid principle** to inject only *interfaces* and not *real implementations*
- Dependency Injection: Is used in application bootstrap (main.go) to inject all the interfaces implementations.
- Decorator: To inject **Services** inside **Driver Adapter** *handler*. By doing this, we *decorate* the Handler with a Service
- Services or Use Cases: Centralize all the business logic of of the application
- Repository: Used to integrate with all **Driven Adapter** like *Databases and External Endpoints*

## Unit Testing

To run all the Unit Testing for this project, just run:

```
go test -cover ./... -coverprofile="cover.out"
go tool cover -func="cover.out"
```

### BDD

Inside `bdd` folder, has an implementation of a BDD test. This test is made by a [Cucumber API](https://github.com/cucumber/godog).
The **BDD** tests will be triggered when running the `go test ./...` in the previous step

```
go test -cover ./... -coverprofile="cover.out"
go tool cover -func="cover.out"
```

This will run all the **Services** unit tests and **Repository** unit Database tests running [Testcontainers](https://testcontainers.com/) database container mocks.

## Docker build and run

This project was built using Docker and Docker Compose. So, to build and run the API, we need to run in the root of the project:

```
docker compose build
```

After the image build finish, run:

```
docker compose up -d
```

The command above may take a while.

After all the containers shows these below status:

```
 ✔ Container hack-video-uploader         Started
 ✔ Container hack-video-processing       Started 
```

we can access `http://localhost:3210/api` endpoints.


## How to use

To use all the endpoints in this API, we can follow these sequence to simulate a customer making an order in a restaurant.
We can separate in three moments.

- Restaurant products manipulation. This is used by the `restaurant owner` to create all the product portfolio with its images and prices
- Customer self service. This is used by the `customer` to choose the products, pay for it and create an order 
- Order preparing and deliver. This is used by the `chef` and `waiter` to check the order status

We will divide in 2 sections: **Restaurant owner** and **Customer order**

### Check app status

After running `Docker` commands, you can check the application status running:

```
docker compose logs app
```

We can see some database errors but at the end of the logs we can see:

```
hack-video-uploader  | 2024/05/27 22:57:35 Video Uploader API has started
```

## AWS ##

The Fast food project uses `AWS Cloud` to host its software components. To know more about the **AWS configuration**, read: [AWS Readme](https://github.com/thiagoluis88git/hack-k8s/infra/README.md)

## Kubernetes

This application has all the K8S YAMLs to be applied in any cluster. 
To read the specific documentation, read: [Kubernetes README](https://github.com/thiagoluis88git/hack-k8s/blob/main/infra/k8s/README.md)

## API endpoints

This section will be used by the API to use this solution


- Cal the POST `http://localhost:3210/auth/login` to Login
- Cal the POST `http://localhost:3210/auth/signup` to Create a user
- Cal the POST `http://localhost:3210/api/upload` to Upload directly a video via Multipart Data
- Cal the GET `http://localhost:3210/api/upload/presign/{cpf}` Get a Presign URL for Upload apart of API Gateway
- Cal the PUT `http://localhost:3210/api/upload/tracking/{trackingID}` Send to proccess a video stored in Presign URL
- Cal the GET `http://localhost:3210/api/trackings/{cpf}` get a list of trackings of a user


## Upload Video

To upload a video, the user can do by two options

- 1. Via **/api/upload**
- 2. Via **Presign a URL**, upload the file by itself and use **upload/tracking/{trackingID}** to start the Video Proccess.

The first option has some drawbacks, for example, it has **10mb** *API Gateway* limits. And it is slower than the second option.
The second option the user has no limit for the upload a video. The user needs to Presign a URL via `presign` endpoint, upload the video via some application and then start the video proccessing by using `put tracking/{trackingID} endpoint`.
