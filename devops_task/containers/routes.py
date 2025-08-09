from flask import Blueprint, render_template, request, jsonify

containers_bp = Blueprint('containers', __name__, template_folder='templates')

@containers_bp.route('/containers', methods=['GET'])
def list_containers():
    image_name = request.args.get('image_name')
    all_flag = request.args.get('all', '0')
    return render_template('containers/containers_list.html', containers=[])

@containers_bp.route('/containers/<container_id>', methods=['DELETE'])
def delete_container(container_id):
    return render_template('containers/containers_list.html', containers=[])

@containers_bp.route('/containers', methods=['POST'])
def create_container():
    data = request.get_json()
    return '', 204

@containers_bp.route('/containers/<container_id>', methods=['GET'])
def container_detail(container_id):
    return jsonify({
        "container_id": container_id,
        "image": "",
        "ports_exposed": [],
        "commands_specified_by_the_user_to_run": "",
        "logs": ""
    })