package plot

// Renderable is a function that can be called to render custom elements on the plot.
type Renderable func(r Renderer, canvasBox Box, defaults Style)
