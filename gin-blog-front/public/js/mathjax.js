window.MathJax = {
  tex: {
    // Inline math delimiters
    inlineMath: [
      ['$', '$'],
      ['\\(', '\\)'],
    ],
    // Display math delimiters
    displayMath: [
      ['$$', '$$'],
      ['\\[', '\\]'],
    ],
  },
  startup: {
    ready() {
      MathJax.startup.defaultReady()
    },
  },
}
