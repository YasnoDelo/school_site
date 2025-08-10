// static/js/mathjax-config.js
window.MathJax = {
    tex: {
      inlineMath: [['$', '$'], ['\\(', '\\)']],
      displayMath: [['$$', '$$'], ['\\[', '\\]']],
      macros: { RR: "\\mathbb{R}" }
    },
    options: {
      skipHtmlTags: ['script','noscript','style','textarea','pre']
    }
  };
  