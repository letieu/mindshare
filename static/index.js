const themeToggle = document.getElementById("theme");
const body = document.body;

themeToggle.addEventListener('click', (event) => {
  event.preventDefault(); // Prevent the default link behavior
  if (body.classList.contains('dark-theme')) {
    body.classList.remove('dark-theme');
    themeToggle.textContent = 'Dark';
    localStorage.setItem('theme', 'light');
  } else {
    body.classList.add('dark-theme');
    themeToggle.textContent = 'Light';
    localStorage.setItem('theme', 'dark');
  }
});

// Load the saved theme
const savedTheme = localStorage.getItem('theme');
if (savedTheme === 'dark') {
  body.classList.add('dark-theme');
  themeToggle.textContent = 'Light';
}
