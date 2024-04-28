/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    'internal/web/views/**/*.templ'
  ],
  theme: {
    extend: {},
    fontFamily: {
      'sans': ['Roboto', 'sans-serif'],
      'roboto': ['Roboto', 'sans-serif']
    },
    // fontSize: {
    //     base: '14px',
    // },
    backgroundColor: '#fcfcfd'
  },
  plugins: []
}
