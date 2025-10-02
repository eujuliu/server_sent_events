<a name="readme-top"></a>

# SSE Service

<p align="center">
  A Server-Sent Events (SSE) service that receives events via a queue and broadcasts them to connected clients.
  <br />
  <a href="https://github.com/eujuliu/taskscheduler-sse/issues">Have a question?</a>
  ·
  <a href="https://github.com/eujuliu/taskscheduler-sse/fork">Request Feature</a>
</p>

<ul>
  <li>
    <a href="#technologies">Technologies</a>
  </li>
  <li>
    <a href="#getting-started">Getting Started</a>
    <ul>
      <li><a href="#prerequisites">Prerequisites</a></li>
      <li><a href="#installation">Installation</a></li>
      <li><a href="#usage">Usage</a></li>
    </ul>
  </li>
  <li><a href="#contributing">Contributing</a></li>
  <li><a href="#author">Author</a></li>
</ul>

## Technologies

This project was built using the `Go` programming language. It utilizes `Docker` and `Docker Compose` for containerization and development environments. For code quality, `golangci-lint` is used for linting and formatting. Development is facilitated by `air` for hot reloading. Testing is handled with Go's built-in testing framework, and additional tools like `pre-commit` for hooks.

## Getting Started

### Prerequisites

For running this project, you will need the following:

- [Go](https://golang.org/dl/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [GNU Make](https://www.gnu.org/software/make/) (usually pre-installed on Linux)

And you need to download the project to your computer.

### Installation

After fulfilling all requirements, you need to install the dependencies. If using Go modules, run:

```bash
go mod download
```

### Usage

To use the project, you can run the following commands using the provided Makefile:

- To build the project:

  ```bash
  make build
  ```

- To run in development mode with hot reloading:

  ```bash
  make dev
  ```

- To run the built binary:

  ```bash
  make run
  ```

- To run tests:

  ```bash
  make test
  ```

- For development with Docker Compose watching:

  ```bash
  make watch
  ```

## Contributing

If you'd like to contribute to this project, please follow these steps:

1. Fork this repository.
2. Create a branch: `git checkout -b feat/your-feature`.
3. Make your changes and commit them: `git commit -m 'Add some feature'`.
4. Push to the original branch: `git push origin feat/your-feature`.
5. Create a pull request.

## Author

<img style="border-radius: 50%;" src="https://avatars.githubusercontent.com/u/49854105?v=4" width="100px;" alt=""/>
<br />
<sub><b>Julio Martins</b></sub></a>

Made by Julio Martins 👋🏽 Contact me!

[![Linkedin Badge](https://img.shields.io/badge/-@ojuliomartins-1262BF?style=for-the-badge&labelColor=1262BF&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/ojuliomartins/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>
