// static/js/gallery.js
document.addEventListener("DOMContentLoaded", function() {
    // Опции по вкусу:
    if (window.lightbox) {
      lightbox.option({
        'resizeDuration': 200,
        'wrapAround': true,
        'albumLabel': 'Изображение %1 из %2'
      });
    }
  });
  