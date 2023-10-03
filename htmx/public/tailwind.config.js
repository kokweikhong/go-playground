/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../views/**/*.{html,js}"],
  plugins: [require("@tailwindcss/forms")],
};
