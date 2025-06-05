function combobox(tree = document) {
  tree.querySelectorAll("[data-combobox]").forEach(comboboxRoot => {
    const combobox = comboboxRoot.querySelector("[role=combobox]")
    const listbox = comboboxRoot.querySelector("[role=listbox]")
    const options = [...listbox.querySelectorAll("[role=option]")]

    const controller = new AbortController()
    const signal = controller.signal

    let isArrowNavigating = false

    const isOpen = () => !listbox.hidden
    
    const getSelectedOptionIndex = () => options.findIndex(option => option.getAttribute("aria-selected") === "true")

    comboboxRoot.addEventListener("focus", e => {
      toggleCombobox(true)
    }, { capture: true, signal })

    comboboxRoot.addEventListener("blur", e => {
      if (!comboboxRoot.contains(document.activeElement)) {
        toggleCombobox(false)
      }
    }, { capture: true, signal })

    comboboxRoot.addEventListener("keydown", e => {
      if (e.key == "ArrowDown" || e.key == "ArrowUp") {
        isArrowNavigating = true
        listbox.classList.add("pointer-events-none")

        const nextOptionIndex = e.key == "ArrowUp" ? getSelectedOptionIndex() - 1 : (getSelectedOptionIndex() + 1) % options.length
        const nextOption = options.at(nextOptionIndex)
        selectOption(nextOption)
        nextOption.scrollIntoView({ block: "nearest", inline: "nearest" })
        setTimeout(() => {
          isArrowNavigating = false
          listbox.classList.remove("pointer-events-none")
        }, 100)
      }

    }, { signal })

    listbox.addEventListener("htmx:beforeSwap", () => {
      controller.abort()
    })

    listbox.addEventListener("mouseover", e => {
      const option = e.target.closest("[role=option]")
      if (option) {
        selectOption(option)
      }
    }, { signal })

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
