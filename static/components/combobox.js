function combobox(tree = document) {
  tree.querySelectorAll("[data-combobox]").forEach(comboboxRoot => {
    const combobox = comboboxRoot.querySelector("[role=combobox]")
    const listbox = comboboxRoot.querySelector("[role=listbox]")
    const options = [...listbox.querySelectorAll("[role=option]")]

    const controller = new AbortController()
    const signal = controller.signal

    listbox.addEventListener("htmx:beforeSwap", () => {
      controller.abort()
    })

    const isOpen = () => !listbox.hidden
    
    const getSelectedOptionIndex = () => options.findIndex(option => option.getAttribute("aria-selected") === "true")

    comboboxRoot.addEventListener("focus", e => {
      toggleCombobox(true)
    }, { capture: true, signal: signal })

    comboboxRoot.addEventListener("blur", e => {
      if (!comboboxRoot.contains(document.activeElement)) {
        toggleCombobox(false)
      }
    }, { capture: true, signal: signal })

    comboboxRoot.addEventListener("keydown", e => {
      switch (e.key) {
        case "ArrowDown": {
          const nextOptionIndex = getSelectedOptionIndex() === options.length - 1 ? 0 : getSelectedOptionIndex() + 1
          const nextOption = options.at(nextOptionIndex)
          selectOption(nextOption)
          nextOption.scrollIntoView({ block: "nearest", inline: "nearest" })
          break
        }
        case "ArrowUp": {
          const prevOption = options.at(getSelectedOptionIndex() - 1)
          selectOption(prevOption)
          prevOption.scrollIntoView({ block: "nearest", inline: "nearest" })
          break
        }
      }
    }, { signal: signal })

    listbox.addEventListener("mouseover", e => {
      const option = e.target.closest("[role=option]")
      if (option) {
        selectOption(option)
      }
    }, { signal: signal })

    toggleCombobox(isOpen())

    function toggleCombobox(open = !isOpen()) {
      if (open) {
        listbox.hidden = false
        combobox.setAttribute("aria-expanded", true)
        selectOption(options[0])
      } else {
        listbox.hidden = true
        combobox.setAttribute("aria-expanded", false)
      }
    }

    function selectOption(option = options[0]) {
      combobox.setAttribute("aria-activedescendant", option.id)

      const unselectedOptions = options.filter(o => !Object.is(o, option))
      unselectedOptions.forEach(o => o.setAttribute("aria-selected", false))

      option.setAttribute("aria-selected", true)
    }
  })
}

addEventListener("DOMContentLoaded", () => combobox())
addEventListener("htmx:afterSwap", () => combobox())
