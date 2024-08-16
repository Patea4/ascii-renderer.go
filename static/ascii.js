document.getElementById('upload-form').addEventListener('submit', function(event) {
  event.preventDefault();

  var formData = new FormData();
  var imageFile = document.getElementById('image').files[0];
  var width = document.getElementById('width-slider').value
  formData.append('image', imageFile)
  formData.append('width', width)

  fetch('/image_upload', {
    method: 'POST',
    body: formData
  })
    .then(response => response.json())
    .then(data => {
      console.log(data)
      var asciiArtP = document.getElementById('ascii-art');

      asciiArtP.textContent = "";

      var asciiString = data.map(row => row.join('')).join('\n');

      asciiArtP.textContent = asciiString;
    })
    .catch(error => console.log("Error: ", error))
});

// Update slider value display
document.getElementById('width-slider').addEventListener('input', function() {
  var sliderValue = document.getElementById('slider-value');
  sliderValue.textContent = this.value;
});
