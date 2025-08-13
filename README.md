# Rorodata Interview Tasks

This repository contains my solutions to two technical tasks from [Rorodata's careers page](https://github.com/rorodata/careers).

## 1. [Backend2 Task](https://github.com/rorodata/careers/blob/master/backend2.md)

**Language:** Go  
**Summary:** A RESTful API for managing clusters and machines in a cloud-like environment (without actual cloud integration) to
- Create and delete clusters (with name and cloud region).
- Create and delete machines within clusters.
- Assign tags to machines at creation.
- Perform actions (start, stop, reboot) on groups of machines based on tags.

Implements the repository pattern for clean separation of data access logic, making it easy to support other databases in the future.

---

## 2. [DevOps Task](https://github.com/rorodata/careers/blob/master/devops.md)

**Stack:** Python, Flask, Jinja2, HTMX, [Docker SDK for Python](https://pypi.org/project/docker/)

**Summary:** A web-based dashboard to interact with Docker, visualize containers and images, and perform common Docker operations such as
- List all containers and available images.
- View details of individual containers such as logs.
- Start new containers with configurable commands and port mappings.
- Stop or remove containers.
- Add or remove Docker images.

