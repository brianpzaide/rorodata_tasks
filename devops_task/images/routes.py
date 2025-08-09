from flask import Blueprint, render_template, request

images_bp = Blueprint('images', __name__, template_folder='templates')

@images_bp.route('/images', methods=['GET'])
def get_images():
    return render_template('images/images_list.html', images=[])

@images_bp.route('/images/<image_name>', methods=['DELETE'])
def delete_image(image_name):
    return render_template('images/images_list.html', images=[])

@images_bp.route('/images', methods=['POST'])
def pull_image():
    image_name = request.json.get('image_name')
    return render_template('images/images_list.html', images=[])