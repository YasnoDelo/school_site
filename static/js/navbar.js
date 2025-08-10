document.addEventListener('DOMContentLoaded', function () {
  var btn = document.getElementById('nav-toggle');
  var menu = document.querySelector('.nav-menu');
  if (!btn || !menu) return;

  var prevFocus = null;
  btn.setAttribute('aria-expanded', 'false');

  // scroll lock helpers (preserve scroll position)
  var scrollAttr = 'data-scroll-y';
  function lockScroll() {
    var scrollY = window.scrollY || window.pageYOffset || document.documentElement.scrollTop || 0;
    document.body.setAttribute(scrollAttr, String(scrollY));
    document.body.style.position = 'fixed';
    document.body.style.top = '-' + scrollY + 'px';
    // compensate scrollbar width
    var scrollbarWidth = window.innerWidth - document.documentElement.clientWidth;
    if (scrollbarWidth > 0) document.body.style.paddingRight = scrollbarWidth + 'px';
    document.body.classList.add('menu-open');
  }
  function unlockScroll() {
    var saved = document.body.getAttribute(scrollAttr);
    var scrollY = saved ? parseInt(saved, 10) : 0;
    document.body.style.position = '';
    document.body.style.top = '';
    document.body.style.paddingRight = '';
    document.body.classList.remove('menu-open');
    window.scrollTo(0, scrollY);
    document.body.removeAttribute(scrollAttr);
  }

  function openMenu() {
    prevFocus = document.activeElement;
    menu.classList.add('open');
    menu.setAttribute('aria-hidden', 'false');
    btn.setAttribute('aria-expanded', 'true');
    lockScroll();
    // focus on first link
    var first = menu.querySelector('a, button');
    if (first) first.focus();
  }

  function closeMenu() {
    menu.classList.remove('open');
    menu.setAttribute('aria-hidden', 'true');
    btn.setAttribute('aria-expanded', 'false');
    unlockScroll();
    if (prevFocus && prevFocus.focus) prevFocus.focus();
  }

  btn.addEventListener('click', function (ev) {
    if (menu.classList.contains('open')) closeMenu(); else openMenu();
  });

  // close on click outside
  document.addEventListener('click', function (ev) {
    if (!menu.classList.contains('open')) return;
    var t = ev.target;
    if (t === btn || btn.contains(t) || menu.contains(t)) return;
    closeMenu();
  });

  // close on Esc
  document.addEventListener('keydown', function (ev) {
    if (ev.key === 'Escape' && menu.classList.contains('open')) {
      ev.preventDefault();
      closeMenu();
    }
  });
});
