from service import service
from urllib.parse import unquote
from flask import Blueprint, render_template, request

images_bp = Blueprint('images', __name__, template_folder='templates')

@images_bp.route('/images', methods=['GET'])
def get_images():
    try:
        images = service.list_images()
        return render_template('images/images_list.html', images=images)
    except service.ServerError as e:
        return render_template('images/images_list.html', err=e)    
    

@images_bp.route('/images/<image_name>', methods=['DELETE'])
def delete_image(image_name):
    image_name = unquote(image_name).strip().lower()
    try:
        if not image_name:
            return render_template('images/images_list.html', err="client error: image_name must be a valid string")
        else:
            service.delete_image(image_name=image_name)
            images = service.list_images()
            return render_template('images/images_list.html', images=images)
    except service.ServerError as e:
        return render_template('images/images_list.html', err=e)

@images_bp.route('/images', methods=['POST'])
def pull_image():
    image_name = request.form.get('image_name') or {}
    try:
        if not image_name or image_name.strip() == "":
            return render_template('images/images_list.html', err="client error: image_name must be a valid string")
        else:
            service.pull_image(image_name=image_name)
            images = service.list_images()
            return render_template('images/images_list.html', images=images)
    except service.ServerError as e:
        return render_template('images/images_list.html', err=e)
    