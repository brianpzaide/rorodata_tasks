from service import service
from urllib.parse import unquote
from flask import Blueprint, render_template, request, jsonify

containers_bp = Blueprint('containers', __name__, template_folder='templates')

@containers_bp.route('/containers', methods=['GET'])
def list_containers():
    image_name = request.args.get('image_name')
    all_flag = request.args.get('all', '')
    print(all_flag)
    all = all_flag == 'true' or all_flag == '1'
    containers = []
    try:
        if (not image_name) or image_name.strip().lower() == "":
            containers = service.list_containers(all=all_flag)
        else:
            containers = service.list_containers(image_name.strip().lower(), all=all)
        return render_template('containers/containers_list.html', containers=containers)
    except service.ServerError as e:
        return render_template('containers/containers_list.html', err=e)
    

@containers_bp.route('/containers/<container_id>', methods=['DELETE'])
def delete_container(container_id):
    try:
        if (not container_id) or container_id.strip().lower() == "":
            return render_template('containers/containers_list.html', err="client error: container_id must be a valid string")
        else:
            service.delete_container(container_id.strip().lower())
            containers = service.list_containers(all=True)
            return render_template('containers/containers_list.html', containers=containers)
    except service.ServerError as e:
        return render_template('containers/containers_list.html', err=e)

@containers_bp.route('/containers/<container_id>', methods=['PUT'])
def stop_container(container_id):
    try:
        if (not container_id) or container_id.strip().lower() == "":
            return render_template('containers/containers_list.html', err="client error: container_id must be a valid string")
        else:
            service.stop_container(container_id.strip().lower())
            containers = service.list_containers(all=True)
            return render_template('containers/containers_list.html', containers=containers)
    except service.ServerError as e:
        return render_template('containers/containers_list.html', err=e)

@containers_bp.route('/containers', methods=['POST'])
def run_container():
    data = request.get_json()
    try:
        print(data)
        service.run_container(name=data['name'], image_name=data['image_name'], commands=[unquote(c) for c in data['commands']], ports=data['ports'])
        print(data)
    except KeyError:
        return '', 400
    except service.ServerError as e:
        return jsonify({"error": str(e)}), 500
    return '', 204

@containers_bp.route('/containers/<container_id>', methods=['GET'])
def container_detail(container_id):
    try:
        if (not container_id) or container_id.strip().lower() == "":
            return render_template('containers/containers_list.html', err="client error: container_id must be a valid string")
        else:
            container_details = service.get_container_details(container_id.strip().lower())
            return render_template('containers/container_details.html', container_details=container_details)
    except service.ServerError as e:
        return render_template('containers/containers_list.html', err=e)

    