
class StampElement extends HTMLElement {
  constructor(){
    super()
    const text = this.innerText
    const shadowRoot = this.attachShadow({mode: "closed"})
    const wrapper = document.createElement("p")

    let bgColor = "#ea999ccc"
    let color = "#000000cc"

    if(this.hasAttribute("color")){
      color = this.getAttribute("color")
    }

    if(this.hasAttribute("background-color")){
      bgColor = this.getAttribute("bgColor")
    }

    let n = false
    for (let entry of text.split(' ')) {
      const span = document.createElement("span")
      span.style.backgroundColor = bgColor
      span.style.color = color
      span.style.marginRight = "0.25rem"
      span.textContent = entry
      let rand = Math.floor(Math.random()*15)
      if(n) {
        rand = -rand
      }

      span.style.rotate = `${rand}deg`
      span.style.display = "inline-block"

      wrapper.appendChild(span)
      n = !n
    }

    shadowRoot.append(wrapper)
  }
}
customElements.define("stam-p", StampElement)
