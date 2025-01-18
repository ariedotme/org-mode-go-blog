class TodoElement extends HTMLElement {
  constructor(){
    super()
    const text = this.innerText
    const shadowRoot = this.attachShadow({mode: "closed"})
    const wrapper = document.createElement("span")
    wrapper.style.color = "#ef9f76"
    wrapper.textContent = `(TODO): ${text}`

    const style = document.createElement("style");
    style.textContent = `
      @keyframes blink {
        0% { background-color: #ffcccb; color: #000; }
        50% { background-color: transparent; color: #ef9f76;}
        100% { background-color: #ffcccb; color: #000;}
      }

      span {
        animation: blink 2s infinite;
      }
    `;

    shadowRoot.append(style,wrapper)
  }
}
customElements.define("to-do", TodoElement)
