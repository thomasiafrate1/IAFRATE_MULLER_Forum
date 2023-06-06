var btn = document.getElementById('btn');
var closePopup = document.getElementById('close-popup');
var createBoxBtn = document.getElementById('create-box');
var popup = document.getElementById('popup');
var titleInput = document.getElementById('title');
var imgInput = document.getElementById('img');
var preview = document.getElementById('preview');
var container = document.getElementById('container');

btn.addEventListener('click', function() {
    popup.style.display = "block";
});

closePopup.addEventListener('click', function() {
    popup.style.display = "none";
});

imgInput.addEventListener('change', function() {
    preview.src = URL.createObjectURL(imgInput.files[0]);
});

createBoxBtn.addEventListener('click', function() {
    var box = document.createElement('div');
    box.className = 'box';

    var title = document.createElement('h2');
    title.textContent = titleInput.value;
    box.appendChild(title);

    var image = document.createElement('img');
    image.src = preview.src;
    box.appendChild(image);

    container.appendChild(box);
    popup.style.display = "none";
});