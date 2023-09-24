/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class', // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        'background': {
          'dark': '#121019',
          'form': '#231A34',
          DEFAULT: '#121019'
        },
        'secondary': {
          'dark': '#1C1726',
          DEFAULT: '#1C1726'
        },
        'primary': {
          'dark': '#FA1E4E',
          DEFAULT: '#FA1E4E'
        },
        'text-dark-theme': {
          'dark': '#8C7DB0',
          'darker': '#4A3E65',
          'auxiliar': '#000836',
          'light': '#BEB6D2',
          DEFAULT: '#8C7DB0',
        },
        'terminal': {
          DEFAULT: '#1C1726',
        }
      },
    }
  },
  plugins: [],
}

