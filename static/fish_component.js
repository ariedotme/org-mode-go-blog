class FishElement extends HTMLElement {
  constructor() {
    super();

    const shadowRoot = this.attachShadow({ mode: "closed" });

    const fishSpan = document.createElement("span");
    fishSpan.textContent = "𓆝";
    fishSpan.style.fontSize = "2rem";
    fishSpan.style.display = "inline-block";

    const style = document.createElement("style");

    style.textContent = `
span {
  position: relative;
  animation: tremble 0.2s infinite ease-in-out;
}

@keyframes tremble {
  0% {
    top: -0.05rem;
    left: -0.2rem;
  }
  25% {
    top: -0.1rem;
    left: -0.25rem;
  }
  50% {
    top: -0.05rem;
    left: -0.15rem;
  }
  75% {
    top: -0rem;
    left: -0.2rem;
  }
  100% {
    top: -0.05rem;
    left: -0.2rem;
  }
}`

    shadowRoot.append(style,fishSpan);
    const fishCharacters = ["𓆝", "𓆟", "𓆞", "𓆝", "𓆟"];

    let currentIndex = 0;
    if(this.hasAttribute("index")) {
      const index = Number(this.getAttribute("index"))
      if(fishCharacters.length > index > 0) {
        currentIndex = index;
      } else {
        let currentIndex = Math.floor(Math.random()*fishCharacters.length);
      }
    }
    fishSpan.textContent = fishCharacters[currentIndex];

    setInterval(() => {
      currentIndex = (currentIndex + 1) % fishCharacters.length;
      fishSpan.textContent = fishCharacters[currentIndex];
    }, 300);
  }
}

customElements.define("fi-sh", FishElement);
