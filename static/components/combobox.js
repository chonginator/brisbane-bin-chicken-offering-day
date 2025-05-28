function combobox(tree = document) {
  tree.querySelectorAll("[data-combobox]").forEach(comboboxRoot => {
    const combobox = comboboxRoot.querySelector("[role=combobox]")
    const listbox = comboboxRoot.querySelector("[role=listbox]")
    const options = [...listbox.querySelectorAll("[role=option]")]

    const isOpen = () => !listbox.hidden

    comboboxRoot.addEventListener("focus", e => {
      toggleCombobox(true)
    }, { capture: true })

    comboboxRoot.addEventListener("blur", e => {
      if (!comboboxRoot.contains(document.activeElement)) {
        toggleCombobox(false)
      }
    }, { capture: true })
    
    toggleCombobox(isOpen())

    function toggleCombobox(open = !isOpen()) {
      if (open) {
        listbox.hidden = false
        combobox.setAttribute("aria-expanded", true)
      } else {
        listbox.hidden = true
        combobox.setAttribute("aria-expanded", false)
      }
    }

    options.forEach(option => {
      option.addEventListener("mouseover", () => {
        selectOption(option)
      })
    })

    function selectOption(option = options[0]) {
      combobox.setAttribute("aria-activedescendant", option.id)
      const unselectedOptions = options.filter(o => !Object.is(o, option))
      unselectedOptions.forEach(o => o.setAttribute("aria-selected", false))
      option.setAttribute("aria-selected", true)
    }
  })
}

addEventListener("htmx:load", e => combobox(e.target))
