function combobox(tree = document) {
  tree.querySelectorAll("[data-combobox]").forEach(comboboxRoot => {
    const combobox = comboboxRoot.querySelector("[role=combobox]")
    const listbox = comboboxRoot.querySelector("[role=listbox]")
    const options = [...listbox.querySelectorAll("[role=option]")]
    const dropdownContent = comboboxRoot.querySelector("#dropdown-content")

    const controller = new AbortController()
    const signal = controller.signal

    let isArrowNavigating = false

    const isOpen = () => combobox.getAttribute("aria-expanded") === "true"
    
    comboboxRoot.addEventListener("focus", e => {
      toggleCombobox(true)
    }, { capture: true, signal })

    comboboxRoot.addEventListener("blur", e => {
      const newFocusTarget = e.relatedTarget
      if (newFocusTarget && comboboxRoot.contains(newFocusTarget)) {
        return
      } 
      toggleCombobox(false)
    }, { capture: true, signal })

    comboboxRoot.addEventListener("keydown", e => {
      switch (e.key) {
        case "ArrowDown":
        case "ArrowUp":
          isArrowNavigating = true
          listbox.classList.add("pointer-events-none")
    
          const nextIndex = e.key == "ArrowUp" ? getSelectedOptionIndex() - 1 : (getSelectedOptionIndex() + 1) % options.length
          const nextOption = options.at(nextIndex)
          selectOption(nextOption)
          nextOption.scrollIntoView({ block: "nearest", inline: "nearest" })
          setTimeout(() => {
            isArrowNavigating = false
            listbox.classList.remove("pointer-events-none")
          }, 100)
          break
        case "Enter":
          console.log(`enter on active element: ${getSelectedOptionIndex()}`)
          const selectedIndex = getSelectedOptionIndex()
          if (selectedIndex === null) {
            break
          }
          const selectedOption = options[selectedIndex]
          const button = selectedOption.querySelector("button")
          if (button) {
            button.click()
            toggleCombobox(false)
          }
          break
        case "Escape":
          toggleCombobox(false)
          break
      }
    }, { signal })

    comboboxRoot.addEventListener("htmx:beforeSwap", () => {
      controller.abort()
    }, { signal })

    listbox.addEventListener("mouseover", e => {
      const option = e.target.closest("[role=option]")
      if (option) {
        selectOption(option)
      }
    }, { signal })

    listbox.addEventListener("click", e => {
      toggleCombobox(false)
    })

    toggleCombobox(isOpen())

    function toggleCombobox(open = !isOpen()) {
      if (open) {
        listbox.scrollTop = 0
        combobox.focus()
        combobox.setAttribute("aria-expanded", true)
        if (options.length == 0) {
          return
        }
        selectOption(options[0])
      } else {
        combobox.blur()
        combobox.setAttribute("aria-expanded", false)
        deselectAllOptions()
      }
    }

    function selectOption(option = options[0]) {
      deselectAllOptions()
      combobox.setAttribute("aria-activedescendant", option.id)
      option.setAttribute("aria-selected", true)
    }
    
    function getSelectedOptionIndex() {
      const activeOptionId = combobox.getAttribute("aria-activedescendant")
      if (!activeOptionId) {
        return null
      }
      return options.findIndex(o => o.getAttribute("id") === activeOptionId)
    }

    function deselectAllOptions() {
      combobox.removeAttribute("aria-activedescendant")
      options.forEach(o => {
        o.setAttribute("aria-selected", false)
      })
    }
  })
}

addEventListener("DOMContentLoaded", () => combobox())

addEventListener("htmx:afterSwap", (e) => {
  if (e.target.closest("[data-combobox]")) {
    combobox()
  }
})
