// static/js/navbar.js
document.addEventListener('DOMContentLoaded', function() {
    var toggle = document.getElementById('nav-toggle');
    var menu = document.getElementById('primary-menu');
    var html = document.documentElement;
  
    if (!toggle || !menu) return;
  
    toggle.addEventListener('click', function() {
      var expanded = toggle.getAttribute('aria-expanded') === 'true';
      if (expanded) {
        toggle.setAttribute('aria-expanded', 'false');
        toggle.setAttribute('aria-label', 'Открыть меню');
        if (menu.classList) menu.classList.remove('open');
        menu.setAttribute('hidden', '');
        // allow scroll
        html.classList.remove('nav-open');
      } else {
        toggle.setAttribute('aria-expanded', 'true');
        toggle.setAttribute('aria-label', 'Закрыть меню');
        if (menu.classList) menu.classList.add('open');
        menu.removeAttribute('hidden');
        // optionally prevent background scroll when menu open:
        html.classList.add('nav-open');
      }
    });
  
    // close menu when clicking outside (mobile)
    document.addEventListener('click', function(e) {
      if (!menu.contains(e.target) && !toggle.contains(e.target) && menu.getAttribute('hidden') === null) {
        // menu is open but click outside
        toggle.setAttribute('aria-expanded', 'false');
        toggle.setAttribute('aria-label', 'Открыть меню');
        if (menu.classList) menu.classList.remove('open');
        menu.setAttribute('hidden', '');
        html.classList.remove('nav-open');
      }
    });
  
    // close on Escape
    document.addEventListener('keydown', function(e) {
      if (e.key === 'Escape' || e.key === 'Esc') {
        toggle.setAttribute('aria-expanded', 'false');
        menu.setAttribute('hidden', '');
        if (menu.classList) menu.classList.remove('open');
        html.classList.remove('nav-open');
      }
    });
  });
  