var btn = document.getElementById('btn');
        var container = document.getElementById('container');
        
        btn.addEventListener('click', function() {
            var box = document.createElement('div');
            box.className = 'box';

            var title = document.createElement('h2');
            title.textContent = 'Titre de la box';
            box.appendChild(title);

            var image = document.createElement('img');
            image.src = 'https://via.placeholder.com/150'; // Remplacez par l'URL de votre image
            box.appendChild(image);

            container.appendChild(box);
        });