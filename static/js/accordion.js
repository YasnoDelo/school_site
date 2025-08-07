document.addEventListener('DOMContentLoaded', function() {
    var toggles = document.querySelectorAll('.section-toggle');
    for (var i = 0; i < toggles.length; i++) {
      (function(btn) {
        btn.addEventListener('click', function() {
          var section = btn.parentNode;
          if (section.classList.contains('open')) {
            section.classList.remove('open');
          } else {
            section.classList.add('open');
          }
        });
      })(toggles[i]);
    }
  });
  