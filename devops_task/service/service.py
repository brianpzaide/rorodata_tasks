
from typing import List, Dict, Union
import docker
from threading import Lock

_docker_client = None
_docker_client_lock = Lock()

def get_docker_client():
    global _docker_client
    with _docker_client_lock:
        if _docker_client is None:
            _docker_client = docker.DockerClient(base_url='unix://var/run/docker.sock')
    return _docker_client

class ServerError(Exception):
    pass

def list_images()->List[str]:
    dc = get_docker_client()
    try:
        return [tag for img in dc.images.list() for tag in img.tags]
    except docker.errors.APIError as e:
         raise ServerError("Failed to list images") from e
    
def list_containers(image_name = "", all=False) -> List[str]:
    dc = get_docker_client()
    try:
        _containers = dc.containers.list(all=all) if image_name == "" else dc.containers.list(filters={
            "ancestor": f"{image_name}:latest"
        })
        containers = []
        for c in _containers:
            containers.append({
                "id": c.id[:6],
                "image": c.image.tags[0],
                "status": c.status,
                "action": "stop" if c.status == "running" else "remove"
            })
        return containers
    except docker.errors.APIError as e:
         raise ServerError(f"Failed to pull the image: {image_name}") from e

def pull_image(image_name: str):
    dc = get_docker_client()
    try:
        dc.images.pull(image_name)
    except docker.errors.APIError as e:
         raise ServerError(f"Failed to pull the image: {image_name}") from e

def delete_image(image_name):
    dc = get_docker_client()
    try:
        dc.images.remove(image=image_name)
    except docker.errors.APIError as e:
         raise ServerError(f"Failed to remove the image: {image_name}") from e


def get_container_details(container_id: str):
    dc = get_docker_client()
    try:
        container = dc.containers.get(container_id=container_id)

        logs = [log for log in container.logs().decode('ascii').split("\n") if log != ""]
        
        
        ports = []
        for k, v in container.ports.items():
            if v:
                ports.append(f"{k}:{v[0]['HostPort']}")
            else:
                ports.append(f"{k}:unmapped")
        
        container_details = {
            "id": container.id,
            "name": container.name,
            "image": container.image.tags[0],
            "ports": ports,
            "logs": logs,
            "status": container.status
        }
        return container_details
    except docker.errors.APIError as e:
         raise ServerError(f"Failed to fetch the container details: {container_id}") from e
    

def run_container(name:str, image_name: str, commands: Union[str, List[str]], ports: Dict[str, str]):
    ports = {f"{k}/tcp": v for k, v in ports.items()}
    dc = get_docker_client()
    try:
        commands = ["/bin/sh", "-c", " && ".join(c for c in commands)]
        dc.containers.run(name= name,image=image_name, command=commands, ports=ports, detach=True)
    except docker.errors.APIError as e:
         print(e)
         raise ServerError(f"Failed to run a container with image: {image_name}") from e

def stop_container(container_id: str):
    dc = get_docker_client()
    try:
        container = dc.containers.get(container_id=container_id)
        container.stop()
    except docker.errors.APIError as e:
         raise ServerError(f"Failed to stop the container with id: {container_id}") from e

def delete_container(container_id: str):
    dc = get_docker_client()
    try:
        container = dc.containers.get(container_id=container_id)
        if container.status == "running":
            container.stop()
        container.remove()
    except docker.errors.APIError as e:
         raise ServerError(f"Failed to delete the container with id: {container_id}") from e
