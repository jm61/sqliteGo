const input = document.getElementById('searchInput')
const suggestionsList = document.getElementById('suggestions')
const form = document.getElementById('searchForm')

input.addEventListener('input', debounce(fetchSuggestions, 300))

function fetchSuggestions() {
  const query = input.value.trim()
  if (query.length < 2) {
    suggestionsList.style.display = 'none'
    return
  }

  fetch(`/search?q=${encodeURIComponent(query)}`)
    .then((response) => response.json())
    .then((data) => {
      suggestionsList.innerHTML = ''
      if (data.length > 0) {
        data.forEach((suggestion) => {
          const li = document.createElement('li')
          li.textContent = suggestion
          li.addEventListener('click', () => {
            input.value = suggestion
            suggestionsList.style.display = 'none'
          })
          suggestionsList.appendChild(li)
        })
        suggestionsList.style.display = 'block'
      } else {
        suggestionsList.style.display = 'none'
      }
    })
    .catch((error) => {
      console.error('Error fetching suggestions:', error)
      suggestionsList.style.display = 'none'
    })
}
// Debounce to avoid too many requests
function debounce(func, delay) {
  let timeoutId
  return function (...args) {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => func.apply(this, args), delay)
  }
}

// Close suggestions on outside click
document.addEventListener('click', (e) => {
  if (!form.contains(e.target)) {
    suggestionsList.style.display = 'none'
  }
})
