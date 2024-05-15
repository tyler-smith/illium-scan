/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    'internal/web/views/*.templ',
    'internal/web/views/**/*.templ'
  ],
  theme: {
    extend: {},
    fontFamily: {
      'sans': ['Roboto', 'sans-serif'],
      'roboto': ['Roboto', 'sans-serif']
    },
    backgroundColor: '#fcfcfd'
  },
  plugins: []
}
