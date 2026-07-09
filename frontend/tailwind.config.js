/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {
      colors: {
        // brand = เขียวตลาดสด (design system)
        brand: {
          50: '#F0FDF4',
          100: '#DCFCE7',
          500: '#22C55E',
          600: '#16A34A',
          700: '#15803D',
        },
      },
      fontFamily: {
        sans: ['Sarabun', 'IBM Plex Sans Thai', 'system-ui', 'sans-serif'],
      },
      maxWidth: {
        // จัดกลางจอให้เหมือนมือถือ (mobile-first)
        app: '448px',
      },
    },
  },
  plugins: [],
}
