
function showVideo(src) {
    const container = document.getElementById('player-container');
    container.innerHTML = `
      <video width="640" height="360" controls>
        <source src="${src}" type="video/mp4">
        Ваш браузер не поддерживает воспроизведение видео.
      </video>
    `;
    // Прокрутить к плееру
    container.scrollIntoView({ behavior: 'smooth' });
  }
  
  // При загрузке: если в URL есть #video.mp4, сразу покажем его
  document.addEventListener('DOMContentLoaded', () => {
    const hash = decodeURIComponent(window.location.hash.slice(1));
    if (hash) {
      showVideo(hash);
    }
  });
  